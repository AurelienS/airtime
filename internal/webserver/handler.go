package webserver

import (
	"context"
	"strconv"

	"github.com/AurelienS/cigare/internal/storage/sqlc"
	"github.com/AurelienS/cigare/web/template/page"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Queries *sqlc.Queries
}

func (h *Handler) GetIndex(c echo.Context) error {
	flights, _ := h.Queries.GetFlights(context.Background())
	test := strconv.Itoa(len(flights))
	return render(c, page.Index(test))
}

func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}
