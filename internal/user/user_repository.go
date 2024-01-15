package user

import (
	"context"

	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/internal/storage/ent"
	"github.com/AurelienS/cigare/internal/storage/ent/user"
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

func (r *Repository) InsertUser(ctx context.Context, user model.User) (model.User, error) {
	u, err := r.client.User.
		Create().
		SetGoogleID(user.GoogleID).
		SetEmail(user.Email).
		SetName(user.Name).
		SetPictureURL(user.PictureURL).
		Save(ctx)
	if err != nil {
		util.Error().Msgf("Failed to insert user %v", user)
		return user, err
	}
	util.Info().Str("user", user.Email).Msg("Inserted user")
	return model.DBToDomainUser(u), nil
}

func (r *Repository) UpdateUser(ctx context.Context, incomingUser model.User) (model.User, error) {
	u, err := r.client.User.
		Query().
		Where(user.GoogleIDEQ(incomingUser.GoogleID)).
		Only(ctx)
	if err != nil {
		return incomingUser, err
	}

	// Then, update the user
	u, err = r.client.User.
		UpdateOneID(u.ID).
		SetEmail(incomingUser.Email).
		SetName(incomingUser.Name).
		SetPictureURL(incomingUser.PictureURL).
		Save(ctx)
	if err != nil {
		util.Error().Msgf("Failed to update user %v", incomingUser)
		return incomingUser, err
	}
	util.Info().Str("user", u.Email).Msg("Updated user")
	return model.DBToDomainUser(u), nil
}

func (r *Repository) UserExists(ctx context.Context, googleID string) bool {
	exists, err := r.client.User.
		Query().
		Where(user.GoogleIDEQ(googleID)).
		Exist(ctx)
	if err != nil {
		util.Error().Msgf("Failed to check if user exists %v", googleID)
		return false
	}
	return exists
}
