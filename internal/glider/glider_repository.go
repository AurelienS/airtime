package glider

import (
	"context"

	"github.com/AurelienS/cigare/internal/log"
	"github.com/AurelienS/cigare/internal/storage"
)

type GliderRepository interface {
	GetGliders(ctx context.Context, user storage.User) ([]storage.Glider, error)
}

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
