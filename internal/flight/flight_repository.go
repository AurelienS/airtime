package flight

import (
	"context"
	"time"

	flightstats "github.com/AurelienS/cigare/internal/flight_statistic"
	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/util"
)

type FlightRepository struct {
	queries storage.Queries
	tm      storage.TransactionManager
}

func NewFlightRepository(queries storage.Queries, tm storage.TransactionManager) FlightRepository {
	return FlightRepository{
		queries: queries,
		tm:      tm,
	}
}

func (r FlightRepository) InsertFlight(
	ctx context.Context,
	flight storage.Flight,
	flightStats flightstats.FlightStatistics,
	user storage.User,

) error {
	util.Info().Str("user", user.Email).Msg("Inserting flight")

	insertFlightParams := storage.InsertFlightParams{
		Date:            flight.Date,
		TakeoffLocation: flight.TakeoffLocation,
		UserID:          user.ID,
		GliderID:        flight.GliderID,
		IgcFilePath:     "not yet",
	}

	transaction := func() error {
		flightStatId, err := r.insertFlightStats(ctx, flightStats)
		if err != nil {
			util.Error().Err(err).Str("user", user.Email).Msg("Failed to insert flight statistics")
			return err
		}
		insertFlightParams.FlightStatisticsID = flightStatId

		_, err = r.queries.InsertFlight(ctx, insertFlightParams)
		if err != nil {
			util.Error().Err(err).Str("user", user.Email).Msg("Failed to insert flight")
		}
		return err
	}
	return r.tm.ExecuteTransaction(ctx, transaction)
}

func (r FlightRepository) GetFlights(ctx context.Context, user storage.User) ([]storage.Flight, error) {
	util.Info().Str("user", user.Email).Msg("Getting flights")

	flights, err := r.queries.GetFlights(ctx, user.ID)
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to get flights")
	}
	return flights, err
}

func (f FlightRepository) insertFlightStats(ctx context.Context, flightStat flightstats.FlightStatistics) (int32, error) {
	util.Info().Msg("Inserting flight statistics")

	param := storage.InsertFlightStatsParams{
		TotalThermicTime:  storage.DurationToPGInterval(flightStat.TotalThermicTime),
		TotalFlightTime:   storage.DurationToPGInterval(flightStat.TotalFlightTime),
		MaxClimb:          int32(flightStat.MaxClimb),
		MaxClimbRate:      flightStat.MaxClimbRate,
		TotalClimb:        int32(flightStat.TotalClimb),
		AverageClimbRate:  flightStat.AverageClimbRate,
		NumberOfThermals:  int32(flightStat.NumberOfThermals),
		PercentageThermic: flightStat.PercentageThermic,
		MaxAltitude:       int32(flightStat.MaxAltitude),
	}
	id, err := f.queries.InsertFlightStats(ctx, param)
	if err != nil {
		util.Error().Err(err).Msg("Failed to insert flight statistics")
		return 0, err
	}

	return id, nil
}

func (f FlightRepository) GetTotalFlightTime(ctx context.Context, userId int) (time.Duration, error) {
	util.Info().Int("user_id", userId).Msg("Getting total flight time")

	flightTimeMicroseconds, err := f.queries.GetTotalFlightTime(ctx, int32(userId))
	if err != nil {
		util.Error().Err(err).Int("user_id", userId).Msg("Failed to get total flight time")
		return 0, err
	}

	return time.Duration(flightTimeMicroseconds) * time.Microsecond, nil
}
