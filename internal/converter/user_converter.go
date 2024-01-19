package converter

import (
	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/internal/storage/ent"
)

func DBToDomainUser(userDB *ent.User) domain.User {
	return domain.User{
		ID:         userDB.ID,
		GoogleID:   userDB.GoogleID,
		Email:      userDB.Email,
		Name:       userDB.Name,
		PictureURL: userDB.PictureURL,
	}
}

func DBToDomainUsers(userDB []*ent.User) []domain.User {
	var users []domain.User
	for _, user := range userDB {
		users = append(users, DBToDomainUser(user))
	}
	return users
}
