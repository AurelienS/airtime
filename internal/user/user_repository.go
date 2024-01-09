package user

import (
	"context"

	"github.com/AurelienS/cigare/internal/storage"
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

func (r *UserRepository) UpdateDefaultGlider(ctx context.Context, defaultGliderId int32, userId int32) error {
	arg := storage.UpdateDefaultGliderParams{
		DefaultGliderID: pgtype.Int4{Int32: defaultGliderId, Valid: true},
		ID:              userId,
	}

	return r.queries.UpdateDefaultGlider(context.Background(), arg)
}
