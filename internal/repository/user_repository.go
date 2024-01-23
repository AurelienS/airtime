package repository

import (
	"context"

	"github.com/AurelienS/cigare/internal/converter"
	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/internal/storage/ent"
	"github.com/AurelienS/cigare/internal/storage/ent/user"
	"github.com/AurelienS/cigare/internal/util"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) UserRepository {
	return UserRepository{
		client: client,
	}
}

func (r *UserRepository) InsertUser(ctx context.Context, user domain.User) (domain.User, error) {
	u, err := r.client.User.
		Create().
		SetGoogleID(user.GoogleID).
		SetEmail(user.Email).
		SetName(user.Name).
		SetPictureURL(user.PictureURL).
		SetTheme(user.Theme).
		Save(ctx)
	if err != nil {
		util.Error().Msgf("Failed to insert user %v", user)
		return user, err
	}
	util.Info().Str("user", user.Email).Msg("Inserted user")
	return converter.DBToDomainUser(u), nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, incomingUser domain.User) (domain.User, error) {
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
		SetTheme(incomingUser.Theme).
		Save(ctx)
	if err != nil {
		util.Error().Msgf("Failed to update user %v", incomingUser)
		return incomingUser, err
	}
	util.Info().Str("user", u.Email).Msg("Updated user")
	return converter.DBToDomainUser(u), nil
}

func (r *UserRepository) UserExists(ctx context.Context, googleID string) bool {
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
