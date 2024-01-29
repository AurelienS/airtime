package viewmodel

type LogbookView struct {
	Flights        []FlightView
	AvailableYears []string
	CurrentYear    string
}

type FlightView struct {
	ID          string
	Fulldate    string
	Date        string
	Location    string
	Duration    string
	AltitudeMax string
	Distance    string
	Link        string
}

type FlightDetailStatView struct {
	Title string
	Value string
}
type FlightDetailView struct {
	FlightView
	UserView
	FlightGeoJSON string
}

type StatView struct {
	Title            string
	AlltimeValue     string
	CurrentYearValue string
}

type StatisticsView struct {
	User UserView

	FlightTimeMonthlyData  []CountDataset
	FlightCountMonthlyData []CountDataset
}

type CountDataset struct {
	Label string
	Data  []int
	Color string
}

type CountData struct {
	Labels   []string
	Datasets []CountDataset
}

type TimeDataset struct {
	Label string
	Data  []float64
	Color string
}

type TimeData struct {
	Labels   []string
	Datasets []TimeDataset
}
