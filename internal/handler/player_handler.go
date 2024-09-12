package handler

import (
	"tan-test-go/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type PlayerHandler struct {
	playerService domain.IPlayerService
}

func NewPlayerHandler(playerService domain.IPlayerService) *PlayerHandler {
	return &PlayerHandler{playerService: playerService}
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
	geojsonData, err := p.playerService.GetPlayersGeoJSON()
	if err != nil {
		return c.Status(500).SendString("Failed to fetch players")
	}

	return c.Type("application/json").SendString(geojsonData)
}
