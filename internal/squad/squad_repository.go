package squad

import (
	"context"

	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/internal/storage/ent"
	userDb "github.com/AurelienS/cigare/internal/storage/ent/user"
	"github.com/AurelienS/cigare/internal/util"
)

type Repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) Repository {
	return Repository{
		client: client,
	}
}

func (r *Repository) InsertSquad(ctx context.Context, name string) (int, error) {
	squad, err := r.client.Squad.
		Create().
		SetName(name).
		Save(ctx)
	if err != nil {
		util.Error().Str("squadName", name).Msg("Failed to insert squad")
		return 0, err
	}
	return squad.ID, nil
}

func (r *Repository) InsertSquadMember(ctx context.Context, squadID, userID int) error {
	_, err := r.client.Squad.
		UpdateOneID(squadID).
		AddMemberIDs(userID).
		Save(ctx)
	if err != nil {
		util.Error().Int("squadID", squadID).Int("userID", userID).Msg("Failed to insert squad member")
		return err
	}

	return nil
}

func (r *Repository) Squads(ctx context.Context, user model.User) ([]model.Squad, error) {
	squads, err := r.client.User.
		Query().
		Where(userDb.IDEQ(user.ID)).
		QuerySquads().
		WithMembers().
		All(ctx)
	if err != nil {
		util.Error().Str("user", user.Email).Msg("Failed to find squads for user")
		return nil, err
	}

	return model.DBToDomainSquads(squads), nil
}
