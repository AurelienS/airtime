package storage

import (
	"context"

	"github.com/AurelienS/cigare/internal/storage"
	"github.com/rs/zerolog/log"
)

type SQLGliderRepository struct {
	Queries storage.Queries
}

func (repo SQLGliderRepository) GetGliders(ctx context.Context, user storage.User) ([]storage.Glider, error) {
	flights, err := repo.Queries.GetGliders(context.Background(), user.ID)
	if err != nil {
		log.Error().Err(err).Str("user", user.Email).Msg("Failed to get gliders")
	}
	return flights, err
}
