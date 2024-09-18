package handler

import (
	"context"
	"encoding/json"
	"tan-test-go/internal/config"
	"tan-test-go/internal/domain"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type GeolocationHandler struct {
	geolocationService domain.IGeolocationService
	rdb                *redis.Client
}

func NewGeolocationHandler(geolocationService domain.IGeolocationService, rdb *redis.Client) *GeolocationHandler {
	return &GeolocationHandler{geolocationService: geolocationService, rdb: rdb}
}

// CreateGeolocation handles player creation
func (g *GeolocationHandler) CreateGeolocations(c *fiber.Ctx) error {
	var input struct {
		Items []domain.Geolocation `json:"items"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := g.geolocationService.CreateGeolocations(input.Items)
	if err != nil {
		if config.IsValidationError(err) {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})

		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Geolocation created successfully"})
}

func (g *GeolocationHandler) GetGeolocationsGeoJSON(c *fiber.Ctx) error {
	ctx := context.Background()
	cacheKey := c.Path()

	// Check if data is in Redis
	val, err := g.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil { // Data not in Redis
		geojsonData, err := g.geolocationService.GetGeolocationsGeoJSON()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Cache the serialized result in Redis with a 60-second expiration
		geojsonBytes, err := json.Marshal(geojsonData)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to marshal GeoJSON",
			})
		}
		g.rdb.Set(ctx, cacheKey, geojsonBytes, 60*time.Second)

		// Return the JSON response
		return c.Status(200).JSON(geojsonData)
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch data from cache",
		})
	}

	// Data found in Redis, unmarshal it and return as JSON
	var cachedGeoJSON interface{}
	if err := json.Unmarshal([]byte(val), &cachedGeoJSON); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to unmarshal cached data",
		})
	}

	return c.Status(200).JSON(cachedGeoJSON)
}
