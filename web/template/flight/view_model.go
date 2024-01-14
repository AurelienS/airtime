package flight

import "github.com/AurelienS/cigare/internal/model"

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
	Gliders         []GliderView
	Flights         []FlightView
	Img             string
	User            model.User
	NumberOfFlight  string
	TotalFlightTime string
}
