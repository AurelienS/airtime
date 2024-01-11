package handler

import (
	"fmt"
	"io"
	"strconv"

	"github.com/AurelienS/cigare/internal/flight"
	"github.com/AurelienS/cigare/internal/glider"
	"github.com/AurelienS/cigare/internal/user"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web/session"
	flightView "github.com/AurelienS/cigare/web/template/flight"
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

	data, err := h.FlightService.GetDashboardData(c.Request().Context(), user)
	if err != nil {
		HandleError(c, err)
	}
	viewData := TransformDashboardToView(data, user)

	return Render(c, page.Flights(viewData))
}

func TransformDashboardToView(data flight.DashboardData, user user.User) flightView.DashboardView {
	var fv []flightView.FlightView
	for _, f := range data.Flights {
		fv = append(fv, flightView.FlightView{
			TakeoffLocation: f.TakeoffLocation,
			Date:            f.Date.Format("02/01 15h04"),
		})
	}

	return flightView.DashboardView{
		Flights:         fv,
		Gliders:         TransformGlidersToView(data.Gliders, user),
		NumberOfFlight:  strconv.Itoa(len(data.Flights)),
		TotalFlightTime: fmt.Sprintf("%d", int(data.TotalFlightTime.Hours())),
	}

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
