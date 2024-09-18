package service

import (
	"encoding/json"
	"errors"
	"tan-test-go/internal/config"
	"tan-test-go/internal/domain"
	"time"

	geojson "github.com/paulmach/go.geojson"
	"github.com/redis/go-redis/v9"
)

type geolocationService struct {
	geolocationRepo domain.IGeolocationRepository
	cache           domain.RedisCacheRepository
	validator       config.Validator
}

func NewGeolocationService(geolocationRepo domain.IGeolocationRepository, cache domain.RedisCacheRepository, validator config.Validator) domain.IGeolocationService {
	return &geolocationService{geolocationRepo, cache, validator}
}

// CreateGeolocations implements domain.IGeolocationService.
func (g *geolocationService) CreateGeolocations(geolocation []domain.Geolocation) error {
	if err := g.validator.ValidateStruct(geolocation); err != nil {
		return err
	}

	if err := g.geolocationRepo.CreateBatch(&geolocation); err != nil {
		return errors.New("Internal Server error")
	}
	return nil
}

// GetGeolocationsGeoJSON implements domain.IGeolocationService.
func (g *geolocationService) GetGeolocationsGeoJSON() (*geojson.FeatureCollection, error) {
	cacheKey := "geolocations:geojson"

	// Check if data is in Redis
	val, err := g.cache.Get(cacheKey)
	if err == redis.Nil { // Data not in Redis
		geolocations, err := g.geolocationRepo.GetGeolocations()
		if err != nil {
			return nil, errors.New("Internal Server error")
		}

		// Create GeoJSON
		featureCollection := geojson.NewFeatureCollection()
		for _, geolocation := range *geolocations {
			point := geojson.NewPointFeature([]float64{geolocation.Lng, geolocation.Lat})
			point.Properties["name"] = geolocation.Name
			point.Properties["description"] = geolocation.Description
			featureCollection.AddFeature(point)
		}

		// Cache the serialized result in Redis with a 60-second expiration
		geojsonBytes, err := json.Marshal(featureCollection)
		if err != nil {
			return nil, errors.New("Failed to marshal GeoJSON")
		}
		if err := g.cache.Set(cacheKey, string(geojsonBytes), 60*time.Second); err != nil {
			return nil, errors.New("Failed to cache GeoJSON")
		}

		return featureCollection, nil
	} else if err != nil {
		return nil, errors.New("Failed to fetch data from cache")
	}

	// Data found in Redis, unmarshal it and return
	var cachedGeoJSON geojson.FeatureCollection
	if err := json.Unmarshal([]byte(val), &cachedGeoJSON); err != nil {
		return nil, errors.New("Failed to unmarshal cached data")
	}

	return &cachedGeoJSON, nil
}
