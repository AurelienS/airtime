package handler

import (
	"strconv"

	"github.com/AurelienS/cigare/internal/flight"
	"github.com/AurelienS/cigare/internal/squad"
	"github.com/AurelienS/cigare/internal/user"
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/template/page"
	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	userService   user.Service
	squadService  squad.Service
	flightService flight.Service
}

func NewDashboardHandler(userService user.Service,
	squadService squad.Service,
	flightService flight.Service,
) DashboardHandler {
	return DashboardHandler{
		userService:   userService,
		squadService:  squadService,
		flightService: flightService,
	}
}

func (h DashboardHandler) GetIndex(c echo.Context) error {
	user := session.GetUserFromContext(c)

	userSquads, err := h.squadService.UserSquads(c.Request().Context(), user)
	if err != nil {
		return err
	}

	isPartOfAtLeastOneSquad := len(userSquads) > 0
	viewbag := page.DashboardView{
		IsPartOfSquad: isPartOfAtLeastOneSquad,
		Squads:        userSquads,
		NumberOfSquad: strconv.Itoa(len(userSquads)),
	}
	return Render(c, page.Dashboard(viewbag))
}

func (h DashboardHandler) PostFlight(c echo.Context) error {
	err := insertFlight(c, h.flightService)
	if err != nil {
		return err
	}
	return h.GetIndex(c)
}
