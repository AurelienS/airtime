package transformer

import (
	"sort"
	"strconv"
	"time"

	"github.com/AurelienS/cigare/web/viewmodel"
)

func TransformCumulativeCount(cumulativeCount map[int]map[time.Month]int) viewmodel.CountData {
	var labels []string
	var values []int

	for _, year := range sortYearsInt(cumulativeCount) {
		yearMap := cumulativeCount[year]
		for _, month := range sortMonthsInt(yearMap) {
			count := yearMap[month]
			labels = append(labels, generateDateLabel(year, month))
			values = append(values, count)
		}
	}

	dataset := viewmodel.CountDataset{
		Label: "count",
		Data:  values,
		Color: "rgb(8,76,223)",
	}

	return viewmodel.CountData{
		Labels:   labels,
		Datasets: []viewmodel.CountDataset{dataset},
	}
}

func TransformChartTimeCumulative(cumulativeMonthlyDuration map[int]map[time.Month]time.Duration) viewmodel.TimeData {
	var labels []string
	var values []float64

	for _, year := range sortYearsDuration(cumulativeMonthlyDuration) {
		yearMap := cumulativeMonthlyDuration[year]
		for _, month := range sortMonthsDuration(yearMap) {
			duration := yearMap[month]
			labels = append(labels, generateDateLabel(year, month))
			values = append(values, duration.Hours())
		}
	}

	dataset := viewmodel.TimeDataset{
		Label: "Time",
		Data:  values,
		Color: "rgb(8,76,223)",
	}

	return viewmodel.TimeData{
		Labels:   labels,
		Datasets: []viewmodel.TimeDataset{dataset},
	}
}

func TransformMonthlyCountToViewmodel(
	monthlyCountByYear map[int]map[time.Month]int,
) []viewmodel.CountDataset {
	var datasets []viewmodel.CountDataset
	colorCounter := 0

	for _, year := range getSortedYears(monthlyCountByYear) {
		dataset := initCountDataset(year, colorCounter)
		colorCounter++

		for month := time.January; month <= time.December; month++ {
			count := monthlyCountByYear[year][month]
			dataset.Data[int(month)-1] = count
		}

		datasets = append(datasets, dataset)
	}

	return datasets
}

func TransformMonthlyTimeToViewmodel(
	monthlyDurationByYear map[int]map[time.Month]time.Duration,
) []viewmodel.TimeDataset {
	var datasets []viewmodel.TimeDataset
	colorCounter := 0

	for _, year := range getSortedYears2(monthlyDurationByYear) {
		dataset := initTimeDataset(year, colorCounter)
		colorCounter++

		for month := time.January; month <= time.December; month++ {
			duration := monthlyDurationByYear[year][month]
			dataset.Data[int(month)-1] = duration.Hours()
		}

		datasets = append(datasets, dataset)
	}

	return datasets
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

func sortYearsInt(yearMap map[int]map[time.Month]int) []int {
	years := make([]int, 0, len(yearMap))
	for year := range yearMap {
		years = append(years, year)
	}
	sort.Ints(years)
	return years
}

func sortYearsDuration(yearMap map[int]map[time.Month]time.Duration) []int {
	years := make([]int, 0, len(yearMap))
	for year := range yearMap {
		years = append(years, year)
	}
	sort.Ints(years)
	return years
}

func sortMonthsInt(monthMap map[time.Month]int) []time.Month {
	months := make([]time.Month, 0, len(monthMap))
	for month := range monthMap {
		months = append(months, month)
	}
	sort.Slice(months, func(i, j int) bool { return months[i] < months[j] })
	return months
}

func sortMonthsDuration(monthMap map[time.Month]time.Duration) []time.Month {
	months := make([]time.Month, 0, len(monthMap))
	for month := range monthMap {
		months = append(months, month)
	}
	sort.Slice(months, func(i, j int) bool { return months[i] < months[j] })
	return months
}

func generateDateLabel(year int, month time.Month) string {
	date := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	return date.Format("01/06")
}

func getSortedYears(yearMap map[int]map[time.Month]int) []int {
	years := make([]int, 0, len(yearMap))
	for year := range yearMap {
		years = append(years, year)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(years)))
	return years
}

func getSortedYears2(yearMap map[int]map[time.Month]time.Duration) []int {
	years := make([]int, 0, len(yearMap))
	for year := range yearMap {
		years = append(years, year)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(years)))
	return years
}

func initTimeDataset(year, colorCounter int) viewmodel.TimeDataset {
	return viewmodel.TimeDataset{
		Label: strconv.Itoa(year),
		Color: datasetColors[colorCounter%len(datasetColors)],
		Data:  make([]float64, 12),
	}
}

func initCountDataset(year, colorCounter int) viewmodel.CountDataset {
	return viewmodel.CountDataset{
		Label: strconv.Itoa(year),
		Color: datasetColors[colorCounter%len(datasetColors)],
		Data:  make([]int, 12),
	}
}
