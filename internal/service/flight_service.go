package service

import (
	"context"
	"time"

	"github.com/AurelienS/airtime/internal/domain"
	"github.com/AurelienS/airtime/internal/repository"
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

func (s FlightService) GetLastFlights(ctx context.Context, count int, user domain.User) ([]domain.Flight, error) {
	return s.flightRepo.GetLastFlights(ctx, count, user)
}

func (s FlightService) RemoveFlight(ctx context.Context, flightID int, user domain.User) error {
	return s.flightRepo.RemoveFlight(ctx, flightID, user)
}

func (s FlightService) RemoveAllFlightsOfUser(ctx context.Context, user domain.User) error {
	return s.flightRepo.RemoveAllFlightsOfUser(ctx, user)
}
