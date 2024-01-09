package flight

import (
	"context"

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
			return err
		}
		insertFlightParams.FlightStatisticsID = flightStatId

		_, err = r.queries.InsertFlight(ctx, insertFlightParams)
		return err
	}
	return r.tm.ExecuteTransaction(ctx, transaction)
}

func (r FlightRepository) GetFlights(ctx context.Context, user storage.User) ([]storage.Flight, error) {
	flights, err := r.queries.GetFlights(ctx, user.ID)
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to get flights")
	}
	return flights, err
}

func (f FlightRepository) insertFlightStats(ctx context.Context, flightStat flightstats.FlightStatistics) (int32, error) {
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
		return 0, err
	}

	return id, nil
}
