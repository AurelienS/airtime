package transformer

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/AurelienS/cigare/internal/logbook"
	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/web/viewmodel"
)

func TransformLogbookToViewModel(
	flights []model.Flight,
	yearStats logbook.Stats,
	allTimeStats logbook.Stats,
	year int,
) viewmodel.LogbookView {
	flightViews := sortAndConvertToViewModel(flights)
	statMain := getMainStat(yearStats, allTimeStats, year)
	statSecondary := getSecondaryStat(yearStats, allTimeStats, year)

	return viewmodel.LogbookView{
		Year:          strconv.Itoa(year),
		Flights:       flightViews,
		StatMain:      statMain,
		StatSecondary: statSecondary,
	}
}

func sortAndConvertToViewModel(flights []model.Flight) []viewmodel.FlightView {
	sort.Slice(flights, func(i, j int) bool {
		return flights[i].Date.After(flights[j].Date)
	})

	var flightViews []viewmodel.FlightView
	for _, f := range flights {
		flightViews = append(flightViews, viewmodel.FlightView{
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

func getDescription(value string, year int) string {
	return fmt.Sprintf("%s en %d", value, year)
}

func getMainStat(yearStats, allTimeStats logbook.Stats, year int) []viewmodel.StatView {
	return []viewmodel.StatView{
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

func getSecondaryStat(yearStats, allTimeStats logbook.Stats, year int) []viewmodel.StatView {
	return []viewmodel.StatView{
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
