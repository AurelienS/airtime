package flight

import (
	"context"

	"github.com/AurelienS/cigare/internal/log"
	"github.com/AurelienS/cigare/internal/storage"
)

type FlightRepository interface {
	InsertFlight(ctx context.Context, flight storage.Flight, user storage.User) error
	GetFlights(ctx context.Context, user storage.User) ([]storage.Flight, error)
}

type SQLFlightRepository struct {
	Queries storage.Queries
}

func NewSQLFlightRepository(queries storage.Queries) SQLFlightRepository {
	return SQLFlightRepository{
		Queries: queries,
	}
}

func (repo SQLFlightRepository) InsertFlight(ctx context.Context, flight storage.Flight, user storage.User) error {

	params := storage.InsertFlightParams{
		Date:            flight.Date,
		TakeoffLocation: flight.TakeoffLocation,
		UserID:          user.ID,
		GliderID:        user.DefaultGliderID.Int32,
		IgcFilePath:     "not yet", // Placeholder path, replace with actual storage path as needed
	}
	err := repo.Queries.InsertFlight(context.Background(), params)
	if err != nil {
		log.Error().Err(err).Str("user", user.Email).Msg("Failed to insert flight into database")
		return err
	}
	return nil
}

func (repo SQLFlightRepository) GetFlights(ctx context.Context, user storage.User) ([]storage.Flight, error) {
	flights, err := repo.Queries.GetFlights(context.Background(), user.ID)
	if err != nil {
		log.Error().Err(err).Str("user", user.Email).Msg("Failed to get flights")
	}
	return flights, err
}
