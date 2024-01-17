package viewmodel

type LogbookView struct {
	StatMain       []StatView
	StatSecondary  []StatView
	Flights        []FlightView
	CurrentYear    string
	AvailableYears []string
}

type FlightView struct {
	Date             string
	TakeoffLocation  string
	TotalFlightTime  string
	MaxAltitude      string
	TotalThermicTime string
	MaxClimbRate     string
}

type StatView struct {
	Title            string
	AlltimeValue     string
	CurrentYearValue string
}
