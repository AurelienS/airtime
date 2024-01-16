package transformer

import (
	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/web/viewmodel"
)

func TransformUserToViewModel(user model.User) viewmodel.UserView {
	return viewmodel.UserView{
		Name:       user.Name,
		PictureURL: user.PictureURL,
	}
}
