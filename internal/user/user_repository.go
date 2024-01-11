package user

import (
	"context"

	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository struct {
	queries storage.Queries
}

func NewUserRepository(queries storage.Queries) UserRepository {
	return UserRepository{
		queries: queries,
	}
}

func (r *UserRepository) UpsertUser(ctx context.Context, user storage.User) (storage.User, error) {
	param := storage.UpsertUserParams{
		GoogleID:   user.GoogleID,
		Email:      user.Email,
		Name:       user.Name,
		PictureUrl: user.PictureUrl,
	}
	updatedUser, err := r.queries.UpsertUser(ctx, param)
	if err != nil {
		util.Error().Msgf("Failed to upsert user %v", param)
	}
	return updatedUser, err
}

func (r *UserRepository) UpdateDefaultGlider(ctx context.Context, defaultGliderId int32, userId int32) error {
	arg := storage.UpdateDefaultGliderParams{
		DefaultGliderID: pgtype.Int4{Int32: defaultGliderId, Valid: true},
		ID:              userId,
	}

	err := r.queries.UpdateDefaultGlider(ctx, arg)
	if err != nil {
		util.Error().Msgf("Failed to update default glider %v", arg)
	}

	return err
}
