package service

import (
	"context"
	"time"

	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/internal/repository"
)

type FlightService struct {
	flightRepo repository.FlightRepository
}

func NewFlightService(
	flightRepo repository.FlightRepository,
) FlightService {
	return FlightService{
		flightRepo: flightRepo,
	}
}

func (s FlightService) GetFlights(
	ctx context.Context,
	start time.Time,
	end time.Time,
	user domain.User,
) ([]domain.Flight, error) {
	return s.flightRepo.GetFlights(ctx, start, end, user)
}

func (s FlightService) GetFlight(ctx context.Context, flightID int, user domain.User) (domain.Flight, error) {
	return s.flightRepo.GetFlight(ctx, flightID, user)
}

func (s FlightService) GetLastFlight(ctx context.Context, user domain.User) (*domain.Flight, error) {
	return s.flightRepo.GetLastFlight(ctx, user)
}

func (s FlightService) GetFlightsForDateRanges(
	flights []domain.Flight,
	dateRanges []domain.DateRange,
) [][]domain.Flight {
	flightsForRanges := make([][]domain.Flight, len(dateRanges))

	for _, flight := range flights {
		for i, dateRange := range dateRanges {
			if (flight.Date.Equal(dateRange.Start) || flight.Date.After(dateRange.Start)) &&
				(flight.Date.Equal(dateRange.End) || flight.Date.Before(dateRange.End)) {
				flightsForRanges[i] = append(flightsForRanges[i], flight)
			}
		}
	}

	return flightsForRanges
}

func (s FlightService) GetFlightsForYear(year int, flights []domain.Flight) []domain.Flight {
	dateRanges := []domain.DateRange{
		{
			Start: time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(year, time.December, 31, 23, 59, 59, 0, time.UTC),
		},
	}
	return s.GetFlightsForDateRanges(flights, dateRanges)[0]
}
