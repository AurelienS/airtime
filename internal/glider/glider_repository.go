package glider

import (
	"context"

	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/util"
)

type GliderRepository interface {
	GetGliders(ctx context.Context, user storage.User) ([]storage.Glider, error)
	AddGlider(ctx context.Context, gliderName string, user storage.User) error
}

type SQLGliderRepository struct {
	Queries storage.Queries
}

func NewSQLGliderRepository(queries storage.Queries) SQLGliderRepository {
	return SQLGliderRepository{
		Queries: queries,
	}
}

func (repo SQLGliderRepository) GetGliders(ctx context.Context, user storage.User) ([]storage.Glider, error) {
	gliders, err := repo.Queries.GetGliders(ctx, user.ID)
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to get gliders")
	}
	return gliders, err
}

func (repo SQLGliderRepository) AddGlider(ctx context.Context, gliderName string, user storage.User) error {
	arg := storage.InsertGliderParams{
		Name:   gliderName,
		UserID: user.ID,
	}
	err := repo.Queries.InsertGlider(ctx, arg)
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to get gliders")
	}
	return err
}
