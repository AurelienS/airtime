package transformer

import (
	"github.com/AurelienS/airtime/internal/domain"
	"github.com/AurelienS/airtime/web/viewmodel"
)

func TransformUserToViewModel(user domain.User) viewmodel.UserView {
	return viewmodel.UserView{
		Name:       user.Name,
		PictureURL: user.PictureURL,
		Theme:      user.Theme,
	}
}
