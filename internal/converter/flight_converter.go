package converter

import (
	"time"

	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/internal/storage/ent"
)

func DBToDomainFlight(flightDB *ent.Flight) domain.Flight {
	return domain.Flight{
		ID:          flightDB.ID,
		Date:        flightDB.Date,
		Location:    flightDB.Location,
		IgcData:     flightDB.IgcData,
		Duration:    time.Duration(flightDB.Duration) * time.Second,
		AltitudeMax: flightDB.AltitudeMax,
		Distance:    flightDB.Distance,
		Pilot:       DBToDomainUser(flightDB.Edges.Pilot),
	}
}

func DBToDomainFlights(flightDBs []*ent.Flight) []domain.Flight {
	var flights []domain.Flight
	for _, f := range flightDBs {
		flights = append(flights, DBToDomainFlight(f))
	}

	return flights
}
