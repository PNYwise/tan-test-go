package service_test

import (
	"encoding/json"
	"tan-test-go/internal/domain"
	mock_test "tan-test-go/internal/repository/_mock"
	"tan-test-go/internal/service"
	"testing"

	geojson "github.com/paulmach/go.geojson"
	"github.com/stretchr/testify/assert"
)

func TestCreatePlayers(t *testing.T) {
	mockRepo := new(mock_test.MockPlayerRepository)
	playerService := service.NewPlayerService(mockRepo)

	players := []domain.Player{
		{PlayerName: "John Doe", Lat: 37.7749, Lng: -122.4194},
	}

	mockRepo.On("CreateBatch", &players).Return(nil)

	err := playerService.CreatePlayers(&players)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestGetPlayersGeoJSON(t *testing.T) {
	mockRepo := new(mock_test.MockPlayerRepository)
	playerService := service.NewPlayerService(mockRepo)

	players := []domain.Player{
		{PlayerName: "John Doe", Lat: 37.7749, Lng: -122.4194},
		{PlayerName: "Jane Doe", Lat: 34.0522, Lng: -118.2437},
	}

	mockRepo.On("GetPlayers").Return(&players, nil)

	featureCollection := geojson.NewFeatureCollection()
	for _, player := range players {
		point := geojson.NewPointFeature([]float64{player.Lng, player.Lat})
		point.Properties["player_name"] = player.PlayerName
		featureCollection.AddFeature(point)
	}
	expectedGeoJSON, _ := json.Marshal(featureCollection)
	expectedGeoJSONString := string(expectedGeoJSON)

	geojsonData, err := playerService.GetPlayersGeoJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, expectedGeoJSONString, geojsonData)

	mockRepo.AssertExpectations(t)
}
