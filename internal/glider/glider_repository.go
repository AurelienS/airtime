package glider

import (
	"context"

	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/user"
	"github.com/AurelienS/cigare/internal/util"
)

type Repository struct {
	Queries storage.Queries
}

func NewRepository(queries storage.Queries) Repository {
	return Repository{
		Queries: queries,
	}
}

func ConvertGliderDBToGlider(gliderDB storage.Glider) Glider {
	var glider Glider

	glider.ID = int(gliderDB.ID)
	glider.Name = gliderDB.Name
	glider.UserID = int(gliderDB.UserID)
	if gliderDB.CreatedAt.Valid {
		glider.CreatedAt = gliderDB.CreatedAt.Time
	}
	if gliderDB.UpdatedAt.Valid {
		glider.UpdatedAt = gliderDB.UpdatedAt.Time
	}

	return glider
}

func (repo Repository) GetGliders(ctx context.Context, user user.User) ([]Glider, error) {
	glidersDB, err := repo.Queries.GetGliders(ctx, int32(user.ID))
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to get gliders")
	}

	var gliders []Glider
	for _, g := range glidersDB {
		gliders = append(gliders, Glider{
			ID:        int(g.ID),
			Name:      g.Name,
			UserID:    int(g.UserID),
			CreatedAt: g.CreatedAt.Time,
			UpdatedAt: g.UpdatedAt.Time,
		})
	}
	return gliders, err
}

func (repo Repository) AddGlider(ctx context.Context, gliderName string, user user.User) error {
	arg := storage.InsertGliderParams{
		Name:   gliderName,
		UserID: int32(user.ID),
	}
	err := repo.Queries.InsertGlider(ctx, arg)
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to get gliders")
	}
	return err
}
