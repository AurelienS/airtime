package user

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (r *Service) UpdateDefaultGlider(ctx context.Context, defaultGliderID int, user User) error {
	user.DefaultGliderID = defaultGliderID
	_, err := r.repo.UpsertUser(ctx, user)
	return err
}

func (r Service) UpsertUser(ctx context.Context, user User) (User, error) {
	return r.repo.UpsertUser(ctx, user)
}
