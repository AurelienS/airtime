package service

import (
	"context"

	"github.com/AurelienS/cigare/internal/storage"
	repo "github.com/AurelienS/cigare/internal/storage/repository"
)

type GliderService struct {
	Repo repo.GliderRepository
}

func (g *GliderService) GetGliders(ctx context.Context, user storage.User) ([]storage.Glider, error) {
	return g.Repo.GetGliders(ctx, user)
}
