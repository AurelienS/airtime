package viewmodel

type LogbookView struct {
	StatMain      []StatView
	StatSecondary []StatView
	Flights       []FlightView
	Year          string
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
	Title       string
	Value       string
	Description string
}
