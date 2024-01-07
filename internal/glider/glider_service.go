package glider

import (
	"context"

	"github.com/AurelienS/cigare/internal/storage"
)

type GliderService struct {
	Repo GliderRepository
}

func (g *GliderService) GetGliders(ctx context.Context, user storage.User) ([]storage.Glider, error) {
	return g.Repo.GetGliders(ctx, user)
}
