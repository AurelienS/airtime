package template

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

type LogbookView struct {
	Stats1  []StatView
	Stats2  []StatView
	Flights []FlightView
	Year    string
}
