package repository

import (
	"context"
	"sort"
	"time"

	"github.com/AurelienS/cigare/internal/converter"
	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/internal/storage/ent"
	"github.com/AurelienS/cigare/internal/storage/ent/flight"
	userDB "github.com/AurelienS/cigare/internal/storage/ent/user"
	"github.com/AurelienS/cigare/internal/util"
)

type FlightRepository struct {
	client *ent.Client
}

func NewFlightRepository(client *ent.Client) FlightRepository {
	return FlightRepository{
		client: client,
	}
}

func (r FlightRepository) InsertFlight(
	ctx context.Context,
	flight domain.Flight,
	user domain.User,
) error {
	_, err := r.client.Flight.
		Create().
		SetDate(flight.Date).
		SetLocation(flight.Location).
		SetIgcData(flight.IgcData).
		SetAltitudeMax(flight.AltitudeMax).
		SetDistance(flight.Distance).
		SetDuration(int(flight.Duration.Seconds())).
		SetPilotID(user.ID).
		Save(ctx)

	util.Error().Err(err).Str("user", user.Email).Msg("Failed to insert flight into database")
	return err
}

func (r FlightRepository) InsertFlights(
	ctx context.Context,
	flights []domain.Flight,
	user domain.User,
) error {
	bulkFlights := make([]*ent.FlightCreate, len(flights))
	for i, flight := range flights {
		bulkFlights[i] = r.client.Flight.
			Create().
			SetDate(flight.Date).
			SetLocation(flight.Location).
			SetIgcData(flight.IgcData).
			SetAltitudeMax(flight.AltitudeMax).
			SetDistance(flight.Distance).
			SetDuration(int(flight.Duration.Seconds())).
			SetPilotID(user.ID)
	}
	err := r.client.Flight.
		CreateBulk(bulkFlights...).
		OnConflict().
		DoNothing().
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r FlightRepository) GetFlight(
	ctx context.Context,
	flightID int,
	user domain.User,
) (domain.Flight, error) {
	flightDB, err := r.client.User.
		Query().
		Where(userDB.IDEQ(user.ID)).
		QueryFlights().
		Where(flight.IDEQ(flightID)).
		WithPilot().
		First(ctx)
	if err != nil {
		return domain.Flight{}, err
	}

	return converter.DBToDomainFlight(flightDB), nil
}

func (r FlightRepository) GetFlights(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
	user domain.User,
) ([]domain.Flight, error) {
	util.Info().Str("user", user.Email).Times("dates", []time.Time{startDate, endDate}).Msg("Getting user flights")

	flightsDB, err := r.client.Flight.
		Query().
		Where(flight.HasPilotWith(userDB.IDEQ(user.ID))).
		Where(flight.DateGTE(startDate), flight.DateLTE(endDate)).
		WithPilot().
		Order(ent.Asc(flight.FieldDate)).
		All(ctx)
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to get flights")
		return nil, err
	}

	return converter.DBToDomainFlights(flightsDB), nil
}

func (r FlightRepository) GetLastFlights(ctx context.Context, count int, user domain.User) ([]domain.Flight, error) {
	util.Info().Str("user", user.Email).Int("flight count", count).Msg("Getting lasts flight")
	lastFlights, err := r.client.Flight.
		Query().
		Where(flight.HasPilotWith(userDB.IDEQ(user.ID))).
		Order(ent.Desc(flight.FieldDate)).
		WithPilot().
		Limit(count).
		All(ctx)
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed getting lasts flights")
		return nil, err
	}

	return converter.DBToDomainFlights(lastFlights), nil
}

func (r FlightRepository) GetFlyingYears(ctx context.Context, user domain.User) ([]int, error) {
	flights, err := r.client.Flight.
		Query().
		Where(flight.HasPilotWith(userDB.IDEQ(user.ID))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	yearSet := make(map[int]bool)
	for _, f := range flights {
		year := f.Date.Year()
		yearSet[year] = true
	}

	uniqueYears := make([]int, 0, len(yearSet))
	for year := range yearSet {
		uniqueYears = append(uniqueYears, year)
	}

	sort.Ints(uniqueYears)

	return uniqueYears, nil
}
