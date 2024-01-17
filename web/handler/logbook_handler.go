package handler

import (
	"fmt"
	"strconv"

	"github.com/AurelienS/cigare/internal/logbook"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/transformer"
	"github.com/AurelienS/cigare/web/view/logbookview"
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

func (h *LogbookHandler) RedirectToLastYearLogbook(c echo.Context) error {
	lastFlight, err := h.LogbookService.GetLastFlight(c.Request().Context(), session.GetUserFromContext(c))
	if err != nil {
		return err
	}
	lastYear := lastFlight.Date.Year()
	redirectTo := fmt.Sprintf("/logbook/%d", lastYear)
	return c.Redirect(301, redirectTo)
}

func (h *LogbookHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	user := session.GetUserFromContext(c)
	yearParam := c.Param("year")

	if yearParam == "year" {
		yearParam = c.FormValue("yearValue")
	}

	year, err := strconv.Atoi(yearParam)
	if err != nil {
		return err
	}

	flyingYears, err := h.LogbookService.GetFlyingYears(ctx, user)
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
		redirectTo := fmt.Sprintf("/logbook/%d", lastYear)
		return c.Redirect(301, redirectTo)
	}

	flights, err := h.LogbookService.GetFlights(ctx, year, user)
	if err != nil {
		return err
	}

	allTimeStats, err := h.LogbookService.GetStatistics(ctx, 0, user)
	if err != nil {
		return err
	}

	yearStats, err := h.LogbookService.GetStatistics(ctx, year, user)
	if err != nil {
		return err
	}

	viewData := transformer.TransformLogbookToViewModel(flights, yearStats, allTimeStats, year, flyingYears)

	return Render(c, logbookview.Logbook(viewData))
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

	return h.Get(c)
}
