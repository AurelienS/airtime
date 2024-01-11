package user

import (
	"context"

	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/util"
)

type Repository struct {
	queries storage.Queries
}

func NewRepository(queries storage.Queries) Repository {
	return Repository{
		queries: queries,
	}
}

func ConvertUserDBToUser(userDB storage.User) User {
	var user User

	user.ID = int(userDB.ID)
	user.GoogleID = userDB.GoogleID
	user.Email = userDB.Email
	user.Name = userDB.Name
	user.PictureURL = userDB.PictureUrl

	if userDB.CreatedAt.Valid {
		user.CreatedAt = userDB.CreatedAt.Time
	}

	return user
}

func (r *Repository) UpsertUser(ctx context.Context, user User) (User, error) {
	param := storage.UpsertUserParams{
		GoogleID:   user.GoogleID,
		Email:      user.Email,
		Name:       user.Name,
		PictureUrl: user.PictureURL,
	}
	updatedUser, err := r.queries.UpsertUser(ctx, param)
	if err != nil {
		util.Error().Msgf("Failed to upsert user %v", param)
	}
	return ConvertUserDBToUser(updatedUser), err
}
