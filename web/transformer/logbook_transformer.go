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
	flyingYears []int,
) viewmodel.LogbookView {
	flightViews := sortAndConvertToViewModel(flights)
	statMain := getMainStat(yearStats, allTimeStats)
	statSecondary := getSecondaryStat(yearStats, allTimeStats)

	flyingYearsString := make([]string, 0, len(flyingYears))
	for _, y := range flyingYears {
		flyingYearsString = append(flyingYearsString, strconv.Itoa(y))
	}

	return viewmodel.LogbookView{
		CurrentYear:    strconv.Itoa(year),
		AvailableYears: flyingYearsString,
		Flights:        flightViews,
		StatMain:       statMain,
		StatSecondary:  statSecondary,
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

func getMainStat(yearStats, allTimeStats logbook.Stats) []viewmodel.StatView {
	return []viewmodel.StatView{
		{
			Title:            "Nombre de vols",
			AlltimeValue:     strconv.Itoa(allTimeStats.FlightCount),
			CurrentYearValue: strconv.Itoa(yearStats.FlightCount),
		},
		{
			Title:            "Temps de vol total",
			AlltimeValue:     prettyDuration(allTimeStats.TotalFlightTime),
			CurrentYearValue: prettyDuration(yearStats.TotalFlightTime),
		},
		{
			Title:            "Montée totale",
			AlltimeValue:     prettyAltitude(allTimeStats.TotalClimb, false),
			CurrentYearValue: prettyAltitude(yearStats.TotalClimb, false),
		},
		{
			Title:            "Temps total en thermique",
			AlltimeValue:     prettyDuration(allTimeStats.TotalThermicTime),
			CurrentYearValue: prettyDuration(yearStats.TotalThermicTime),
		},
		{
			Title:            "Nombre total de thermiques",
			AlltimeValue:     strconv.Itoa(allTimeStats.TotalNumberOfThermals),
			CurrentYearValue: strconv.Itoa(yearStats.TotalNumberOfThermals),
		},
	}
}

func getSecondaryStat(yearStats, allTimeStats logbook.Stats) []viewmodel.StatView {
	return []viewmodel.StatView{
		{
			Title:            "Durée moyenne de vol",
			AlltimeValue:     prettyDuration(allTimeStats.AverageFlightLength),
			CurrentYearValue: prettyDuration(yearStats.AverageFlightLength),
		},
		{
			Title:            "Durée maximale de vol",
			AlltimeValue:     prettyDuration(allTimeStats.MaxFlightLength),
			CurrentYearValue: prettyDuration(yearStats.MaxFlightLength),
		},
		{
			Title:            "Durée minimale de vol",
			AlltimeValue:     prettyDuration(allTimeStats.MinFlightLength),
			CurrentYearValue: prettyDuration(yearStats.MinFlightLength),
		},
		{
			Title:            "Altitude maximale",
			AlltimeValue:     prettyAltitude(allTimeStats.MaxAltitude, true),
			CurrentYearValue: prettyAltitude(yearStats.MaxAltitude, true),
		},
		{
			Title:            "Plus grand thermique",
			AlltimeValue:     prettyAltitude(allTimeStats.MaxClimb, true),
			CurrentYearValue: prettyAltitude(yearStats.MaxClimb, true),
		},

		{
			Title:            "Taux de montée maximal",
			AlltimeValue:     fmt.Sprintf("%.2f", allTimeStats.MaxClimbRate) + " m/s",
			CurrentYearValue: fmt.Sprintf("%.2f", yearStats.MaxClimbRate) + " m/s",
		},
	}
}

func prettyDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh%d", hours, minutes)
	}
	return fmt.Sprintf("%d min", minutes)
}

func prettyAltitude(alt int, forceMeter bool) string {
	if forceMeter {
		return strconv.Itoa(alt) + " m"
	}
	km := alt / 1000
	m := alt % 1000

	if km > 0 {
		return fmt.Sprintf("%d km", km)
	}
	return fmt.Sprintf("%d m", m)
}
