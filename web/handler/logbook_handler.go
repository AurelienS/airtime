package handler

import (
	"fmt"
	"strconv"

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
	userview := transformer.TransformUserToViewModel(user)

	flyingYears, err := h.statisticService.GetFlyingYears(ctx, user)
	if err != nil {
		return err
	}

	year, err := h.getRequestedYear(c, flyingYears)
	if err != nil {
		return err
	}

	flights, yearStats, allTimeStats, err := h.statisticService.GetFlightStats(ctx, user, year)
	if err != nil {
		return Render(c, logbook.Index(viewmodel.LogbookView{}, userview))
	}

	viewData := transformer.TransformLogbookToViewModel(
		&flights,
		yearStats,
		allTimeStats,
		year,
		flyingYears,
		c.Get("flight_added") != nil,
	)

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

	view := viewmodel.FlightDetailView{
		UserView:   transformer.TransformUserToViewModel(user),
		FlightView: transformer.TransformFlightToViewModel(flight),
		Stats:      transformer.TransformStatToViewModel(flight.Statistic),
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
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to process and insert flight")
		return err
	}

	c.Set("flight_added", "Flight processed and added successfully")
	c.Response().Header().Set("HX-Redirect", "/")
	return nil
}

func (h *LogbookHandler) getRequestedYear(c echo.Context, flyingYears []int) (int, error) {
	yearParam := c.Param("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil || !yearInSlice(year, flyingYears) {
		return 0, fmt.Errorf("invalid year: %s", yearParam) // default to the last year if not specified or invalid
	}
	return year, nil
}

func yearInSlice(year int, slice []int) bool {
	for _, y := range slice {
		if y == year {
			return true
		}
	}
	return false
}
