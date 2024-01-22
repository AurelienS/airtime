package converter

import (
	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/internal/storage/ent"
)

func DBToDomainFlight(flightDB *ent.Flight) domain.Flight {
	return domain.Flight{
		ID:              flightDB.ID,
		Date:            flightDB.Date,
		TakeoffLocation: flightDB.TakeoffLocation,
		IgcFilePath:     flightDB.IgcFilePath,
		Pilot:           DBToDomainUser(flightDB.Edges.Pilot),
		Statistic:       DBToDomainFlightStatistic(flightDB.Edges.Statistic),
	}
}

func DBToDomainFlights(flightDBs []*ent.Flight) []domain.Flight {
	var flights []domain.Flight
	for _, f := range flightDBs {
		flights = append(flights, DBToDomainFlight(f))
	}

	return flights
}
