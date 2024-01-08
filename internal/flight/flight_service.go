package flight

import (
	"context"

	"github.com/AurelienS/cigare/internal/storage"
	"github.com/ezgliding/goigc/pkg/igc"
)

type FlightService struct {
	Repo FlightRepository
}

func NewFlightService(repository FlightRepository) FlightService {
	return FlightService{
		Repo: repository,
	}
}

func (s *FlightService) UploadFlight(ctx context.Context, byteContent []byte, user storage.User) error {
	track, err := igc.Parse(string(byteContent))
	if err != nil {
		return err
	}
	flight := ConvertToMyFlight(track)
	flight.UserID = user.ID

	return s.Repo.InsertFlight(ctx, flight, user)
}

func (s *FlightService) GetFlights(ctx context.Context, user storage.User) ([]storage.Flight, error) {
	return s.Repo.GetFlights(context.Background(), user)
}
