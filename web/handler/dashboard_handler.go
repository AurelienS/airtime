package handler

import (
	"time"

	"github.com/AurelienS/airtime/internal/service"
	"github.com/AurelienS/airtime/web/session"
	"github.com/AurelienS/airtime/web/transformer"
	"github.com/AurelienS/airtime/web/view/dashboard"
	"github.com/AurelienS/airtime/web/viewmodel"
	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	flightService    service.FlightService
	statisticService service.StatisticService
}

func NewDashboardHandler(
	flightService service.FlightService,
	statisticService service.StatisticService,
) DashboardHandler {
	return DashboardHandler{
		flightService:    flightService,
		statisticService: statisticService,
	}
}

func (h DashboardHandler) GetIndex(c echo.Context) error {
	user := session.GetUserFromContext(c)
	ctx := c.Request().Context()

	allTimeStats, err := h.statisticService.GetGlobalStats(ctx, user, time.Time{}, time.Now())
	if err != nil {
		return err
	}

	startOfCurrentYear := time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	currentYearStats, err := h.statisticService.GetGlobalStats(ctx, user, startOfCurrentYear, time.Now())
	if err != nil {
		return err
	}

	sitesStats := viewmodel.DashboardSitesStatsView{}

	lastFlights, err := h.flightService.GetLastFlights(ctx, 8, user)
	if err != nil {
		return err
	}

	view := transformer.TransformDashboardToViewModel(allTimeStats, currentYearStats, lastFlights, sitesStats, user)

	return Render(c, dashboard.Index(view))
}
