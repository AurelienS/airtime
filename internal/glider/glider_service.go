package glider

import (
	"context"

	"github.com/AurelienS/cigare/internal/storage"
)

type GliderService struct {
	repo GliderRepository
}

func NewGliderService(repository GliderRepository) *GliderService {
	return &GliderService{
		repo: repository,
	}
}

func (g *GliderService) GetGliders(ctx context.Context, user storage.User) ([]storage.Glider, error) {
	return g.repo.GetGliders(ctx, user)
}

func (g *GliderService) AddGlider(ctx context.Context, gliderName string, user storage.User) error {
	return g.repo.AddGlider(ctx, gliderName, user)
}
