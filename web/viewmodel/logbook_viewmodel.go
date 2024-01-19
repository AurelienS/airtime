package viewmodel

type LogbookView struct {
	StatMain       []StatView
	StatSecondary  []StatView
	Flights        []FlightView
	CurrentYear    string
	AvailableYears []string
	IsFlightAdded  bool
}

type FlightView struct {
	ID               string
	Date             string
	TakeoffLocation  string
	TotalFlightTime  string
	MaxAltitude      string
	TotalThermicTime string
	MaxClimbRate     string
	Link             string
}

type FlightDetailStatView struct {
	Title string
	Value string
}
type FlightDetailView struct {
	FlightView
	UserView
	Stats []FlightDetailStatView
}

type StatView struct {
	Title            string
	AlltimeValue     string
	CurrentYearValue string
}

type ProgressionView struct {
	User UserView

	FlightTimeMonthlyData  []ChartDataset
	FlightCountMonthlyData []ChartDataset
}

type ChartDataset struct {
	Label string
	Data  []int
	Color string
}

type ChartData struct {
	Labels   []string
	Datasets []ChartDataset
}
