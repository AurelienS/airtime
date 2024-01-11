package handler

import (
	"fmt"
	"io"

	"github.com/AurelienS/cigare/internal/flight"
	"github.com/AurelienS/cigare/internal/glider"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/template/page"
	"github.com/labstack/echo/v4"
)

type FlightHandler struct {
	FlightService flight.FlightService
	GliderService glider.GliderService
}

func NewFlightHandler(flightService flight.FlightService, GliderService glider.GliderService) FlightHandler {
	return FlightHandler{
		FlightService: flightService,
		GliderService: GliderService,
	}
}

func (h *FlightHandler) GetIndexPage(c echo.Context) error {
	user := session.GetUserFromContext(c)
	context := c.Request().Context()
	flights, err := h.FlightService.GetFlights(context, user)
	if err != nil {
		return HandleError(c, err)
	}

	gliders, err := h.GliderService.GetGliders(context, user)
	if err != nil {
		return HandleError(c, err)
	}

	totalFlightTime, err := h.FlightService.GetTotalFlightTime(context, int(user.ID))
	fmt.Println("file: flight_handler.go ~ line 41 ~ err : ", err)
	fmt.Println("file: flight_handler.go ~ line 40 ~ totalFlightTime : ", totalFlightTime)

	return Render(c, page.Flights(flights, gliders, totalFlightTime))
}

func (h *FlightHandler) PostFlight(c echo.Context) error {
	file, err := c.FormFile("igcfile")
	if err != nil {
		util.Error().Err(err).Msg("Failed to get IGC file from form")
		return HandleError(c, err)
	}

	src, err := file.Open()
	if err != nil {
		util.Error().Err(err).Str("filename", file.Filename).Msg("Failed to open IGC file")
		return HandleError(c, err)
	}
	defer src.Close()

	byteContent, err := io.ReadAll(src)
	if err != nil {
		util.Error().Err(err).Str("filename", file.Filename).Msg("Failed to read IGC file")
		return HandleError(c, err)
	}

	user := session.GetUserFromContext(c)

	err = h.FlightService.AddFlight(c.Request().Context(), byteContent, user)
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to insert flight into database")
		return HandleError(c, err)
	}

	util.Info().Str("user", user.Email).Str("filename", file.Filename).Msg("File parsed and flight record created successfully")

	return h.GetIndexPage(c)
}
