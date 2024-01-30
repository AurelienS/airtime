package handler

import (
	"time"

	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/internal/service"
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/transformer"
	"github.com/AurelienS/cigare/web/view/statistics"
	"github.com/AurelienS/cigare/web/view/statistics/chart"
	"github.com/labstack/echo/v4"
)

type StatisticsHandler struct {
	statisticService service.StatisticService
}

func NewStatisticsHandler(statisticService service.StatisticService) StatisticsHandler {
	return StatisticsHandler{statisticService: statisticService}
}

func (h *StatisticsHandler) GetIndex(c echo.Context) error {
	user := session.GetUserFromContext(c)

	return Render(c, statistics.Index(transformer.TransformUserToViewModel(user)))
}

func (h *StatisticsHandler) GetCountDistinct(c echo.Context) error {
	user := session.GetUserFromContext(c)
	monthlyStatsByYear, err := h.statisticService.GetMonthlyStatisticsByYear(
		c.Request().Context(),
		user)
	if err != nil {
		return err
	}

	flightCountExtractor := func(stats domain.MultipleFlightStats) int {
		return len(stats.Flights)
	}

	view := transformer.TransformMonthlyCountToViewmodel(monthlyStatsByYear, flightCountExtractor)
	return Render(c, chart.CountDistinct(view))
}

func (h *StatisticsHandler) GetCountCumul(c echo.Context) error {
	user := session.GetUserFromContext(c)
	start := time.Time{}
	end := time.Now()
	flightCounts, err := h.statisticService.GetCumulativeFlightCount(
		c.Request().Context(),
		user,
		start,
		end)
	if err != nil {
		return err
	}

	chartData := transformer.TransformChartCountCumulative(flightCounts)
	return Render(c, chart.CountCumul(chartData))
}

func (h *StatisticsHandler) GetTimeDistinct(c echo.Context) error {
	user := session.GetUserFromContext(c)
	monthlyStatsByYear, err := h.statisticService.GetMonthlyStatisticsByYear(
		c.Request().Context(),
		user)
	if err != nil {
		return err
	}

	totalFlightTimeExtractor := func(stats domain.MultipleFlightStats) float64 {
		return stats.DurationTotal.Hours()
	}

	view := transformer.TransformMonthlyTimeToViewmodel(monthlyStatsByYear, totalFlightTimeExtractor)
	return Render(c, chart.TimeDistinct(view))
}

func (h *StatisticsHandler) GetTimeCumul(c echo.Context) error {
	user := session.GetUserFromContext(c)
	start := time.Time{}
	end := time.Now()
	flightDurations, err := h.statisticService.GetCumulativeFlightDuration(
		c.Request().Context(),
		user,
		start,
		end)
	if err != nil {
		return err
	}

	chartData := transformer.TransformChartTimeCumulative(flightDurations)
	return Render(c, chart.TimeCumul(chartData))
}
