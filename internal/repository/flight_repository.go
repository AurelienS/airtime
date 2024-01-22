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
	flightStats domain.FlightStatistic,
	user domain.User,
) error {
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
		SetTotalDistance(int(flightStats.TotalDistance)).
		Save(ctx)
	if err != nil {
		r := tx.Rollback()
		if r != nil {
			util.
				Error().
				Err(err).
				Str("user", user.Email).
				Msg("Failed to insert flight statistic into database AND failed to rollback transaction")
			return r
		}
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to insert statistic flight into database")
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
			util.
				Error().
				Err(err).
				Str("user", user.Email).
				Msg("Failed to insert flight into database AND failed to rollback transaction")
			return r
		}
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to insert flight into database")
		return err
	}

	return tx.Commit()
}

func (r FlightRepository) InsertFlights(
	ctx context.Context,
	flights []domain.Flight,
	flightStats []domain.FlightStatistic,
	user domain.User,
) error {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}

	// Batch insert flight statistics
	bulkStats := make([]*ent.FlightStatisticCreate, len(flightStats))
	for i, flightStat := range flightStats {
		bulkStats[i] = tx.FlightStatistic.
			Create().
			SetTotalThermicTime(int(flightStat.TotalThermicTime.Seconds())).
			SetTotalFlightTime(int(flightStat.TotalFlightTime.Seconds())).
			SetMaxClimb(flightStat.MaxClimb).
			SetMaxClimbRate(flightStat.MaxClimbRate).
			SetTotalClimb(flightStat.TotalClimb).
			SetAverageClimbRate(flightStat.AverageClimbRate).
			SetNumberOfThermals(flightStat.NumberOfThermals).
			SetPercentageThermic(flightStat.PercentageThermic).
			SetMaxAltitude(flightStat.MaxAltitude).
			SetTotalDistance(int(flightStat.TotalDistance))
	}
	stats, err := tx.FlightStatistic.
		CreateBulk(bulkStats...).Save(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Batch insert flights
	bulkFlights := make([]*ent.FlightCreate, len(flights))
	for i, flight := range flights {
		bulkFlights[i] = tx.Flight.
			Create().
			SetDate(flight.Date).
			SetTakeoffLocation(flight.TakeoffLocation).
			SetIgcFilePath("not yet").
			SetPilotID(user.ID).
			SetStatistic(stats[i])
	}
	err = tx.Flight.
		CreateBulk(bulkFlights...).
		OnConflict().
		DoNothing().
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
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
		WithStatistic().
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
		WithStatistic().
		Order(ent.Desc(flight.FieldDate)).
		All(ctx)
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to get flights")
		return nil, err
	}

	return converter.DBToDomainFlights(flightsDB), nil
}

// If there is no last flight, it return nil without an error.
func (r FlightRepository) GetLastFlight(ctx context.Context, user domain.User) (*domain.Flight, error) {
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
			return nil, nil //nolint:nilnil
		}
		util.Error().Err(err).Str("user", user.Email).Msg("Failed querying user")
		return nil, err
	}

	domainModel := converter.DBToDomainFlight(flt)
	return &domainModel, nil
}

func (r FlightRepository) GetLastFlights(ctx context.Context, count int, user domain.User) ([]domain.Flight, error) {
	util.Info().Str("user", user.Email).Int("flight count", count).Msg("Getting lasts flight")
	lastFlights, err := r.client.Flight.
		Query().
		Where(flight.HasPilotWith(userDB.IDEQ(user.ID))).
		Order(ent.Desc(flight.FieldDate)).
		WithPilot().
		WithStatistic().
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
