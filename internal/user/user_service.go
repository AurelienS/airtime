package user

import (
	"context"
	"strconv"

	"github.com/AurelienS/cigare/internal/storage"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return UserService{
		repo: repo,
	}
}

func (r *UserService) UpdateDefaultGlider(ctx context.Context, defaultGliderId string, user storage.User) error {
	gliderId, err := strconv.Atoi(defaultGliderId)
	if err != nil {
		return err
	}

	return r.repo.UpdateDefaultGlider(ctx, int32(gliderId), user.ID)
}

func (r UserService) UpsertUser(ctx context.Context, user storage.User) (storage.User, error) {
	return r.repo.UpsertUser(ctx, user)
}
