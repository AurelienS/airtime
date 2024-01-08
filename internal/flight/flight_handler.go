package flight

import (
	"context"
	"io"
	"net/http"

	"github.com/AurelienS/cigare/internal/auth"
	"github.com/AurelienS/cigare/internal/glider"
	"github.com/AurelienS/cigare/internal/log"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web/template/page"
	"github.com/labstack/echo/v4"
)

type FlightHandler struct {
	FlightService FlightService
	GliderService glider.GliderService
}

func NewFlightHandler(flightService FlightService, GliderService glider.GliderService) FlightHandler {
	return FlightHandler{
		FlightService: flightService,
		GliderService: GliderService,
	}
}

func (h *FlightHandler) GetIndexPage(c echo.Context) error {
	user := auth.GetUserFromContext(c)
	flights, err := h.FlightService.GetFlights(context.Background(), user)
	if err != nil {
		return util.HandleError(c, err)
	}

	gliders, err := h.GliderService.GetGliders(context.Background(), user)
	if err != nil {
		return util.HandleError(c, err)
	}

	return util.Render(c, page.Flights(flights, gliders))
}

func (h *FlightHandler) PostFlight(c echo.Context) error {
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
