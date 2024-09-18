package domain

import geojson "github.com/paulmach/go.geojson"

type Geolocation struct {
	ID          uint    `json:"id" validate:"omitempty,gt=0"`
	Name        string  `json:"name" validate:"required,min=2,max=255"`
	Description string  `json:"description" validate:"omitempty"`
	Lat         float64 `json:"lat" validate:"required,gte=-90,lte=90"`
	Lng         float64 `json:"lng" validate:"required,gte=-180,lte=180"`
}

type IGeolocationService interface {
	CreateGeolocations(players []Geolocation) error
	GetGeolocationsGeoJSON() (*geojson.FeatureCollection, error)
}

type IGeolocationRepository interface {
	CreateBatch(*[]Geolocation) error
	GetGeolocations() (*[]Geolocation, error)
}
