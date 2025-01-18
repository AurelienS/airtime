package handler

import (
	"github.com/AurelienS/airtime/internal/domain"
	"github.com/AurelienS/airtime/internal/service"
	"github.com/AurelienS/airtime/web/session"
	"github.com/AurelienS/airtime/web/transformer"
	"github.com/AurelienS/airtime/web/view/statistics"
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

func (h *StatisticsHandler) GetCountByMonth(c echo.Context) error {
	user := session.GetUserFromContext(c)
	stats, err := h.statisticService.ComputeStatistics(
		c.Request().Context(),
		user,
	)
	if err != nil {
		return err
	}

	chartItems := domain.ConvertDateCountToChartDataItem(stats.MonthlyCount)
	view := transformer.TransformMultiDatasetsToViewmodel(chartItems)

	return Render(c, statistics.CountByMonth(view))
}

func (h *StatisticsHandler) GetCountByYear(c echo.Context) error {
	user := session.GetUserFromContext(c)
	stats, err := h.statisticService.ComputeStatistics(
		c.Request().Context(),
		user,
	)
	if err != nil {
		return err
	}

	chartItems := domain.ConvertDateCountToChartDataItem(stats.YearlyCount)
	view := transformer.TransformSingleDatasetToViewmodel(chartItems)

	return Render(c, statistics.CountByYear(view))
}

func (h *StatisticsHandler) GetCountCumulative(c echo.Context) error {
	user := session.GetUserFromContext(c)
	stats, err := h.statisticService.ComputeStatistics(
		c.Request().Context(),
		user,
	)
	if err != nil {
		return err
	}

	chartItems := domain.ConvertDateCountToChartDataItem(stats.CumulativeMonthlyCount)
	view := transformer.TransformSingleDatasetToViewmodel(chartItems)

	return Render(c, statistics.CountCumul(view))
}

func (h *StatisticsHandler) GetDurationByMonth(c echo.Context) error {
	user := session.GetUserFromContext(c)
	stats, err := h.statisticService.ComputeStatistics(
		c.Request().Context(),
		user,
	)
	if err != nil {
		return err
	}

	chartItems := domain.ConvertDateDurationToChartDataItem(stats.MonthlyDuration)
	view := transformer.TransformMultiDatasetsToViewmodel(chartItems)

	return Render(c, statistics.DurationByMonth(view))
}

func (h *StatisticsHandler) GetDurationByYear(c echo.Context) error {
	user := session.GetUserFromContext(c)
	stats, err := h.statisticService.ComputeStatistics(
		c.Request().Context(),
		user,
	)
	if err != nil {
		return err
	}

	chartItems := domain.ConvertDateDurationToChartDataItem(stats.YearlyDuration)
	view := transformer.TransformSingleDatasetToViewmodel(chartItems)

	return Render(c, statistics.DurationByYear(view))
}

func (h *StatisticsHandler) GetDurationCumulative(c echo.Context) error {
	user := session.GetUserFromContext(c)
	stats, err := h.statisticService.ComputeStatistics(
		c.Request().Context(),
		user,
	)
	if err != nil {
		return err
	}

	chartItems := domain.ConvertDateDurationToChartDataItem(stats.CumulativeMonthlyDuration)
	view := transformer.TransformSingleDatasetToViewmodel(chartItems)

	return Render(c, statistics.DurationCumul(view))
}
