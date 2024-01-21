package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/internal/service"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/transformer"
	"github.com/AurelienS/cigare/web/view/logbookview"
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

	year := h.getRequestedYear(c, flyingYears)

	if !yearInSlice(year, flyingYears) {
		return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/logbook/%d", flyingYears[len(flyingYears)-1]))
	}

	flights, yearStats, allTimeStats, err := h.getFlightStats(ctx, user, year)
	if err != nil {
		return Render(c, logbookview.TabLog(viewmodel.LogbookView{}, userview))
	}

	viewData := transformer.TransformLogbookToViewModel(
		&flights,
		yearStats,
		allTimeStats,
		year,
		flyingYears,
		c.Get("flight_added") != nil,
	)

	return Render(c, logbookview.TabLog(viewData, userview))
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
	return Render(c, logbookview.Flight(view))
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

	return h.GetLogbook(c)
}

func (h *LogbookHandler) getRequestedYear(c echo.Context, flyingYears []int) int {
	yearParam := c.Param("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil || !yearInSlice(year, flyingYears) {
		return flyingYears[len(flyingYears)-1] // default to the last year if not specified or invalid
	}
	return year
}

func (h *LogbookHandler) getFlightStats(
	ctx context.Context,
	user domain.User,
	year int,
) ([]domain.Flight, service.StatsAggregated, service.StatsAggregated, error) {
	allFlights, err := h.flightService.GetFlights(
		ctx,
		time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Now(),
		user,
	)
	if err != nil {
		return nil, service.StatsAggregated{}, service.StatsAggregated{}, err
	}

	yearFlights := h.flightService.GetFlightsForYear(year, allFlights)
	yearStats := h.statisticService.ComputeAggregateStatistics(yearFlights)
	allTimeStats := h.statisticService.ComputeAggregateStatistics(allFlights)
	return yearFlights, yearStats, allTimeStats, nil
}

func yearInSlice(year int, slice []int) bool {
	for _, y := range slice {
		if y == year {
			return true
		}
	}
	return false
}
