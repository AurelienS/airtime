package handler

import (
	"net/http"

	"github.com/AurelienS/cigare/internal/storage/sqlc"
	"github.com/AurelienS/cigare/web/template/page"
	"github.com/labstack/echo/v4"
)

type FlightHandler struct {
	Queries *sqlc.Queries
}

func (h *FlightHandler) GetIndex(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/flights")
}

func (h *FlightHandler) GetGliders(c echo.Context) error {
	return Render(c, page.Gliders())
}

func (h *FlightHandler) GetFlights(c echo.Context) error {
	return Render(c, page.Flights())
}
