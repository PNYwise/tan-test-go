package repository

import (
	"context"
	"log"
	"tan-test-go/internal/domain"

	"github.com/jackc/pgx/v5"
)

type playerRepo struct {
	db  *pgx.Conn
	ctx context.Context
}

func NewPlayerRepository(ctx context.Context, db *pgx.Conn) domain.IPlayerRepository {
	return &playerRepo{db, ctx}
}

// CreateBatch implements domain.IPlayerRepository.
func (p *playerRepo) CreateBatch(players *[]domain.Player) error {
	tx, err := p.db.Begin(p.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(p.ctx)

	for _, player := range *players {
		_, err := tx.Exec(
			p.ctx, "INSERT INTO players (player_name, lat, lng) VALUES ($1, $2, $3)",
			player.PlayerName, player.Lat, player.Lng)
		if err != nil {
			log.Printf("Error inserting player: %v", err)
			return err
		}
	}

	err = tx.Commit(p.ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetPlayers implements domain.IPlayerRepository.
func (p *playerRepo) GetPlayers() (*[]domain.Player, error) {
	rows, err := p.db.Query(p.ctx, "SELECT id, player_name, lat, lng FROM players")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []domain.Player
	for rows.Next() {
		var player domain.Player
		if err := rows.Scan(&player.ID, &player.PlayerName, &player.Lat, &player.Lng); err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return &players, nil
}
