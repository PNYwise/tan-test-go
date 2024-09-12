package handler

import (
	"context"
	"tan-test-go/internal/domain"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type PlayerHandler struct {
	playerService domain.IPlayerService
	rdb           *redis.Client
}

func NewPlayerHandler(playerService domain.IPlayerService, rdb *redis.Client) *PlayerHandler {
	return &PlayerHandler{playerService: playerService, rdb: rdb}
}

// CreatePlayers handles player creation
func (p *PlayerHandler) CreatePlayers(c *fiber.Ctx) error {
	var input struct {
		Items []domain.Player `json:"items"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := p.playerService.CreatePlayers(&input.Items)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create players"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Players created successfully"})
}

func (p *PlayerHandler) GetPlayersGeoJSON(c *fiber.Ctx) error {
	ctx := context.Background()
	cacheKey := c.Path()

	// Check if data is in Redis
	val, err := p.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil { // Data not in Redis
		geojsonData, err := p.playerService.GetPlayersGeoJSON()
		if err != nil {
			return c.Status(500).SendString("Failed to fetch players")
		}

		// Cache result in Redis with a 30-second expiration
		p.rdb.Set(ctx, cacheKey, geojsonData, 30*time.Second)

		return c.Type("application/json").SendString(geojsonData)
	} else if err != nil {
		return c.Status(500).SendString("Failed to fetch data from cache")
	}

	// Data found in Redis, return it
	return c.Type("application/json").SendString(val)
}
