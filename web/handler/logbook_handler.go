package handler

import (
	"fmt"
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
	LogbookService service.LogbookService
}

func NewLogbookHandler(logbookService service.LogbookService) LogbookHandler {
	return LogbookHandler{
		LogbookService: logbookService,
	}
}

func (h *LogbookHandler) GetTabLog(c echo.Context) error {
	ctx := c.Request().Context()
	user := session.GetUserFromContext(c)
	userview := transformer.TransformUserToViewModel(user)
	yearParam := c.Param("year")
	fmt.Println("file: logbook_handler.go ~ line 39 ~ func ~ c.ParamValues() : ", c.ParamValues())
	fmt.Println("file: logbook_handler.go ~ line 33 ~ func ~ yearParam : ", yearParam)
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
		fmt.Println(
			"file: logbook_handler.go ~ line 69 ~ if!flyingYearIncludeYear&&len ~ flyingYearIncludeYear : ",
			flyingYearIncludeYear,
		)
		lastYear := flyingYears[len(flyingYears)-1]
		redirectTo := fmt.Sprintf("/logbook/log/%d", lastYear)
		return c.Redirect(301, redirectTo)
	}

	allFlights, err := h.LogbookService.GetFlights(
		ctx,
		time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Now(),
		user,
	)
	if err != nil {
		return err
	}

	startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	endOfYear := time.Date(year, time.December, 31, 23, 59, 59, 0, time.UTC)

	dateRanges := []DateRange{
		{Start: startOfYear, End: endOfYear},
	}

	rangedFlights := getFlightsForDateRanges(allFlights, dateRanges)

	yearStats := domain.ComputeAggregateStatistics(rangedFlights[0])
	allTimeStats := domain.ComputeAggregateStatistics(allFlights)

	viewData := transformer.TransformLogbookToViewModel(
		&allFlights,
		yearStats,
		allTimeStats,
		year,
		flyingYears,
		isFlightAdded)

	return Render(c, logbookview.TabLog(viewData, userview))
}

type DateRange struct {
	Start time.Time
	End   time.Time
}

func getFlightsForDateRanges(flights []domain.Flight, dateRanges []DateRange) [][]domain.Flight {
	flightsForRanges := make([][]domain.Flight, len(dateRanges))

	for _, flight := range flights {
		for i, dateRange := range dateRanges {
			if (flight.Date.Equal(dateRange.Start) || flight.Date.After(dateRange.Start)) &&
				(flight.Date.Equal(dateRange.End) || flight.Date.Before(dateRange.End)) {
				flightsForRanges[i] = append(flightsForRanges[i], flight)
			}
		}
	}

	return flightsForRanges
}

func (h *LogbookHandler) GetTabProgression(c echo.Context) error {
	user := session.GetUserFromContext(c)
	statsYearMonth, err := h.LogbookService.GetStatisticsByYearAndMonth(
		c.Request().Context(),
		user)
	if err != nil {
		return err
	}

	totalFlightTimeExtractor := func(stats domain.StatsAggregated) int {
		return int(stats.TotalFlightTime.Hours())
	}
	flightCountExtractor := func(stats domain.StatsAggregated) int {
		return stats.FlightCount
	}

	view := viewmodel.ProgressionView{
		User:                   transformer.TransformUserToViewModel(user),
		FlightTimeMonthlyData:  transformer.TransformChartViewModel(statsYearMonth, totalFlightTimeExtractor),
		FlightCountMonthlyData: transformer.TransformChartViewModel(statsYearMonth, flightCountExtractor),
	}
	return Render(c, logbookview.TabProgression(view))
}

func (h *LogbookHandler) GetFlight(c echo.Context) error {
	user := session.GetUserFromContext(c)

	flightIDParam := c.Param("flight")
	flightID, err := strconv.Atoi(flightIDParam)
	if err != nil {
		return err
	}
	flight, err := h.LogbookService.GetFlight(c.Request().Context(), flightID, user)
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

	err = h.LogbookService.ProcessAndAddFlight(c.Request().Context(), file, user)
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to process and insert flight")
		return err
	}

	c.Set("flight_added", "Flight processed and added successfully")

	return h.GetTabLog(c)
}
