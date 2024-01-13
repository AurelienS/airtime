package squad

import (
	"context"

	"github.com/AurelienS/cigare/internal/model"
)

type Service struct {
	squadRepo Repository
}

func NewService(squadRepo Repository) Service {
	return Service{
		squadRepo: squadRepo,
	}
}

func (s Service) CreateSquad(ctx context.Context, name string, user model.User) error {
	squadID, err := s.squadRepo.InsertSquad(ctx, name)
	if err != nil {
		return err
	}

	err = s.squadRepo.InsertSquadMember(ctx, squadID, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) UserSquads(ctx context.Context, user model.User) ([]model.Squad, error) {
	return s.squadRepo.Squads(ctx, user)
}
