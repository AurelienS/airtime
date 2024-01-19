package logbook

import (
	"context"
	"sort"
	"time"

	"github.com/AurelienS/cigare/internal/converter"
	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/internal/storage/ent"
	"github.com/AurelienS/cigare/internal/storage/ent/flight"
	"github.com/AurelienS/cigare/internal/storage/ent/flightstatistic"
	userDB "github.com/AurelienS/cigare/internal/storage/ent/user"
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
	flightStats model.FlightStatistic,
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
		r := tx.Rollback()
		if r != nil {
			return r
		}
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
		r := tx.Rollback()
		if r != nil {
			return r
		}
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetFlights(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
	user model.User,
) ([]model.Flight, error) {
	util.Info().Str("user", user.Email).Msg("Getting user flights")

	flightsDB, err := r.client.User.
		Query().
		Where(userDB.IDEQ(user.ID)).
		QueryFlights().
		Where(flight.DateGTE(startDate), flight.DateLTE(endDate)).
		WithStatistic().
		WithPilot().
		All(ctx)
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to get flights")
		return nil, err
	}

	var flights []model.Flight
	for _, f := range flightsDB {
		flights = append(flights, converter.DBToDomainFlight(f))
	}

	return flights, nil
}

func (r *Repository) GetStatistics(ctx context.Context,
	startDate time.Time,
	endDate time.Time,
	user model.User,
) ([]model.FlightStatistic, error) {
	util.Info().Str("user", user.Email).Msg("Getting statistics")

	stats, err := r.client.FlightStatistic.
		Query().
		Where(
			flightstatistic.HasFlightWith(
				flight.HasPilotWith(userDB.IDEQ(user.ID)),
				flight.DateGTE(startDate),
				flight.DateLTE(endDate),
			),
		).
		All(ctx)

	return converter.DBToDomainFlightStatistics(stats), err
}

// If there is no last flight, it return nil without an error.
func (r *Repository) GetLastFlight(ctx context.Context, user model.User) (*model.Flight, error) {
	util.Info().Str("user", user.Email).Msg("Getting last flight")

	flt, err := r.client.Flight.
		Query().
		Where(flight.HasPilotWith(userDB.IDEQ(user.ID))).
		Order(ent.Desc(flight.FieldDate)).
		WithPilot().
		WithStatistic().
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			util.Debug().Str("user", user.Email).Msg("No flight found")
			return nil, nil // Not a error but no flight
		} else {
			util.Error().Err(err).Str("user", user.Email).Msg("Failed querying user")
			return nil, err
		}
	}

	domainModel := converter.DBToDomainFlight(flt)
	return &domainModel, nil
}

func (r Repository) GetFlyingYears(ctx context.Context, user model.User) ([]int, error) {
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
