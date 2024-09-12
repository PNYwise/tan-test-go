package service

import (
	"context"
	"encoding/json"
	"tan-test-go/internal/domain"
	"time"

	geojson "github.com/paulmach/go.geojson"
	"github.com/redis/go-redis/v9"
)

type playerService struct {
	playerRepo domain.IPlayerRepository
	rdb        *redis.Client
}

func NewPlayerService(playerRepo domain.IPlayerRepository, rdb *redis.Client) domain.IPlayerService {
	return &playerService{playerRepo: playerRepo, rdb: rdb}
}

// CreatePlayers implements domain.IPlayerService.
func (p *playerService) CreatePlayers(players *[]domain.Player) error {
	return p.playerRepo.CreateBatch(players)
}

// GetPlayersGeoJSON implements domain.IPlayerService.
func (p *playerService) GetPlayersGeoJSON() (string, error) {
	cacheKey := "players_geojson"

	ctx := context.Background()
	// Check if data is in Redis
	val, err := p.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		players, err := p.playerRepo.GetPlayers()
		if err != nil {
			return "", err
		}

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

		geojsonData := string(geojsonBytes)

		// Cache with a 30-second expiration
		p.rdb.Set(ctx, cacheKey, geojsonData, 30*time.Second)

		return geojsonData, nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}
