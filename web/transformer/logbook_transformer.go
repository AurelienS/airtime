package transformer

import (
	"sort"
	"strconv"
	"time"

	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/internal/service"
	"github.com/AurelienS/cigare/web/viewmodel"
)

func TransformLogbookToViewModel(
	flights []domain.Flight,
	flyingYears []int,
	year int,
) viewmodel.LogbookView {
	var flightViews []viewmodel.FlightView
	for _, f := range flights {
		flightViews = append(flightViews, TransformFlightToViewmodel(f))
	}

	flyingYearsString := make([]string, 0, len(flyingYears))
	for _, y := range flyingYears {
		flyingYearsString = append(flyingYearsString, strconv.Itoa(y))
	}

	return viewmodel.LogbookView{
		CurrentYear:    strconv.Itoa(year),
		AvailableYears: flyingYearsString,
		Flights:        flightViews,
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

type StatExtractor func(stats domain.MultipleFlightStats) int

func TransformChartViewModel(statsYearMonth service.StatsYearMonth, extractor StatExtractor) []viewmodel.ChartDataset {
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
