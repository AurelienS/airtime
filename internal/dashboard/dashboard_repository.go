package user

import (
	"github.com/AurelienS/cigare/internal/storage"
)

type Repository struct {
	queries storage.Queries
}

func NewRepository(queries storage.Queries) Repository {
	return Repository{
		queries: queries,
	}
}
