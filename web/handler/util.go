package handler

import (
	"net/http"

	"github.com/AurelienS/airtime/internal/util"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

func HandleError(c echo.Context, err error) {
	util.Error().Msgf("Error encountered: %s\n", err)
	err = c.String(http.StatusInternalServerError, "Internal Server Error")
	if err != nil {
		util.Error().Msg("Cannot set status code 500 to response")
	}
}
