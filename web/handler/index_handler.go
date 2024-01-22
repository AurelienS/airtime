package handler

import (
	"net/http"

	"github.com/AurelienS/cigare/internal/service"
	"github.com/AurelienS/cigare/web/view"
	"github.com/labstack/echo/v4"
)

type IndexHandler struct {
	flightService service.FlightService
}

func NewIndexHandler(flightService service.FlightService) IndexHandler {
	return IndexHandler{
		flightService: flightService,
	}
}

func (h IndexHandler) Get(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/dashboard")
}

func (h IndexHandler) Dummy(c echo.Context) error {
	return Render(c, view.Dummy())
}

func (h IndexHandler) Landing(c echo.Context) error {
	return Render(c, view.Landing())
}
