package flight

import (
	"context"

	"github.com/AurelienS/cigare/internal/webserver/handler"
	"github.com/AurelienS/cigare/internal/webserver/middleware"
	"github.com/AurelienS/cigare/web/template/flight"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (h *FlightHandler) GetFlights(c echo.Context) error {
	user := middleware.GetUserFromContext(c)
	flights, err := h.FlightService.GetFlights(context.Background(), user)
	if err != nil {
		return handler.HandleError(c, err)
	}
	log.Info().Str("user", user.Email).Msg("Fetched flights successfully")
	return handler.Render(c, flight.FlightRecords(flights))
}

func (h *FlightHandler) GetGlidersCard(c echo.Context) error {
	user := middleware.GetUserFromContext(c)
	gliders, err := h.GliderService.GetGliders(context.Background(), user)
	if err != nil {
		return handler.HandleError(c, err)
	}
	log.Info().Str("user", user.Email).Msg("Fetched gliders successfully")
	return handler.Render(c, flight.GliderCard(gliders))
}
