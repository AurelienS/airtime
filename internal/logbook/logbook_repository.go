package logbook

import (
	"context"
	"time"

	flightstats "github.com/AurelienS/cigare/internal/flight_statistic"
	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/internal/storage/ent"
	"github.com/AurelienS/cigare/internal/storage/ent/flight"
	userDb "github.com/AurelienS/cigare/internal/storage/ent/user"

	"github.com/AurelienS/cigare/internal/storage/ent/flightstatistic"
	"github.com/AurelienS/cigare/internal/util"
)

type Repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) Repository {
	return Repository{
		client: client,
	}
}

func (r Repository) InsertFlight(
	ctx context.Context,
	flight model.Flight,
	flightStats flightstats.FlightStatistic,
	user model.User,
) error {
	util.Info().Str("user", user.Email).Msg("Inserting flight")

	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}

	fs, err := tx.FlightStatistic.
		Create().
		SetTotalThermicTime(int(flightStats.TotalThermicTime.Seconds())).
		SetTotalFlightTime(int(flightStats.TotalFlightTime.Seconds())).
		SetMaxClimb(flightStats.MaxClimb).
		SetMaxClimbRate(flightStats.MaxClimbRate).
		SetTotalClimb(flightStats.TotalClimb).
		SetAverageClimbRate(flightStats.AverageClimbRate).
		SetNumberOfThermals(flightStats.NumberOfThermals).
		SetPercentageThermic(flightStats.PercentageThermic).
		SetMaxAltitude(flightStats.MaxAltitude).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Flight.
		Create().
		SetDate(flight.Date).
		SetTakeoffLocation(flight.TakeoffLocation).
		SetIgcFilePath("not yet").
		SetPilotID(user.ID).
		SetStatistic(fs).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetFlights(ctx context.Context, user model.User) ([]model.Flight, error) {
	util.Info().Str("user", user.Email).Msg("Getting user flights")

	// Fetch flights associated with the user
	flightsDB, err := r.client.User.
		Query().
		Where(userDb.IDEQ(user.ID)).
		QueryFlights().
		WithStatistic().
		All(ctx)
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to get flights")
		return nil, err
	}

	var flights []model.Flight
	for _, f := range flightsDB {
		flights = append(flights, model.DBToDomainFlight(f))
	}

	return flights, nil
}

func (r *Repository) GetTotalFlightTime(ctx context.Context, userID int) (time.Duration, error) {
	util.Info().Int("user_id", userID).Msg("Getting total flight time")

	// Perform the sum aggregation on the total flight time
	agg, err := r.client.FlightStatistic.
		Query().
		Where(flightstatistic.HasFlightWith(flight.HasPilotWith(userDb.IDEQ(userID)))).
		Aggregate(ent.Sum(flightstatistic.FieldTotalFlightTime)).
		Int(ctx)
	if err != nil {
		util.Error().Err(err).Int("user_id", userID).Msg("Failed to get total flight time")
		return 0, err
	}

	totalFlightTime := time.Duration(agg) * time.Second
	return totalFlightTime, nil
}
