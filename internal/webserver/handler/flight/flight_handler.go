package flight

import (
	"github.com/AurelienS/cigare/internal/service"
)

type FlightHandler struct {
	FlightService service.FlightService
	GliderService service.GliderService
}

func NewFlightHandler(flightService service.FlightService, gliderService service.GliderService) *FlightHandler {
	return &FlightHandler{
		FlightService: flightService,
		GliderService: gliderService,
	}
}
