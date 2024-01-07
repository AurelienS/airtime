package service

import (
	"context"

	"github.com/AurelienS/cigare/internal/storage"
	repo "github.com/AurelienS/cigare/internal/storage/repository"
	"github.com/AurelienS/cigare/pkg/model"
	"github.com/ezgliding/goigc/pkg/igc"
)

type FlightService struct {
	Repo repo.FlightRepository
}

func (s *FlightService) UploadFlight(ctx context.Context, byteContent []byte, user storage.User) error {
	track, err := igc.Parse(string(byteContent))
	if err != nil {
		return err
	}
	flight := model.ConvertToMyFlight(track)
	flight.UserID = user.ID

	return s.Repo.InsertFlight(ctx, flight, user)
}

func (s *FlightService) GetFlights(ctx context.Context, user storage.User) ([]storage.Flight, error) {
	return s.Repo.GetFlights(context.Background(), user)
}
