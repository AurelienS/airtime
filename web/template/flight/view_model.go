package flight

import (
	"github.com/AurelienS/cigare/internal/model"
)

type GliderView struct {
	ID           int
	Name         string
	LinkToUpdate string
	IsSelected   bool
}

//nolint:revive
type FlightView struct {
	Date             string
	TakeoffLocation  string
	TotalFlightTime  string
	MaxAltitude      string
	TotalThermicTime string
	MaxClimbRate     string
}

type DashboardView struct {
	Gliders         []GliderView
	Flights         []FlightView
	Img             string
	User            model.User
	NumberOfFlight  string
	TotalFlightTime string
}
