package glider

import (
	"context"

	"github.com/AurelienS/cigare/internal/user"
)

type Service struct {
	repo Repository
}

func NewService(repository Repository) Service {
	return Service{
		repo: repository,
	}
}

func (g *Service) GetGliders(ctx context.Context, user user.User) ([]Glider, error) {
	return g.repo.GetGliders(ctx, user)

}

func (g *Service) AddGlider(ctx context.Context, gliderName string, user user.User) error {
	return g.repo.AddGlider(ctx, gliderName, user)
}
