package flight

type GliderView struct {
	ID           int
	Name         string
	LinkToUpdate string
	IsSelected   bool
}

type FlightView struct {
	TakeoffLocation string
	Date            string
}

type DashboardView struct {
	Gliders []GliderView
	Flights []FlightView

	NumberOfFlight  string
	TotalFlightTime string
}
