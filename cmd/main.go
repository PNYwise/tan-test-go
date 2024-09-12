package main

import (
	"context"
	"fmt"
	"log"
	"tan-test-go/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	app := fiber.New()

	app.Use(cors.New(cors.Config{}))

	conf := config.New()

	ctx := context.Background()
	// Initialize the db
	db := config.DbConn(ctx, conf)
	defer db.Close(ctx)

	redis := config.RedisConn(ctx, conf)
	_ = redis

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong!")
	})

	appPort := conf.GetString("app.port")
	appHost := conf.GetString("app.host")
	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", appHost, appPort)))
}
