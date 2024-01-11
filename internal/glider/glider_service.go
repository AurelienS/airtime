package glider

import (
	"context"

	"github.com/AurelienS/cigare/internal/user"
)

type GliderService struct {
	repo GliderRepository
}

func NewGliderService(repository GliderRepository) GliderService {
	return GliderService{
		repo: repository,
	}
}

func (g *GliderService) GetGliders(ctx context.Context, user user.User) ([]Glider, error) {
	return g.repo.GetGliders(ctx, user)

}

func (g *GliderService) AddGlider(ctx context.Context, gliderName string, user user.User) error {
	return g.repo.AddGlider(ctx, gliderName, user)
}
