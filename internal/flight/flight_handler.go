package flight

import (
	"context"
	"io"
	"net/http"

	"github.com/AurelienS/cigare/internal/auth"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web/template/flight"
	"github.com/AurelienS/cigare/web/template/page"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type FlightHandler struct {
	FlightService FlightService
}

func NewFlightHandler(flightService FlightService) *FlightHandler {
	return &FlightHandler{
		FlightService: flightService,
	}
}

/* **********************************
 *            PAGES
 ********************************** */

func (h *FlightHandler) GetIndexPage(c echo.Context) error {
	log.Info().Msg("Redirecting to flights page")
	return c.Redirect(http.StatusFound, "/flights")
}

func (h *FlightHandler) GetGlidersPage(c echo.Context) error {
	log.Info().Msg("Rendering gliders page")
	return util.Render(c, page.Gliders())
}

func (h *FlightHandler) GetFlightsPage(c echo.Context) error {
	log.Info().Msg("Rendering flights page")
	return util.Render(c, page.Flights())
}

/* **********************************
 *            DATA
 ********************************** */

func (h *FlightHandler) GetFlights(c echo.Context) error {
	user := auth.GetUserFromContext(c)
	flights, err := h.FlightService.GetFlights(context.Background(), user)
	if err != nil {
		return util.HandleError(c, err)
	}
	log.Info().Str("user", user.Email).Msg("Fetched flights successfully")
	return util.Render(c, flight.FlightRecords(flights))
}

func (h *FlightHandler) Upload(c echo.Context) error {
	file, err := c.FormFile("igcfile")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get IGC file from form")
		return util.HandleError(c, err)
	}

	src, err := file.Open()
	if err != nil {
		log.Error().Err(err).Str("filename", file.Filename).Msg("Failed to open IGC file")
		return util.HandleError(c, err)
	}
	defer src.Close()

	byteContent, err := io.ReadAll(src)
	if err != nil {
		log.Error().Err(err).Str("filename", file.Filename).Msg("Failed to read IGC file")
		return util.HandleError(c, err)
	}

	user := auth.GetUserFromContext(c)

	err = h.FlightService.UploadFlight(c.Request().Context(), byteContent, user)
	if err != nil {
		log.Error().Err(err).Str("user", user.Email).Msg("Failed to insert flight into database")
		return util.HandleError(c, err)
	}

	log.Info().Str("user", user.Email).Str("filename", file.Filename).Msg("File parsed and flight record created successfully")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "File parsed successfully",
	})
}
