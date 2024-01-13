package user

import (
	"context"

	"github.com/AurelienS/cigare/internal/model"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (r Service) UpsertUser(ctx context.Context, user model.User) (model.User, error) {
	return r.repo.UpsertUser(ctx, user)
}
