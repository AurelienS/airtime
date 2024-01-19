package transformer

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/web/viewmodel"
)

func TransformLogbookToViewModel(
	flights []model.Flight,
	yearStats model.StatsAggregated,
	allTimeStats model.StatsAggregated,
	year int,
	flyingYears []int,
	isFlightAdded bool,
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
		IsFlightAdded:  isFlightAdded,
	}
}

var datasetColors = []string{
	"rgb(255, 99, 132)",  // Rouge clair
	"rgb(54, 162, 235)",  // Bleu clair
	"rgb(255, 206, 86)",  // Jaune
	"rgb(75, 192, 192)",  // Turquoise
	"rgb(153, 102, 255)", // Violet
	"rgb(255, 159, 64)",  // Orange
	"rgb(199, 199, 199)", // Gris
	"rgb(83, 109, 254)",  // Bleu foncé
	"rgb(255, 99, 71)",   // Corail
	"rgb(144, 238, 144)", // Vert clair
	"rgb(255, 215, 0)",   // Or
	"rgb(218, 165, 32)",  // Bronze
	"rgb(106, 90, 205)",  // Ardoise
	"rgb(255, 127, 80)",  // Saumon
	"rgb(0, 128, 128)",   // Sarcelle
	"rgb(0, 255, 127)",   // Vert printemps
	"rgb(255, 182, 193)", // Rose clair
	"rgb(107, 142, 35)",  // Olive
	"rgb(75, 0, 130)",    // Indigo
	"rgb(255, 69, 0)",    // Rouge orangé
}

type StatExtractor func(stats model.StatsAggregated) int

func TransformStatsViewModel(statsYearMonth model.StatsYearMonth, extractor StatExtractor) []viewmodel.ChartDataset {
	datasets := []viewmodel.ChartDataset{}

	// Create a slice of years to sort
	years := make([]int, 0, len(statsYearMonth))
	for year := range statsYearMonth {
		years = append(years, year)
	}

	// Sort years slice in reverse order
	sort.Sort(sort.Reverse(sort.IntSlice(years)))

	colorCounter := 0

	for _, year := range years {
		monthStatsMap := statsYearMonth[year]

		// Create a slice of months to sort
		months := make([]time.Month, 0, len(monthStatsMap))
		for month := range monthStatsMap {
			months = append(months, month)
		}

		// Sort months slice
		sort.Slice(months, func(i, j int) bool {
			return months[i] < months[j]
		})

		monthDataset := viewmodel.ChartDataset{
			Label: strconv.Itoa(year),
			Color: datasetColors[colorCounter%len(datasetColors)],
			Data:  []int{},
		}
		colorCounter++

		// Append stats for each month in sorted order
		for _, month := range months {
			stats := monthStatsMap[month]
			// Use the extractor function to get the specific stat
			monthDataset.Data = append(monthDataset.Data, extractor(stats))
		}

		datasets = append(datasets, monthDataset)
	}

	return datasets
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
			TotalThermicTime: prettyDuration(f.Statistic.TotalThermicTime, false),
			TotalFlightTime:  prettyDuration(f.Statistic.TotalFlightTime, true),
			MaxClimbRate:     strconv.FormatFloat(f.Statistic.MaxClimbRate, 'f', 1, 64),
			MaxAltitude:      strconv.Itoa(f.Statistic.MaxAltitude),
		})
	}
	return flightViews
}

func getMainStat(yearStats, allTimeStats model.StatsAggregated) []viewmodel.StatView {
	return []viewmodel.StatView{
		{
			Title:            "Nombre de vols",
			AlltimeValue:     strconv.Itoa(allTimeStats.FlightCount),
			CurrentYearValue: strconv.Itoa(yearStats.FlightCount),
		},
		{
			Title:            "Temps de vol total",
			AlltimeValue:     prettyDuration(allTimeStats.TotalFlightTime, true),
			CurrentYearValue: prettyDuration(yearStats.TotalFlightTime, false),
		},
		{
			Title:            "Montée totale",
			AlltimeValue:     prettyAltitude(allTimeStats.TotalClimb, false),
			CurrentYearValue: prettyAltitude(yearStats.TotalClimb, false),
		},
		{
			Title:            "Temps total en thermique",
			AlltimeValue:     prettyDuration(allTimeStats.TotalThermicTime, true),
			CurrentYearValue: prettyDuration(yearStats.TotalThermicTime, false),
		},
		{
			Title:            "Nombre total de thermiques",
			AlltimeValue:     strconv.Itoa(allTimeStats.TotalNumberOfThermals),
			CurrentYearValue: strconv.Itoa(yearStats.TotalNumberOfThermals),
		},
	}
}

func getSecondaryStat(yearStats, allTimeStats model.StatsAggregated) []viewmodel.StatView {
	return []viewmodel.StatView{
		{
			Title:            "Durée moyenne de vol",
			AlltimeValue:     prettyDuration(allTimeStats.AverageFlightLength, false),
			CurrentYearValue: prettyDuration(yearStats.AverageFlightLength, false),
		},
		{
			Title:            "Durée maximale de vol",
			AlltimeValue:     prettyDuration(allTimeStats.MaxFlightLength, false),
			CurrentYearValue: prettyDuration(yearStats.MaxFlightLength, false),
		},
		{
			Title:            "Durée minimale de vol",
			AlltimeValue:     prettyDuration(allTimeStats.MinFlightLength, false),
			CurrentYearValue: prettyDuration(yearStats.MinFlightLength, false),
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

func prettyDuration(d time.Duration, onlyHour bool) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if onlyHour && hours > 0 {
		return fmt.Sprintf("%d h", hours)
	}

	if hours > 0 {
		return fmt.Sprintf("%dh%02d", hours, minutes)
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
