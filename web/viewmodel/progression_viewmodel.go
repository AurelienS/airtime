package viewmodel

type ChartDataset struct {
	Label string
	Data  []string
	Color string
}

type ChartData struct {
	Labels   []string
	Datasets []ChartDataset
}
