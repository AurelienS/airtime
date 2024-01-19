package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/AurelienS/cigare/internal/logbook"
	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/transformer"
	"github.com/AurelienS/cigare/web/view/logbookview"
	"github.com/AurelienS/cigare/web/viewmodel"
	"github.com/labstack/echo/v4"
)

type LogbookHandler struct {
	LogbookService logbook.Service
}

func NewLogbookHandler(logbookService logbook.Service) LogbookHandler {
	return LogbookHandler{
		LogbookService: logbookService,
	}
}

func (h *LogbookHandler) GetTabLog(c echo.Context) error {
	ctx := c.Request().Context()
	user := session.GetUserFromContext(c)
	userview := transformer.TransformUserToViewModel(user)
	yearParam := c.Param("year")
	isFlightAdded := c.Get("flight_added") != nil

	flyingYears, err := h.LogbookService.GetFlyingYears(ctx, user)
	if err != nil {
		return err
	}

	numberOfYearFlying := len(flyingYears)
	if numberOfYearFlying == 0 {
		return Render(c, logbookview.TabLog(viewmodel.LogbookView{}, userview))
	}

	if yearParam == "" {
		if numberOfYearFlying == 1 {
			yearParam = strconv.Itoa(flyingYears[0])
		} else {
			yearParam = strconv.Itoa(flyingYears[numberOfYearFlying-1])
		}
	}

	year, err := strconv.Atoi(yearParam)
	if err != nil {
		return err
	}

	flyingYearIncludeYear := false
	for _, fy := range flyingYears {
		if fy == year {
			flyingYearIncludeYear = true
			break
		}
	}

	if !flyingYearIncludeYear && len(flyingYears) > 0 {
		lastYear := flyingYears[len(flyingYears)-1]
		redirectTo := fmt.Sprintf("/logbook/log/%d", lastYear)
		return c.Redirect(301, redirectTo)
	}

	flights, err := h.LogbookService.GetFlights(ctx, year, user)
	if err != nil {
		return err
	}

	allTimeStats, err := h.LogbookService.GetStatistics(
		ctx,
		time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Now(),
		user,
	)
	if err != nil {
		return err
	}

	yearStats, err := h.LogbookService.GetStatistics(
		ctx,
		time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(year, time.December, 31, 23, 59, 59, 0, time.UTC),
		user,
	)
	if err != nil {
		return err
	}

	viewData := transformer.TransformLogbookToViewModel(flights, yearStats, allTimeStats, year, flyingYears, isFlightAdded)

	return Render(c, logbookview.TabLog(viewData, userview))
}

func (h *LogbookHandler) GetTabProgression(c echo.Context) error {
	user := session.GetUserFromContext(c)
	statsYearMonth, err := h.LogbookService.GetStatisticsByYearAndMonth(
		c.Request().Context(),
		user)
	if err != nil {
		return err
	}

	totalFlightTimeExtractor := func(stats model.StatsAggregated) int {
		return int(stats.TotalFlightTime.Hours())
	}
	flightCountExtractor := func(stats model.StatsAggregated) int {
		return stats.FlightCount
	}

	view := viewmodel.ProgressionView{
		User:                   transformer.TransformUserToViewModel(user),
		FlightTimeMonthlyData:  transformer.TransformStatsViewModel(statsYearMonth, totalFlightTimeExtractor),
		FlightCountMonthlyData: transformer.TransformStatsViewModel(statsYearMonth, flightCountExtractor),
	}
	return Render(c, logbookview.TabProgression(view))
}

func (h *LogbookHandler) GetFlight(c echo.Context) error {
	var f model.Flight
	return Render(c, logbookview.Flight(f))
}

func (h *LogbookHandler) PostFlight(c echo.Context) error {
	file, err := c.FormFile("igcfile")
	if err != nil {
		util.Error().Err(err).Msg("Failed to get IGC file from form")
		return err
	}

	user := session.GetUserFromContext(c)

	err = h.LogbookService.ProcessAndAddFlight(c.Request().Context(), file, user)
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to process and insert flight")
		return err
	}

	c.Set("flight_added", "Flight processed and added successfully")

	return h.GetTabLog(c)
}
