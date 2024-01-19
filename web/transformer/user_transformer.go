package transformer

import (
	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/web/viewmodel"
)

func TransformUserToViewModel(user domain.User) viewmodel.UserView {
	return viewmodel.UserView{
		Name:       user.Name,
		PictureURL: user.PictureURL,
	}
}
