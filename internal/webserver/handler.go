package webserver

import (
	"github.com/AurelienS/cigare/web/template"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Handler struct {
}

func (h *Handler) GetIndex(c echo.Context) error {
	return render(c, template.Index())
}

func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}
