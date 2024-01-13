package flight

type GliderView struct {
	ID           int
	Name         string
	LinkToUpdate string
	IsSelected   bool
}

//nolint:revive
type FlightView struct {
	TakeoffLocation string
	Date            string
}

type DashboardView struct {
	Gliders []GliderView
	Flights []FlightView
	Img     string

	NumberOfFlight  string
	TotalFlightTime string
}
