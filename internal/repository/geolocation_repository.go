package repository

import (
	"context"
	"log"
	"tan-test-go/internal/domain"

	"github.com/jackc/pgx/v5"
)

type geolocationRepo struct {
	db  *pgx.Conn
	ctx context.Context
}

func NewGeolocationRepository(ctx context.Context, db *pgx.Conn) domain.IGeolocationRepository {
	return &geolocationRepo{db, ctx}
}

// CreateBatch implements domain.IGeolocationRepository.
func (g *geolocationRepo) CreateBatch(geolocations *[]domain.Geolocation) error {
	tx, err := g.db.Begin(g.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(g.ctx)

	for _, geolocation := range *geolocations {
		_, err := tx.Exec(
			g.ctx, "INSERT INTO geolocations (name, description, lat, lng) VALUES ($1, $2, $3, $4)",
			geolocation.Name, geolocation.Description, geolocation.Lat, geolocation.Lng)
		if err != nil {
			log.Printf("Error inserting geolocation: %v", err)
			return err
		}
	}

	err = tx.Commit(g.ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetGeolocations implements domain.IPlayerRepository.
func (g *geolocationRepo) GetGeolocations() (*[]domain.Geolocation, error) {
	rows, err := g.db.Query(g.ctx, "SELECT id, name, description, lat, lng FROM geolocations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var geolocations []domain.Geolocation
	for rows.Next() {
		var geolocation domain.Geolocation
		if err := rows.Scan(&geolocation.ID, &geolocation.Name, &geolocation.Description, &geolocation.Lat, &geolocation.Lng); err != nil {
			return nil, err
		}
		geolocations = append(geolocations, geolocation)
	}

	return &geolocations, nil
}
