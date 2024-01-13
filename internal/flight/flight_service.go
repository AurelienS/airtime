package flight

import (
	"context"
	"io"
	"mime/multipart"
	"time"

	flightstats "github.com/AurelienS/cigare/internal/flight_statistic"
	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/internal/util"
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

func (s *Service) ProcessAndAddFlight(ctx context.Context, file *multipart.FileHeader, user model.User) error {
	src, err := file.Open()
	if err != nil {
		util.Error().Err(err).Str("filename", file.Filename).Msg("Failed to open IGC file")
		return err
	}
	defer src.Close()

	byteContent, err := io.ReadAll(src)
	if err != nil {
		util.Error().Err(err).Str("filename", file.Filename).Msg("Failed to read IGC file")
		return err
	}

	track, err := igc.Parse(string(byteContent))
	if err != nil {
		return err
	}

	flight := TrackToFlight(track, user)
	stats := flightstats.NewFlightStatistics(track.Points)
	err = s.flightRepo.InsertFlight(ctx, flight, stats, user)
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to insert flight into database")
		return err
	}

	util.Info().Str("user", user.Email).Str("filename", file.Filename).
		Msg("File processed and flight record created successfully")

	return nil
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
