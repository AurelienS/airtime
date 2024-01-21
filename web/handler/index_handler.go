package handler

import (
	"fmt"
	"net/http"

	"github.com/AurelienS/cigare/internal/service"
	"github.com/AurelienS/cigare/web/session"
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
	return h.redirectToLogbook(c)
}

func (h IndexHandler) Dummy(c echo.Context) error {
	return Render(c, view.Dummy())
}

func (h IndexHandler) Landing(c echo.Context) error {
	return Render(c, view.Landing())
}

func (h IndexHandler) redirectToLogbook(c echo.Context) error {
	user := session.GetUserFromContext(c)
	lastFlight, err := h.flightService.GetLastFlight(c.Request().Context(), user)
	if err != nil {
		return err
	}
	if lastFlight == nil {
		return c.Redirect(http.StatusFound, "/onboarding")
	}
	lastYear := lastFlight.Date.Year()
	redirectTo := fmt.Sprintf("/logbook/%d", lastYear)
	return c.Redirect(http.StatusFound, redirectTo)
}
