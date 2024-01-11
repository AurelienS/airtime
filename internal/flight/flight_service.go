package flight

import (
	"context"
	"time"

	flightstats "github.com/AurelienS/cigare/internal/flight_statistic"
	"github.com/AurelienS/cigare/internal/glider"
	"github.com/AurelienS/cigare/internal/user"
	"github.com/ezgliding/goigc/pkg/igc"
)

type Service struct {
	flightRepo    Repository
	gliderService glider.Service
}

func NewService(
	flightRepo Repository,
	gliderService glider.Service,
) Service {
	return Service{
		flightRepo:    flightRepo,
		gliderService: gliderService,
	}
}

func (s *Service) AddFlight(ctx context.Context, byteContent []byte, user user.User) error {
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
	Flights         []Flight
	Gliders         []glider.Glider
	TotalFlightTime time.Duration
	NumberOfFlight  int
}

func (s Service) GetDashboardData(ctx context.Context, user user.User) (DashboardData, error) {
	var data DashboardData

	flights, err := s.flightRepo.GetFlights(ctx, user)
	if err != nil {
		return data, err
	}

	gliders, err := s.gliderService.GetGliders(ctx, user)
	if err != nil {
		return data, err
	}

	totalFlightTime, err := s.flightRepo.GetTotalFlightTime(ctx, user.ID)
	if err != nil {
		return data, err
	}

	data = DashboardData{
		Flights:         flights,
		Gliders:         gliders,
		TotalFlightTime: totalFlightTime,
		NumberOfFlight:  len(flights),
	}
	return data, nil
}
