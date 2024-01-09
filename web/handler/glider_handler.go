package handler

import (
	"context"

	"github.com/AurelienS/cigare/internal/glider"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/template/flight"
	"github.com/labstack/echo/v4"
)

type GliderHandler struct {
	GliderService glider.GliderService
}

func NewGliderHandler(gliderService glider.GliderService) GliderHandler {
	return GliderHandler{
		GliderService: gliderService,
	}
}

func (h *GliderHandler) PostGlider(c echo.Context) error {
	user := session.GetUserFromContext(c)
	gliderName := c.FormValue("gliderName")

	err := h.GliderService.AddGlider(context.Background(), gliderName, user)
	if err != nil {
		return HandleError(c, err)
	}

	return h.GetGlidersCard(c)
}

func (h *GliderHandler) GetGlidersCard(c echo.Context) error {
	user := session.GetUserFromContext(c)
	gliders, err := h.GliderService.GetGliders(context.Background(), user)
	if err != nil {
		return HandleError(c, err)
	}
	util.Info().Str("user", user.Email).Msg("Fetched gliders successfully")
	return Render(c, flight.GliderCard(gliders))
}
