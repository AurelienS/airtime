package handler

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/AurelienS/cigare/internal/logbook"
	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/template"
	logbookView "github.com/AurelienS/cigare/web/template/logbook"
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

func (h *LogbookHandler) RedirectToCurrentYearLogbook(c echo.Context) error {
	return c.Redirect(301, "/logbook/"+time.Now().Format("2006"))
}

func getDescription(value string, year int) string {
	return fmt.Sprintf("%s en %d", value, year)
}

func sortAndConvertToViewModel(flights []model.Flight) []template.FlightView {
	sort.Slice(flights, func(i, j int) bool {
		return flights[i].Date.After(flights[j].Date)
	})

	var flightViews []template.FlightView
	for _, f := range flights {
		flightViews = append(flightViews, template.FlightView{
			TakeoffLocation:  f.TakeoffLocation,
			Date:             f.Date.Local().Format("02/01 15:04"),
			TotalThermicTime: prettyDuration(f.Statistic.TotalThermicTime),
			TotalFlightTime:  prettyDuration(f.Statistic.TotalFlightTime),
			MaxClimbRate:     strconv.FormatFloat(f.Statistic.MaxClimbRate, 'f', 1, 64),
			MaxAltitude:      strconv.Itoa(f.Statistic.MaxAltitude),
		})
	}
	return flightViews
}

func getStatCat1(yearStats, allTimeStats logbook.Stats, year int) []template.StatView {
	return []template.StatView{
		{
			Title:       "Nombre de vols",
			Value:       strconv.Itoa(allTimeStats.FlightCount),
			Description: getDescription(strconv.Itoa(yearStats.FlightCount), year),
		},
		{
			Title:       "Temps de vol total",
			Value:       prettyDuration(allTimeStats.TotalFlightTime),
			Description: getDescription(prettyDuration(yearStats.TotalFlightTime), year),
		},
		{
			Title:       "Montée totale",
			Value:       prettyAltitude(allTimeStats.TotalClimb),
			Description: getDescription(prettyAltitude(yearStats.TotalClimb), year),
		},
		{
			Title:       "Temps total en thermique",
			Value:       prettyDuration(allTimeStats.TotalThermicTime),
			Description: getDescription(prettyDuration(yearStats.TotalThermicTime), year),
		},
		{
			Title:       "Nombre total de thermiques",
			Value:       strconv.Itoa(allTimeStats.TotalNumberOfThermals),
			Description: getDescription(strconv.Itoa(yearStats.TotalNumberOfThermals), year),
		},
	}
}

func getStatCat2(yearStats, allTimeStats logbook.Stats, year int) []template.StatView {
	return []template.StatView{
		{
			Title:       "Durée moyenne de vol",
			Value:       prettyDuration(allTimeStats.AverageFlightLength),
			Description: getDescription(prettyDuration(yearStats.AverageFlightLength), year),
		},
		{
			Title:       "Durée maximale de vol",
			Value:       prettyDuration(allTimeStats.MaxFlightLength),
			Description: getDescription(prettyDuration(yearStats.MaxFlightLength), year),
		},
		{
			Title:       "Durée minimale de vol",
			Value:       prettyDuration(allTimeStats.MinFlightLength),
			Description: getDescription(prettyDuration(yearStats.MinFlightLength), year),
		},
		{
			Title:       "Altitude maximale",
			Value:       strconv.Itoa(allTimeStats.MaxAltitude) + "m",
			Description: getDescription(strconv.Itoa(yearStats.MaxAltitude)+"m", year),
		},
		{
			Title:       "Plus grand thermique",
			Value:       strconv.Itoa(allTimeStats.MaxClimb) + "m",
			Description: getDescription(strconv.Itoa(yearStats.MaxClimb)+"m", year),
		},

		{
			Title:       "Taux de montée maximal",
			Value:       fmt.Sprintf("%.2f", allTimeStats.MaxClimbRate) + "m/s",
			Description: getDescription(fmt.Sprintf("%.2f", yearStats.MaxClimbRate)+"m/s", year),
		},
	}
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

	flightViews := sortAndConvertToViewModel(flights)

	allTimeStats, err := h.LogbookService.GetStatistics(c.Request().Context(), 0, user)
	if err != nil {
		return err
	}

	yearStats, err := h.LogbookService.GetStatistics(c.Request().Context(), year, user)
	if err != nil {
		return err
	}

	statCat1 := getStatCat1(yearStats, allTimeStats, year)
	statCat2 := getStatCat2(yearStats, allTimeStats, year)

	viewData := template.LogbookView{
		Year:    yearParam,
		Flights: flightViews,
		Stats1:  statCat1,
		Stats2:  statCat2,
	}

	return Render(c, logbookView.Logbook(viewData))
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
