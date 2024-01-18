package handler

import (
	"net/http"

	"github.com/AurelienS/cigare/web/view"
	"github.com/labstack/echo/v4"
)

type IndexHandler struct{}

func (h IndexHandler) Get(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/logbook/log/0")
}

func (h IndexHandler) Dummy(c echo.Context) error {
	return Render(c, view.Dummy())
}

func (h IndexHandler) Landing(c echo.Context) error {
	return Render(c, view.Landing())
}
