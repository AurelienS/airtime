package service

import (
	"context"

	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return UserService{
		repo: repo,
	}
}

func (s UserService) UpsertUser(ctx context.Context, userModel domain.User) (domain.User, error) {
	exists := s.repo.UserExists(ctx, userModel.GoogleID)
	if exists {
		return s.repo.UpdateUser(ctx, userModel)
	}
	return s.repo.InsertUser(ctx, userModel)
}
