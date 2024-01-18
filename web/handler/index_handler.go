package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type IndexHandler struct{}

func (h IndexHandler) Get(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/logbook/log/0")
}
