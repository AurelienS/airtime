package squad

import (
	"context"

	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/user"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository struct {
	queries storage.Queries
}

func NewRepository(queries storage.Queries) Repository {
	return Repository{
		queries: queries,
	}
}

func ConvertSquadDBToSquad(squadDB storage.Squad) Squad {
	return Squad{
		ID:        int(squadDB.ID),
		Name:      squadDB.Name,
		CreatedAt: squadDB.CreatedAt.Time,
	}
}

func (r Repository) InsertSquad(ctx context.Context, name string) (Squad, error) {
	squadDB, err := r.queries.InsertSquad(ctx, name)
	squad := ConvertSquadDBToSquad(squadDB)
	if err != nil {
		util.Error().Str("squadName", name).Msg("Failed to insert squad")
		return squad, err
	}
	return squad, nil
}

func (r Repository) InsertSquadMember(ctx context.Context, squadID, userID int, admin bool) error {
	arg := storage.InsertSquadMemberParams{
		SquadID: int32(squadID),
		UserID:  int32(userID),
		Admin:   pgtype.Bool{Bool: admin, Valid: true},
	}
	err := r.queries.InsertSquadMember(ctx, arg)
	if err != nil {
		util.Error().Int("squadID", squadID).Int("userID", userID).Msg("Failed to insert squad member")
		return err
	}
	return nil
}

func (r Repository) Squads(ctx context.Context, user user.User) ([]Squad, error) {
	var squads []Squad
	squadsDB, err := r.queries.FindAllSquadForUser(ctx, int32(user.ID))
	if err != nil {
		util.Error().Str("user", user.Email).Msg("Failed to find squads for user")
		return squads, err
	}
	for _, squadDB := range squadsDB {
		squads = append(squads, ConvertSquadDBToSquad(squadDB))
	}

	return squads, nil
}
