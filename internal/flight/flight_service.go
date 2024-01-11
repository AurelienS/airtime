package flight

import (
	"context"
	"time"

	flightstats "github.com/AurelienS/cigare/internal/flight_statistic"
	"github.com/AurelienS/cigare/internal/storage"
	"github.com/ezgliding/goigc/pkg/igc"
)

type FlightService struct {
	flightRepo FlightRepository
}

func NewFlightService(
	flightRepo FlightRepository,

) FlightService {
	return FlightService{
		flightRepo: flightRepo,
	}
}

func (s *FlightService) AddFlight(ctx context.Context, byteContent []byte, user storage.User) error {
	track, err := igc.Parse(string(byteContent))
	if err != nil {
		return err
	}

	flight := TrackToFlight(track, user)
	stats := flightstats.NewFlightStatistics(track.Points)
	err = s.flightRepo.InsertFlight(ctx, flight, stats, user)

	return err
}

func (s *FlightService) GetFlights(ctx context.Context, user storage.User) ([]storage.Flight, error) {
	return s.flightRepo.GetFlights(ctx, user)
}

func (s FlightService) GetTotalFlightTime(ctx context.Context, userId int) (time.Duration, error) {
	totalTime, err := s.flightRepo.GetTotalFlightTime(ctx, userId)
	if err != nil {
		return 0, err
	}

	return totalTime, nil
}
