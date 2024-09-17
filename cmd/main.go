package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"tan-test-go/internal/config"
	"tan-test-go/internal/handler"
	"tan-test-go/internal/repository"
	"tan-test-go/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	createLogDirectoryIfNotExists()

	// Initialize the logger
	logger, err := config.NewLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf("Failed to flush logs: %v", err)
		}
	}()

	// Initialize Fiber app
	app := fiber.New()
	app.Use(cors.New(cors.Config{}))

	// Load configuration
	conf := config.New()

	ctx := context.Background()
	// Initialize the database connection
	db := config.DbConn(ctx, conf)
	defer db.Close(ctx)

	// Initialize Redis connection
	redis := config.RedisConn(ctx, conf)

	// Initialize repository and service with logger
	geolocationRepo := repository.NewGeolocationRepository(ctx, db, logger)
	geolocationService := service.NewGeolocationService(geolocationRepo)

	// Initialize handler with logger
	geolocationHandler := handler.NewGeolocationHandler(geolocationService, redis)

	// Define routes
	api := app.Group("/api")
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong!")
	})
	api.Post("/create-batch", geolocationHandler.CreateGeolocations)
	api.Get("/map-data", geolocationHandler.GetGeolocationsGeoJSON)

	// Start the server
	appPort := conf.GetString("app.port")
	appHost := conf.GetString("app.host")
	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", appHost, appPort)))
}

func createLogDirectoryIfNotExists() {
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
		os.Exit(1)
	}
}
