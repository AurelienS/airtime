package storage

import (
	"context"

	"github.com/AurelienS/cigare/internal/storage"
)

type FlightRepository interface {
	InsertFlight(ctx context.Context, flight storage.Flight, user storage.User) error
	GetFlights(ctx context.Context, user storage.User) ([]storage.Flight, error)
}

type GliderRepository interface {
	GetGliders(ctx context.Context, user storage.User) ([]storage.Glider, error)
}
