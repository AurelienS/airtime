package handler

import (
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/transformer"
	"github.com/AurelienS/cigare/web/view"
	"github.com/labstack/echo/v4"
)

type IndexHandler struct{}

func (h IndexHandler) Get(e echo.Context) error {
	user := session.GetUserFromContext(e)
	userview := transformer.TransformUserToViewModel(user)
	return Render(e, view.Index(userview))
}
