package handler

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/AurelienS/cigare/internal/logbook"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/template/component"
	"github.com/AurelienS/cigare/web/template/flight"
	"github.com/AurelienS/cigare/web/template/page"
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

func prettyDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh%d", hours, minutes)
	}
	return fmt.Sprintf("%dmin", minutes)
}

func prettyAltitude(alt int) string {
	km := alt / 1000
	m := alt % 1000

	if km > 0 {
		return fmt.Sprintf("%dkm", km)
	}
	return fmt.Sprintf("%dm", m)
}

func (h *LogbookHandler) GetPage(c echo.Context) error {
	user := session.GetUserFromContext(c)

	flights, err := h.LogbookService.GetFlights(c.Request().Context(), user)
	if err != nil {
		return err
	}

	sort.Slice(flights, func(i, j int) bool {
		return flights[i].Date.After(flights[j].Date)
	})

	var flightViews []flight.FlightView
	for _, f := range flights {
		flightViews = append(flightViews, flight.FlightView{
			TakeoffLocation:  f.TakeoffLocation,
			Date:             f.Date.Local().Format("02/01 15:04"),
			TotalThermicTime: prettyDuration(f.Statistic.TotalThermicTime),
			TotalFlightTime:  prettyDuration(f.Statistic.TotalFlightTime),
			MaxClimbRate:     strconv.FormatFloat(f.Statistic.MaxClimbRate, 'f', 1, 64),
			MaxAltitude:      strconv.Itoa(f.Statistic.MaxAltitude),
		})
	}

	stats, err := h.LogbookService.GetStatistics(c.Request().Context(), user)
	if err != nil {
		return err
	}

	viewData := page.LogbookView{
		Flights: flightViews,
		Stats1: []component.StatView{
			{
				Title:       "Nombre de vols",
				Value:       strconv.Itoa(stats.FlightCount),
				Description: "XX cette année",
			},
			{
				Title:       "Temps de vol total",
				Value:       prettyDuration(stats.TotalFlightTime),
				Description: "XX cette année",
			},
			{
				Title:       "Montée totale",
				Value:       prettyAltitude(stats.TotalClimb),
				Description: "XX cette année",
			},
			{
				Title:       "Temps total en thermique",
				Value:       prettyDuration(stats.TotalThermicTime),
				Description: "XX cette année",
			},
			{
				Title:       "Nombre total de thermiques",
				Value:       strconv.Itoa(stats.TotalNumberOfThermals),
				Description: "XX cette année",
			},
		},
		Stats2: []component.StatView{
			{
				Title:       "Durée moyenne de vol",
				Value:       prettyDuration(stats.AverageFlightLength),
				Description: "XX cette année",
			},
			{
				Title:       "Durée maximale de vol",
				Value:       prettyDuration(stats.MaxFlightLength),
				Description: "XX cette année",
			},
			{
				Title:       "Durée minimale de vol",
				Value:       prettyDuration(stats.MinFlightLength),
				Description: "XX cette année",
			},
			{
				Title:       "Altitude maximale",
				Value:       strconv.Itoa(stats.MaxAltitude) + "m",
				Description: "XX cette année",
			},
			{
				Title:       "Plus grand thermique",
				Value:       strconv.Itoa(stats.MaxClimb) + "m",
				Description: "XX cette année",
			},

			{
				Title:       "Taux de montée maximal",
				Value:       fmt.Sprintf("%.2f", stats.MaxClimbRate) + "m/s",
				Description: "XX cette année",
			},
		},
	}

	return Render(c, page.Logbook(viewData))
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

	return h.GetPage(c)
}
