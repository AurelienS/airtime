package user

import (
	"context"

	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository struct {
	queries storage.Queries
}

func NewUserRepository(queries storage.Queries) UserRepository {
	return UserRepository{
		queries: queries,
	}
}

func ConvertUserDBToUser(userDB storage.User) User {
	var user User

	user.ID = int(userDB.ID)
	user.GoogleID = userDB.GoogleID
	user.Email = userDB.Email
	user.Name = userDB.Name
	user.PictureUrl = userDB.PictureUrl

	if userDB.DefaultGliderID.Valid {
		user.DefaultGliderID = int(userDB.DefaultGliderID.Int32)
	}

	if userDB.CreatedAt.Valid {
		user.CreatedAt = userDB.CreatedAt.Time
	}
	if userDB.UpdatedAt.Valid {
		user.UpdatedAt = userDB.UpdatedAt.Time
	}

	return user
}

func (r *UserRepository) UpsertUser(ctx context.Context, user User) (User, error) {
	param := storage.UpsertUserParams{
		GoogleID:        user.GoogleID,
		Email:           user.Email,
		Name:            user.Name,
		PictureUrl:      user.PictureUrl,
		DefaultGliderID: pgtype.Int4{Int32: int32(user.DefaultGliderID), Valid: user.DefaultGliderID > 0},
	}
	updatedUser, err := r.queries.UpsertUser(ctx, param)
	if err != nil {
		util.Error().Msgf("Failed to upsert user %v", param)
	}
	return ConvertUserDBToUser(updatedUser), err
}
