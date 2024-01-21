package handler

import (
	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/internal/service"
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/transformer"
	"github.com/AurelienS/cigare/web/view/progression"
	"github.com/AurelienS/cigare/web/viewmodel"
	"github.com/labstack/echo/v4"
)

type ProgressionHandler struct {
	statisticService service.StatisticService
}

func NewProgressionHandler(statisticService service.StatisticService) ProgressionHandler {
	return ProgressionHandler{statisticService: statisticService}
}

func (h *ProgressionHandler) GetProgression(c echo.Context) error {
	user := session.GetUserFromContext(c)
	statsYearMonth, err := h.statisticService.GetStatisticsByYearAndMonth(
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
	return Render(c, progression.Index(view))
}
