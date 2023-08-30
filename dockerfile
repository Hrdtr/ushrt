FROM golang:1.21.0 AS build

WORKDIR /go/src/ushrt
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .


FROM alpine:latest as release

WORKDIR /ushrt
COPY --from=build /go/src/ushrt/app .
RUN apk -U upgrade \
  && apk add --no-cache dumb-init ca-certificates \
  && chmod +x /ushrt/app
EXPOSE 3000

ENTRYPOINT ["/usr/bin/dumb-init", "--"]