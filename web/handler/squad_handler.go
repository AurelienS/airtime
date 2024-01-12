package handler

import (
	"net/http"

	"github.com/AurelienS/cigare/internal/squad"
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/template/page"
	"github.com/labstack/echo/v4"
)

type SquadHandler struct {
	squadService squad.Service
}

func NewSquadHandler(squadService squad.Service) SquadHandler {
	return SquadHandler{
		squadService: squadService,
	}
}

func (h SquadHandler) GetCreateSquad(c echo.Context) error {
	return Render(c, page.CreateSquad())
}

func (h SquadHandler) PostSquad(c echo.Context) error {
	squadName := c.FormValue("squad_name")
	user := session.GetUserFromContext(c)

	err := h.squadService.CreateSquad(c.Request().Context(), squadName, user)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}
