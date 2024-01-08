package user

import (
	"context"
	"database/sql"

	"github.com/AurelienS/cigare/internal/storage"
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
		DefaultGliderID: sql.NullInt32{Int32: defaultGliderId, Valid: true},
		ID:              userId,
	}

	return r.queries.UpdateDefaultGlider(context.Background(), arg)
}
