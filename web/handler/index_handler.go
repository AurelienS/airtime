package handler

import (
	"github.com/AurelienS/cigare/web/template/page"
	"github.com/labstack/echo/v4"
)

type IndexHandler struct{}

func (h IndexHandler) Get(e echo.Context) error {
	return Render(e, page.Index())
}
