package repository

import (
	"context"
	"tan-test-go/internal/domain"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type geolocationRepo struct {
	db     *pgx.Conn
	ctx    context.Context
	logger *zap.Logger
}

func NewGeolocationRepository(ctx context.Context, db *pgx.Conn, logger *zap.Logger) domain.IGeolocationRepository {
	return &geolocationRepo{db, ctx, logger}
}

// Assuming zap logger is initialized in geolocationRepo struct
// For example, add it as a field: logger *zap.Logger

// CreateBatch implements domain.IGeolocationRepository.
func (g *geolocationRepo) CreateBatch(geolocations *[]domain.Geolocation) error {
	tx, err := g.db.Begin(g.ctx)
	if err != nil {
		g.logger.Error("Error starting transaction", zap.Error(err))
		return err
	}
	defer func() {
		if err := tx.Rollback(g.ctx); err != nil {
			g.logger.Error("Error rolling back transaction", zap.Error(err))
		}
	}()

	for _, geolocation := range *geolocations {
		_, err := tx.Exec(
			g.ctx, "INSERT INTO geolocations (name, description, lat, lng) VALUES ($1, $2, $3, $4)",
			geolocation.Name, geolocation.Description, geolocation.Lat, geolocation.Lng)
		if err != nil {
			g.logger.Error("Error inserting geolocation", zap.Error(err), zap.Any("geolocation", geolocation))
			return err
		}
	}

	err = tx.Commit(g.ctx)
	if err != nil {
		g.logger.Error("Error committing transaction", zap.Error(err))
		return err
	}
	return nil
}

// GetGeolocations implements domain.IPlayerRepository.
func (g *geolocationRepo) GetGeolocations() (*[]domain.Geolocation, error) {
	rows, err := g.db.Query(g.ctx, "SELECT id, name, description1, lat, lng FROM geolocations")
	if err != nil {
		g.logger.Error("Error fetching geolocations", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var geolocations []domain.Geolocation
	for rows.Next() {
		var geolocation domain.Geolocation
		if err := rows.Scan(&geolocation.ID, &geolocation.Name, &geolocation.Description, &geolocation.Lat, &geolocation.Lng); err != nil {
			g.logger.Error("Error scanning geolocation row", zap.Error(err))
			return nil, err
		}
		geolocations = append(geolocations, geolocation)
	}

	return &geolocations, nil
}
