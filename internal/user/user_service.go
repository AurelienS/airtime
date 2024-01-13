package user

import (
	"context"
	"fmt"

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

func (s Service) UpsertUser(ctx context.Context, userModel model.User) (model.User, error) {
	fmt.Println("file: user_service.go ~ line 20 ~ func ~ UpsertUser : ")
	exists := s.repo.UserExists(ctx, userModel.GoogleID)
	if exists {
		return s.repo.UpdateUser(ctx, userModel)
	}
	return s.repo.InsertUser(ctx, userModel)
}
