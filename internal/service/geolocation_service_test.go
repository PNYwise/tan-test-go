package service_test

import (
	"encoding/json"
	"errors"
	"tan-test-go/internal/config"
	"tan-test-go/internal/domain"
	mock_test "tan-test-go/internal/repository/_mock"
	"tan-test-go/internal/service"
	"testing"
	"time"

	geojson "github.com/paulmach/go.geojson"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateGeolocations(t *testing.T) {
	mockRepo := new(mock_test.MockGeolocationRepository)
	mockValidator := new(config.MockValidator)
	mockRedisRepo := new(mock_test.MockRedisRepository)

	geolocationService := service.NewGeolocationService(mockRepo, mockRedisRepo, mockValidator)

	geolocations := []domain.Geolocation{
		{Name: "John Doe", Description: "test Geolocation", Lat: 37.7749, Lng: -122.4194},
	}

	t.Run("Validation Error", func(t *testing.T) {
		// Set up expectations
		mockValidator.On("ValidateStruct", geolocations).Return(errors.New("validation error")).Once()

		// Call the method
		err := geolocationService.CreateGeolocations(geolocations)

		// Assertions
		assert.Error(t, err)

		// Verify expectations
		mockValidator.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		// Set up expectations
		mockValidator.On("ValidateStruct", geolocations).Return(nil).Once()
		mockRepo.On("CreateBatch", &geolocations).Return(errors.New("repository error")).Once()

		// Call the method
		err := geolocationService.CreateGeolocations(geolocations)

		// Assertions
		assert.Error(t, err)

		// Verify expectations
		mockValidator.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Valid Geolocations", func(t *testing.T) {
		// Set up expectations
		mockValidator.On("ValidateStruct", geolocations).Return(nil).Once()
		mockRepo.On("CreateBatch", &geolocations).Return(nil).Once()

		// Call the method
		err := geolocationService.CreateGeolocations(geolocations)

		// Assertions
		assert.NoError(t, err)

		// Verify expectations
		mockValidator.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetGeolocationsGeoJSON(t *testing.T) {
	mockRepo := new(mock_test.MockGeolocationRepository)
	mockValidator := new(config.MockValidator)
	mockRedisRepo := new(mock_test.MockRedisRepository)

	geolocationService := service.NewGeolocationService(mockRepo, mockRedisRepo, mockValidator)

	geolocations := []domain.Geolocation{
		{Name: "John Doe", Description: "Test Location", Lat: 37.7749, Lng: -122.4194},
	}

	t.Run("Cache Miss", func(t *testing.T) {
		// Set up expectations
		mockRedisRepo.On("Get", "geolocations:geojson").Return("", redis.Nil).Once()
		mockRepo.On("GetGeolocations").Return(&geolocations, nil).Once()
		mockRedisRepo.On("Set", "geolocations:geojson", mock.AnythingOfType("string"), 60*time.Second).Return(nil).Once()

		// Call the method
		geojsonData, err := geolocationService.GetGeolocationsGeoJSON()

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, geojsonData)

		// Verify expectations
		mockRedisRepo.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Cache Hit", func(t *testing.T) {
		// Create a GeoJSON FeatureCollection for testing
		featureCollection := geojson.NewFeatureCollection()
		point := geojson.NewPointFeature([]float64{-122.4194, 37.7749})
		point.Properties["name"] = "John Doe"
		point.Properties["description"] = "Test Location"
		featureCollection.AddFeature(point)
		geojsonBytes, _ := json.Marshal(featureCollection)

		// Set up expectations
		mockRedisRepo.On("Get", "geolocations:geojson").Return(string(geojsonBytes), nil).Once()

		// Call the method
		geojsonData, err := geolocationService.GetGeolocationsGeoJSON()

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, geojsonData)

		// Verify expectations
		mockRedisRepo.AssertExpectations(t)
	})

	t.Run("Cache Error", func(t *testing.T) {
		// Set up expectations
		mockRedisRepo.On("Get", "geolocations:geojson").Return("", errors.New("cache error")).Once()

		// Call the method
		geojsonData, err := geolocationService.GetGeolocationsGeoJSON()

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, geojsonData)

		// Verify expectations
		mockRedisRepo.AssertExpectations(t)
	})
}
