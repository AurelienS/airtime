package handler

import (
	"strconv"
	"time"

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

func (h *LogbookHandler) RedirectToCurrentYearLogbook(c echo.Context) error {
	return c.Redirect(301, "/logbook/"+time.Now().Format("2006"))
}

func (h *LogbookHandler) Get(c echo.Context) error {
	user := session.GetUserFromContext(c)
	yearParam := c.Param("year")

	if yearParam == "year" {
		yearParam = c.FormValue("yearValue")
	}

	year, err := strconv.Atoi(yearParam)
	if err != nil {
		return err
	}

	flights, err := h.LogbookService.GetFlights(c.Request().Context(), year, user)
	if err != nil {
		return err
	}

	allTimeStats, err := h.LogbookService.GetStatistics(c.Request().Context(), 0, user)
	if err != nil {
		return err
	}

	yearStats, err := h.LogbookService.GetStatistics(c.Request().Context(), year, user)
	if err != nil {
		return err
	}

	viewData := transformer.TransformLogbookToViewModel(flights, yearStats, allTimeStats, year)

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
