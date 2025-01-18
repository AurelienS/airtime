package transformer

import (
	"fmt"
	"sort"
	"time"

	"github.com/AurelienS/airtime/internal/domain"
	"github.com/AurelienS/airtime/web/viewmodel"
)

func TransformMultiDatasetsToViewmodel(dataItems []domain.ChartDataItem) viewmodel.ChartData {
	yearlyData := make(map[int][]string)

	for _, item := range dataItems {
		year, month, _ := item.GetDate().Date()
		if _, exists := yearlyData[year]; !exists {
			yearlyData[year] = make([]string, 12)
		}
		yearlyData[year][month-1] = item.GetValue()
	}

	var years []int
	for year := range yearlyData {
		years = append(years, year)
	}
	sort.Slice(years, func(i, j int) bool { return years[i] > years[j] })

	labels := make([]string, 0, 12)
	for i := 1; i <= 12; i++ {
		month := time.Month(i)
		labels = append(labels, frenchMonthName(month))
	}

	datasets := make([]viewmodel.ChartDataset, 0, len(yearlyData))
	colorIndex := 0
	for _, year := range years {
		datasets = append(datasets, viewmodel.ChartDataset{
			Label: fmt.Sprintf("%d", year),
			Data:  yearlyData[year],
			Color: datasetColors[colorIndex%len(datasetColors)],
		})
		colorIndex++
	}

	return viewmodel.ChartData{
		Labels:   labels,
		Datasets: datasets,
	}
}

func TransformSingleDatasetToViewmodel(dataItems []domain.ChartDataItem) viewmodel.ChartData {
	var labels []string
	var data []string

	for _, item := range dataItems {
		labels = append(labels, item.GetDate().Format("Jan 2006"))
		data = append(data, item.GetValue())
	}

	dataset := viewmodel.ChartDataset{
		Data:  data,
		Color: datasetColors[0],
	}

	return viewmodel.ChartData{
		Labels:   labels,
		Datasets: []viewmodel.ChartDataset{dataset},
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

func frenchMonthName(month time.Month) string {
	frenchMonths := map[time.Month]string{
		time.January:   "Janvier",
		time.February:  "Février",
		time.March:     "Mars",
		time.April:     "Avril",
		time.May:       "Mai",
		time.June:      "Juin",
		time.July:      "Juillet",
		time.August:    "Août",
		time.September: "Septembre",
		time.October:   "Octobre",
		time.November:  "Novembre",
		time.December:  "Décembre",
	}

	return frenchMonths[month]
}
