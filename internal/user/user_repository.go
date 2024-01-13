package user

import (
	"context"

	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/internal/storage/ent"
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

func (r *Repository) UpsertUser(ctx context.Context, user model.User) (model.User, error) {
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
	return model.DBToDomainUser(u), nil
}
