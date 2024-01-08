package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/AurelienS/cigare/internal/log"
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
		return fmt.Errorf("could not convert defaultGliderId param to int: %v", err)
	}

	err = r.repo.UpdateDefaultGlider(context.Background(), int32(gliderId), user.ID)
	if err != nil {
		log.Error().Err(err).Str("user", user.Email).Msg("Failed to update default_glider_id")
	}
	return err
}
