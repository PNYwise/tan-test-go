package domain

import geojson "github.com/paulmach/go.geojson"

type Geolocation struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
}

type IGeolocationService interface {
	CreateGeolocations(players *[]Geolocation) error
	GetGeolocationsGeoJSON() (*geojson.FeatureCollection, error)
}

type IGeolocationRepository interface {
	CreateBatch(*[]Geolocation) error
	GetGeolocations() (*[]Geolocation, error)
}
