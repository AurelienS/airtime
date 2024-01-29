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
	"rgb(8, 76, 223)",
	"rgb(223, 8, 8)",
	"rgb(8, 223, 8)",
	"rgb(223, 223, 8)",
	"rgb(8, 8, 223)",
	"rgb(223, 8, 223)",
	"rgb(8, 223, 223)",
	"rgb(223, 128, 8)",
	"rgb(128, 8, 223)",
	"rgb(8, 223, 128)",
	"rgb(128, 223, 8)",
	"rgb(223, 8, 128)",
	"rgb(8, 128, 223)",
	"rgb(128, 128, 8)",
	"rgb(8, 128, 128)",
	"rgb(128, 8, 128)",
	"rgb(223, 128, 223)",
	"rgb(128, 223, 128)",
	"rgb(223, 128, 128)",
	"rgb(128, 128, 223)",
}

type (
	StatIntExtractor   func(stats domain.MultipleFlightStats) int
	StatFloatExtractor func(stats domain.MultipleFlightStats) float64
)

func TransformChartViewModel(
	statsYearMonth service.StatsYearMonth,
	extractor StatIntExtractor,
) []viewmodel.CountDataset {
	datasets := []viewmodel.CountDataset{}

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

		monthDataset := viewmodel.CountDataset{
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

func TransformChartTimeViewModel(
	statsYearMonth service.StatsYearMonth,
	extractor StatFloatExtractor,
) []viewmodel.TimeDataset {
	datasets := []viewmodel.TimeDataset{}

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

		monthDataset := viewmodel.TimeDataset{
			Label: strconv.Itoa(year),
			Color: datasetColors[colorCounter%len(datasetColors)],
			Data:  []float64{},
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
