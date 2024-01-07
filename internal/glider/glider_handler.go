package glider

import (
	"context"

	"github.com/AurelienS/cigare/internal/auth"
	"github.com/AurelienS/cigare/internal/log"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web/template/flight"
	"github.com/labstack/echo/v4"
)

type GliderHandler struct {
	GliderService GliderService
}

func NewGliderHandler(flightService *GliderService) *GliderHandler {
	return &GliderHandler{
		GliderService: *flightService,
	}
}

func (h *GliderHandler) PostGlider(c echo.Context) error {
	user := auth.GetUserFromContext(c)
	gliderName := c.FormValue("gliderName")

	err := h.GliderService.AddGlider(context.Background(), gliderName, user)
	if err != nil {
		return util.HandleError(c, err)
	}

	return h.GetGlidersCard(c)
}

func (h *GliderHandler) GetGlidersCard(c echo.Context) error {
	user := auth.GetUserFromContext(c)
	gliders, err := h.GliderService.GetGliders(context.Background(), user)
	if err != nil {
		return util.HandleError(c, err)
	}
	log.Info().Str("user", user.Email).Msg("Fetched gliders successfully")
	return util.Render(c, flight.GliderCard(gliders))
}
