package main

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"time"
	"unsafe"

	"github.com/Hrdtr/ushrt/config"
	"github.com/Hrdtr/ushrt/db"
	_ "github.com/Hrdtr/ushrt/docs"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/jellydator/ttlcache/v3"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func InternalServerError(c *fiber.Ctx, e error) error {
	if config.APP_ENV == "development" {
		println(e.Error())
	}
	return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "An unexpected error occurred while processing your request",
	})
}

type GetOrCreateLinkRequestBody struct {
	ID          string `json:"id" validate:"omitempty,max=32"`
	OriginalUrl string `json:"original_url" format:"url" validate:"required,url"`
}
type GetOrCreateLinkResponse struct {
	db.Link
	ShortenedUrl string `json:"shortened_url" format:"url"`
}

var validate = validator.New()
var cache *ttlcache.Cache[string, db.Link]

// @title			UShrt
// @version		1.0
// @description	Dead simple headless url shortener service for your app
// @contact.name	Herdi Tr.
// @contact.email	iam@icm.hrdtr.dev
// @BasePath		/
func main() {
	if config.APP_CACHE_TTL != "" {
		ttl, err := time.ParseDuration(config.APP_CACHE_TTL)
		if err != nil {
			panic(err)
		}
		cache = ttlcache.New[string, db.Link](
			ttlcache.WithTTL[string, db.Link](ttl),
		)
		go cache.Start() // starts automatic expired item deletion
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.APP_CORS_ALLOW_ORIGINS,
		AllowHeaders: "Origin, Content-Type, Accept, UShrt-API-Key",
	}))
	app.Get("/swagger/*", swagger.New(swagger.Config{
		CustomStyle: "body { margin: 0 } .swagger-ui .topbar { display: none } .swagger-ui .info { margin-top: 20px }",
	}))
	InitRoutes(app)

	log.Fatal(app.Listen(":3000"))
}

func InitRoutes(app *fiber.App) {
	app.Get("/:id", ResolveURL)

	api := app.Group("/api")
	api.Use(func(c *fiber.Ctx) error {
		if c.Get("UShrt-API-Key") != config.APP_API_KEY {
			return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "Invalid API key",
			})
		}
		return c.Next()
	})
	api.Post("/", GetOrCreateLink)
}

// ResolveURL godoc
// @Summary	Resolve original url
// @Param		id	path	string	true	"Link ID"
// @Success	301
// @Failure	404	{object}	ErrorResponse
// @Router		/{id}  [get]
func ResolveURL(c *fiber.Ctx) error {
	if cache != nil && cache.Has("ID:"+c.Params("id")) {
		item := cache.Get("ID:" + c.Params("id"))
		return c.Redirect(item.Value().OriginalUrl, 301)
	}

	link, err := db.Q.GetLink(context.Background(), c.Params("id"))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "Requested link not found",
			})
		}
		return InternalServerError(c, err)
	}
	cache.Set("ID:"+c.Params("id"), link, ttlcache.DefaultTTL)
	return c.Redirect(link.OriginalUrl, 301)
}

func RandStringBytesMaskImprSrcUnsafe(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6
		letterIdxMask = 1<<letterIdxBits - 1
		letterIdxMax  = 63 / letterIdxBits
	)
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

// GetOrCreateLink godoc
// @Summary	Retrieve an existing short link or create a new one
// @Tags		api
// @Param		UShrt-API-Key	header	string	true	"API Key"
// @Param		body	body	GetOrCreateLinkRequestBody	true	"Payload"
// @Success	200 {object}  GetOrCreateLinkResponse
// @Failure	400	{object}	ErrorResponse
// @Failure	401	{object}	ErrorResponse
// @Router		/api   [post]
func GetOrCreateLink(c *fiber.Ctx) error {
	payload := new(GetOrCreateLinkRequestBody)
	if err := c.BodyParser(payload); err != nil {
		return InternalServerError(c, err)
	}
	if err := validate.Struct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
	}

	if cache != nil && cache.Has("OriginalUrl:"+payload.OriginalUrl) {
		item := cache.Get("OriginalUrl:" + payload.OriginalUrl)
		return c.Status(fiber.StatusOK).JSON(GetOrCreateLinkResponse{
			Link:         item.Value(),
			ShortenedUrl: config.APP_BASE_URL + "/" + item.Value().ID,
		})
	}

	link, err := db.Q.GetLinkByOriginalUrl(context.Background(), payload.OriginalUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			newID := payload.ID
			if newID == "" {
				newID = RandStringBytesMaskImprSrcUnsafe(10)
			}
			newlyCreatedLink, err := db.Q.CreateLink(context.Background(), db.CreateLinkParams{
				ID:          newID,
				OriginalUrl: payload.OriginalUrl,
			})
			if err != nil {
				return InternalServerError(c, err)
			}
			return c.Status(fiber.StatusOK).JSON(GetOrCreateLinkResponse{
				Link:         newlyCreatedLink,
				ShortenedUrl: config.APP_BASE_URL + "/" + newlyCreatedLink.ID,
			})
		}
		return InternalServerError(c, err)
	}
	cache.Set("OriginalUrl:"+payload.OriginalUrl, link, ttlcache.DefaultTTL)
	return c.Status(fiber.StatusOK).JSON(GetOrCreateLinkResponse{
		Link:         link,
		ShortenedUrl: config.APP_BASE_URL + "/" + link.ID,
	})
}
