package flight

import (
	"net/http"

	"github.com/AurelienS/cigare/internal/webserver/handler"
	"github.com/AurelienS/cigare/web/template/page"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (h *FlightHandler) GetIndexPage(c echo.Context) error {
	log.Info().Msg("Redirecting to flights page")
	return c.Redirect(http.StatusFound, "/flights")
}

func (h *FlightHandler) GetGlidersPage(c echo.Context) error {
	log.Info().Msg("Rendering gliders page")
	return handler.Render(c, page.Gliders())
}

func (h *FlightHandler) GetFlightsPage(c echo.Context) error {
	log.Info().Msg("Rendering flights page")
	return handler.Render(c, page.Flights())
}
