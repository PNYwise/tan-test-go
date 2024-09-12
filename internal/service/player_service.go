package service

import (
	"encoding/json"
	"tan-test-go/internal/domain"

	geojson "github.com/paulmach/go.geojson"
)

type playerService struct {
	playerRepo domain.IPlayerRepository
}

func NewPlayerService(playerRepo domain.IPlayerRepository) domain.IPlayerService {
	return &playerService{playerRepo: playerRepo}
}

// CreatePlayers implements domain.IPlayerService.
func (p *playerService) CreatePlayers(players *[]domain.Player) error {
	return p.playerRepo.CreateBatch(players)
}

// GetPlayersGeoJSON implements domain.IPlayerService.
func (p *playerService) GetPlayersGeoJSON() (string, error) {
	players, err := p.playerRepo.GetPlayers()
	if err != nil {
		return "", err
	}
	// Create GeoJSON
	featureCollection := geojson.NewFeatureCollection()
	for _, player := range *players {
		point := geojson.NewPointFeature([]float64{player.Lng, player.Lat})
		point.Properties["player_name"] = player.PlayerName
		featureCollection.AddFeature(point)
	}

	geojsonBytes, err := json.Marshal(featureCollection)
	if err != nil {
		return "", err
	}

	return string(geojsonBytes), nil
}
