package webserver

import (
	"net/http"

	"github.com/AurelienS/cigare/internal/storage/sqlc"
	"github.com/AurelienS/cigare/web/template/page"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Queries *sqlc.Queries
}

func (h *Handler) GetIndex(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/flights")
}

func (h *Handler) GetGliders(c echo.Context) error {
	return render(c, page.Gliders())
}

func (h *Handler) GetFlights(c echo.Context) error {
	return render(c, page.Flights())
}

func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}
