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

	squads = membersToSquads(squadAndMembers)
	return squads, nil
}

func membersToSquads(squadAndMembers []storage.FindAllSquadForUserRow) []Squad {
	membersBySquadID := make(map[int][]Member)
	squadInfoByID := make(map[int]storage.FindAllSquadForUserRow)

	for _, memberRow := range squadAndMembers {
		squadID := int(memberRow.Squad.ID)
		member := ConvertMemberDBToMember(memberRow.SquadMember)
		membersBySquadID[squadID] = append(membersBySquadID[squadID], member)
		if _, exists := squadInfoByID[squadID]; !exists {
			squadInfoByID[squadID] = memberRow
		}
	}

	squads := make([]Squad, 0, len(membersBySquadID))
	for squadID, members := range membersBySquadID {
		squadInfo := squadInfoByID[squadID]
		squad := Squad{
			ID:        squadID,
			Name:      squadInfo.Squad.Name,
			Members:   members,
			CreatedAt: squadInfo.Squad.CreatedAt.Time,
		}
		squads = append(squads, squad)
	}

	return squads
}
