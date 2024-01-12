package squad

import (
	"context"
	"time"

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

func ConvertMemberDBToMember(memberDB storage.SquadMember) Member {
	return Member{
		ID:       int(memberDB.ID),
		SquadID:  int(memberDB.SquadID),
		UserID:   int(memberDB.UserID),
		Admin:    memberDB.Admin.Bool,
		JoinedAt: memberDB.JoinedAt.Time,
	}
}

func (r Repository) InsertSquad(ctx context.Context, name string) (int, error) {
	id, err := r.queries.InsertSquad(ctx, name)
	if err != nil {
		util.Error().Str("squadName", name).Msg("Failed to insert squad")
		return 0, err
	}
	return int(id), nil
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
	squadAndMembers, err := r.queries.FindAllSquadForUser(ctx, int32(user.ID))
	if err != nil {
		util.Error().Str("user", user.Email).Msg("Failed to find squads for user")
		return squads, err
	}

	squads = assembleSquadsFromRows(squadAndMembers)
	return squads, nil
}

func GroupBy[K comparable, T any](slice []T, keySelector func(T) K) map[K][]T {
	grouped := make(map[K][]T)
	for _, item := range slice {
		key := keySelector(item)
		grouped[key] = append(grouped[key], item)
	}
	return grouped
}

func assembleSquadsFromRows(rows []storage.FindAllSquadForUserRow) []Squad {
	squadsMap := GroupBy[int, storage.FindAllSquadForUserRow](
		rows,
		func(row storage.FindAllSquadForUserRow) int { return int(row.Squad.ID) },
	)

	var assembledSquads []Squad
	for squadID, squadRows := range squadsMap {
		var squadMembers []Member
		var squadName string
		var squadCreationTime time.Time

		// Assuming that all rows for the same squad will have identical squad info fields
		if len(squadRows) > 0 {
			squadInfo := squadRows[0]
			squadName = squadInfo.Squad.Name
			squadCreationTime = squadInfo.Squad.CreatedAt.Time
		}

		for _, squadMemberRow := range squadRows {
			member := ConvertMemberDBToMember(squadMemberRow.SquadMember)
			squadMembers = append(squadMembers, member)
		}

		squad := Squad{
			ID:        squadID,
			Name:      squadName,
			Members:   squadMembers,
			CreatedAt: squadCreationTime,
		}
		assembledSquads = append(assembledSquads, squad)
	}

	return assembledSquads
}
