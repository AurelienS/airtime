package viewmodel

type ChartDataset struct {
	Label string
	Data  []float64
	Color string
}

type ChartData struct {
	Labels   []string
	Datasets []ChartDataset
}
