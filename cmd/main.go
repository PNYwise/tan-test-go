package main

import (
	"context"
	"fmt"
	"log"
	"tan-test-go/internal/config"
	"tan-test-go/internal/handler"
	"tan-test-go/internal/repository"
	"tan-test-go/internal/service"

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

	// Initialize repository and service
	playerRepo := repository.NewPlayerRepostory(ctx, db)
	playerService := service.NewPlayerService(playerRepo, redis)

	// Initialize player handler
	playerHandler := handler.NewPlayerHandler(playerService)

	// Routes
	api := app.Group("/api")
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong!")
	})
	api.Post("/create-batch", playerHandler.CreatePlayers)
	api.Get("/map-data", playerHandler.GetPlayersGeoJSON)

	appPort := conf.GetString("app.port")
	appHost := conf.GetString("app.host")
	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", appHost, appPort)))
}
