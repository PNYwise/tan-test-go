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
	geolocationRepo := repository.NewGeolocationRepository(ctx, db)
	geolocationService := service.NewGeolocationService(geolocationRepo)

	// Initialize handler
	geolocationHandler := handler.NewGeolocationHandler(geolocationService, redis)

	// Routes
	api := app.Group("/api")
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong!")
	})
	api.Post("/create-batch", geolocationHandler.CreateGeolocations)
	api.Get("/map-data", geolocationHandler.GetGeolocationsGeoJSON)

	appPort := conf.GetString("app.port")
	appHost := conf.GetString("app.host")
	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", appHost, appPort)))
}
