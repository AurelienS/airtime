package user

import (
	"context"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return UserService{
		repo: repo,
	}
}

func (r *UserService) UpdateDefaultGlider(ctx context.Context, defaultGliderId int, user User) error {
	return r.repo.UpdateDefaultGlider(ctx, defaultGliderId, user.ID)
}

func (r UserService) UpsertUser(ctx context.Context, user User) (User, error) {
	return r.repo.UpsertUser(ctx, user)
}
