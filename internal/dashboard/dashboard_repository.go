package dashboard

import (
	"github.com/AurelienS/cigare/internal/storage/ent"
)

type Repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) Repository {
	return Repository{
		client: client,
	}
}
