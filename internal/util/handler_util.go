package util

import (
	"net/http"

	"github.com/AurelienS/cigare/internal/log"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

func HandleError(c echo.Context, err error) error {
	log.Error().Msgf("Error encountered: %s\n", err)
	return c.String(http.StatusInternalServerError, "Internal Server Error")
}
