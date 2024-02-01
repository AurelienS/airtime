package handler

import (
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
	monthlyCountByYear, err := h.statisticService.ComputeStatistics(
		c.Request().Context(),
		user,
		[]service.StatisticType{service.MonthlyCount},
	)
	if err != nil {
		return err
	}

	view := transformer.TransformMonthlyCountToViewmodel(monthlyCountByYear.MonthlyCount)

	return Render(c, chart.CountDistinct(view))
}

func (h *StatisticsHandler) GetCountCumul(c echo.Context) error {
	user := session.GetUserFromContext(c)

	cumulativeMonthlyCount, err := h.statisticService.ComputeStatistics(
		c.Request().Context(),
		user,
		[]service.StatisticType{service.CumulativeMonthlyCount},
	)
	if err != nil {
		return err
	}

	chartData := transformer.TransformCumulativeCount(cumulativeMonthlyCount.CumulativeMonthlyCount)
	return Render(c, chart.CountCumul(chartData))
}

func (h *StatisticsHandler) GetTimeDistinct(c echo.Context) error {
	user := session.GetUserFromContext(c)
	monthlyDurationByYear, err := h.statisticService.ComputeStatistics(
		c.Request().Context(),
		user,
		[]service.StatisticType{
			service.MonthlyDuration,
		},
	)
	if err != nil {
		return err
	}

	view := transformer.TransformMonthlyTimeToViewmodel(monthlyDurationByYear.MonthlyDuration)
	return Render(c, chart.TimeDistinct(view))
}

func (h *StatisticsHandler) GetTimeCumul(c echo.Context) error {
	user := session.GetUserFromContext(c)

	cumulativeMonthlyDuration, err := h.statisticService.ComputeStatistics(
		c.Request().Context(),
		user,
		[]service.StatisticType{service.CumulativeMonthlyDuration},
	)
	if err != nil {
		return err
	}

	chartData := transformer.TransformChartTimeCumulative(cumulativeMonthlyDuration.CumulativeMonthlyDuration)
	return Render(c, chart.TimeCumul(chartData))
}
