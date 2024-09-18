package service_test

import (
	"errors"
	"tan-test-go/internal/config"
	"tan-test-go/internal/domain"
	mock_test "tan-test-go/internal/repository/_mock"
	"tan-test-go/internal/service"
	"testing"

	geojson "github.com/paulmach/go.geojson"
	"github.com/stretchr/testify/assert"
)

func TestCreateGeolocations(t *testing.T) {
	mockRepo := new(mock_test.MockGeolocationRepository)
	mockValidator := new(config.MockValidator)

	geolocationService := service.NewGeolocationService(mockRepo, mockValidator)

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

func TestGetGeolocationGeoJSON(t *testing.T) {
	mockRepo := new(mock_test.MockGeolocationRepository)
	mockValidator := new(config.MockValidator)
	geolocationService := service.NewGeolocationService(mockRepo, mockValidator)

	t.Run("Valid Geolocations", func(t *testing.T) {
		geolocations := []domain.Geolocation{
			{Name: "My home", Description: "test Geolocation 1", Lat: 37.7749, Lng: -122.4194},
			{Name: "my home 2", Description: "test Geolocation 2", Lat: 34.0522, Lng: -118.2437},
		}

		// Mock the repository to return geolocations
		mockRepo.On("GetGeolocations").Return(&geolocations, nil)

		// Create the expected GeoJSON FeatureCollection
		expectedFeatureCollection := geojson.NewFeatureCollection()
		for _, geolocation := range geolocations {
			point := geojson.NewPointFeature([]float64{geolocation.Lng, geolocation.Lat})
			point.Properties["name"] = geolocation.Name
			point.Properties["description"] = geolocation.Description
			expectedFeatureCollection.AddFeature(point)
		}

		// Call the service
		geojsonData, err := geolocationService.GetGeolocationsGeoJSON()
		assert.NoError(t, err)

		// Compare the returned FeatureCollection with the expected one
		assert.Equal(t, expectedFeatureCollection, geojsonData)

		// Ensure the mock expectations were met
		mockRepo.AssertExpectations(t)
	})

}
