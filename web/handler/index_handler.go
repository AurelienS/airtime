package handler

import (
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/template/page"
	"github.com/labstack/echo/v4"
)

type IndexHandler struct{}

func (h IndexHandler) Get(e echo.Context) error {
	user := session.GetUserFromContext(e)
	return Render(e, page.Index(user))
}
