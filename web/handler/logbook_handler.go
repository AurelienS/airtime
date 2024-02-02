package handler

import (
	"strconv"
	"time"

	"github.com/AurelienS/cigare/internal/service"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/transformer"
	"github.com/AurelienS/cigare/web/view/logbook"
	"github.com/AurelienS/cigare/web/viewmodel"
	"github.com/labstack/echo/v4"
)

type LogbookHandler struct {
	logbookService   service.LogbookService
	statisticService service.StatisticService
	flightService    service.FlightService
}

func NewLogbookHandler(
	logbookService service.LogbookService,
	statisticService service.StatisticService,
	flightService service.FlightService,
) LogbookHandler {
	return LogbookHandler{
		logbookService:   logbookService,
		statisticService: statisticService,
		flightService:    flightService,
	}
}

func (h *LogbookHandler) GetLogbook(c echo.Context) error {
	ctx := c.Request().Context()
	user := session.GetUserFromContext(c)

	flyingYears, err := h.statisticService.GetFlyingYears(ctx, user)
	if err != nil {
		return err
	}

	yearParam := c.Param("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		if yearParam == "" && len(flyingYears) > 0 {
			year = flyingYears[len(flyingYears)-1]
		} else {
			return err
		}
	}
	startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	endOfYear := time.Date(year, time.December, 31, 23, 59, 0, 0, time.UTC)
	yearFlights, err := h.flightService.GetFlights(
		ctx,
		startOfYear,
		endOfYear,
		user,
	)
	if err != nil {
		return err
	}

	viewData := transformer.TransformLogbookToViewModel(
		yearFlights,
		flyingYears,
		year,
	)
	userview := transformer.TransformUserToViewModel(user)

	return Render(c, logbook.Index(viewData, userview))
}

func (h *LogbookHandler) GetFlight(c echo.Context) error {
	user := session.GetUserFromContext(c)

	flightIDParam := c.Param("flight")
	flightID, err := strconv.Atoi(flightIDParam)
	if err != nil {
		return err
	}
	flight, err := h.flightService.GetFlight(c.Request().Context(), flightID, user)
	if err != nil {
		return err
	}

	geoJSON, err := flight.GenerateGeoJSON()
	if err != nil {
		return err
	}
	view := viewmodel.FlightDetailView{
		UserView:      transformer.TransformUserToViewModel(user),
		FlightView:    transformer.TransformFlightToViewmodel(flight),
		FlightGeoJSON: geoJSON,
	}
	return Render(c, logbook.Flight(view))
}

func (h *LogbookHandler) PostFlight(c echo.Context) error {
	file, err := c.FormFile("igcfile")
	if err != nil {
		util.Error().Err(err).Msg("Failed to get IGC file from form")
		return err
	}

	user := session.GetUserFromContext(c)

	err = h.logbookService.AddIGCFlight(c.Request().Context(), file, user)
	if err != nil {
		return err
	}

	c.Set("flight_added", "Flight processed and added successfully")
	c.Response().Header().Set("HX-Redirect", "/")
	return nil
}

func (h *LogbookHandler) DeleteFlight(c echo.Context) error {
	user := session.GetUserFromContext(c)

	id := c.Param("id")
	flightID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	err = h.flightService.RemoveFlight(c.Request().Context(), flightID, user)
	if err != nil {
		return err
	}
	return nil
}
