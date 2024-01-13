package flight

import (
	"context"
	"time"

	flightstats "github.com/AurelienS/cigare/internal/flight_statistic"
	"github.com/AurelienS/cigare/internal/model"
	"github.com/ezgliding/goigc/pkg/igc"
)

type Service struct {
	flightRepo Repository
}

func NewService(
	flightRepo Repository,
) Service {
	return Service{
		flightRepo: flightRepo,
	}
}

func (s *Service) AddFlight(ctx context.Context, byteContent []byte, user model.User) error {
	track, err := igc.Parse(string(byteContent))
	if err != nil {
		return err
	}

	flight := TrackToFlight(track, user)
	stats := flightstats.NewFlightStatistics(track.Points)
	err = s.flightRepo.InsertFlight(ctx, flight, stats, user)

	return err
}

type DashboardData struct {
	Flights         []model.Flight
	TotalFlightTime time.Duration
	NumberOfFlight  int
}

func (s Service) GetDashboardData(ctx context.Context, user model.User) (DashboardData, error) {
	var data DashboardData

	flights, err := s.flightRepo.GetFlights(ctx, user)
	if err != nil {
		return data, err
	}

	totalFlightTime, err := s.flightRepo.GetTotalFlightTime(ctx, user.ID)
	if err != nil {
		return data, err
	}

	data = DashboardData{
		Flights:         flights,
		TotalFlightTime: totalFlightTime,
		NumberOfFlight:  len(flights),
	}
	return data, nil
}
