package handler

import (
	"github.com/AurelienS/cigare/internal/squad"
	"github.com/AurelienS/cigare/internal/user"
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/template/page"
	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	userService  user.Service
	squadService squad.Service
}

func NewDashboardHandler(userService user.Service, squadService squad.Service) DashboardHandler {
	return DashboardHandler{
		userService:  userService,
		squadService: squadService,
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
	}
	return Render(c, page.Dashboard(viewbag))
}
