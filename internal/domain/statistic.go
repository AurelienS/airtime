package domain

import "time"

type DateCount struct {
	Date  time.Time
	Count int
}

func (dc DateCount) GetDate() time.Time {
	return dc.Date
}

func (dc DateCount) GetValue() float64 {
	return float64(dc.Count)
}

type DateDuration struct {
	Date     time.Time
	Duration time.Duration
}

func (dd DateDuration) GetDate() time.Time {
	return dd.Date
}

func (dd DateDuration) GetValue() float64 {
	return dd.Duration.Hours()
}

type Statistics struct {
	MonthlyCount              []DateCount
	YearlyCount               []DateCount
	CumulativeMonthlyCount    []DateCount
	MonthlyDuration           []DateDuration
	YearlyDuration            []DateDuration
	CumulativeMonthlyDuration []DateDuration
}

type ChartDataItem interface {
	GetDate() time.Time
	GetValue() float64
}

func ConvertDateCountToChartDataItem(dateCounts []DateCount) []ChartDataItem {
	items := make([]ChartDataItem, len(dateCounts))
	for i, dc := range dateCounts {
		items[i] = dc
	}
	return items
}

func ConvertDateDurationToChartDataItem(dateDurations []DateDuration) []ChartDataItem {
	items := make([]ChartDataItem, len(dateDurations))
	for i, dc := range dateDurations {
		items[i] = dc
	}
	return items
}
