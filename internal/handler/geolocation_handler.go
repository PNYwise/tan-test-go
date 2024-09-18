package handler

import (
	"tan-test-go/internal/config"
	"tan-test-go/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type GeolocationHandler struct {
	geolocationService domain.IGeolocationService
}

func NewGeolocationHandler(geolocationService domain.IGeolocationService) *GeolocationHandler {
	return &GeolocationHandler{geolocationService: geolocationService}
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
	geojsonData, err := g.geolocationService.GetGeolocationsGeoJSON()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(geojsonData)
}
