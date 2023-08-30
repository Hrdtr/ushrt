package db

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Hrdtr/ushrt/config"
	_ "github.com/lib/pq"
)

//go:generate sqlc generate

var (
	DB *sql.DB
	Q  = New(Get())
)

func Get() *sql.DB {
	if DB == nil {
		POSTGRES_PORT, _ := strconv.ParseUint(config.POSTGRES_PORT, 10, 32)
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", config.POSTGRES_HOST, POSTGRES_PORT, config.POSTGRES_USER, config.POSTGRES_PASSWORD, config.POSTGRES_DB, config.POSTGRES_SSL_MODE)
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			panic(err.Error())
		}
		DB = db
	}

	return DB
}
