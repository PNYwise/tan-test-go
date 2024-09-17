package service

import (
	"errors"
	"tan-test-go/internal/domain"

	geojson "github.com/paulmach/go.geojson"
)

type geolocationService struct {
	geolocationRepo domain.IGeolocationRepository
}

func NewGeolocationService(geolocationRepo domain.IGeolocationRepository) domain.IGeolocationService {
	return &geolocationService{geolocationRepo: geolocationRepo}
}

// CreateGeolocations implements domain.IGeolocationService.
func (g *geolocationService) CreateGeolocations(players *[]domain.Geolocation) error {
	if err := g.geolocationRepo.CreateBatch(players); err != nil {
		return errors.New("Internal Server error")
	}
	return nil
}

// GetGeolocationsGeoJSON implements domain.IGeolocationService.
func (g *geolocationService) GetGeolocationsGeoJSON() (*geojson.FeatureCollection, error) {
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

	return featureCollection, nil
}
