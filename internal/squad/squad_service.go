package squad

import (
	"context"

	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/user"
)

type Service struct {
	squadRepo Repository
	tm        storage.TransactionManager
}

func NewService(squadRepo Repository, tm storage.TransactionManager) Service {
	return Service{
		squadRepo: squadRepo,
		tm:        tm,
	}
}

func (s Service) CreateSquad(ctx context.Context, name string, user user.User) error {
	transaction := func() error {
		squadID, err := s.squadRepo.InsertSquad(ctx, name)
		if err != nil {
			return err
		}

		err = s.squadRepo.InsertSquadMember(ctx, squadID, user.ID, true)
		if err != nil {
			return err
		}
		return nil
	}

	return s.tm.ExecuteTransaction(ctx, transaction)
}

func (s Service) UserSquads(ctx context.Context, user user.User) ([]Squad, error) {
	return s.squadRepo.Squads(ctx, user)
}
