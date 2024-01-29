package transformer

import (
	"github.com/AurelienS/cigare/internal/service"
	"github.com/AurelienS/cigare/web/viewmodel"
)

func TransformChartCountCumulative(flightCounts []service.FlightCount) viewmodel.CountData {
	labels := make([]string, 0, len(flightCounts))
	values := make([]int, 0, len(flightCounts))

	for _, flightCount := range flightCounts {
		labels = append(labels, flightCount.Date.Format("01/06"))
		values = append(values, flightCount.Count)
	}
	a := viewmodel.CountDataset{
		Label: "count",
		Data:  values,
		Color: "rgb(8,76,223)",
	}

	return viewmodel.CountData{
		Labels:   labels,
		Datasets: []viewmodel.CountDataset{a},
	}
}

func TransformChartTimeCumulative(flightCounts []service.FlightDuration) viewmodel.TimeData {
	labels := make([]string, 0, len(flightCounts))
	values := make([]float64, 0, len(flightCounts))

	for _, flightCount := range flightCounts {
		labels = append(labels, flightCount.Date.Format("01/06"))
		values = append(values, flightCount.Duration.Hours())
	}
	a := viewmodel.TimeDataset{
		Label: "count",
		Data:  values,
		Color: "rgb(8,76,223)",
	}

	return viewmodel.TimeData{
		Labels:   labels,
		Datasets: []viewmodel.TimeDataset{a},
	}
}
