package flight

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	flightstats "github.com/AurelienS/cigare/internal/flight_statistic"
	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/user"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository struct {
	queries storage.Queries
	tm      storage.TransactionManager
}

func NewRepository(queries storage.Queries, tm storage.TransactionManager) Repository {
	return Repository{
		queries: queries,
		tm:      tm,
	}
}

func ConvertFlightDBToFlight(flightDB storage.Flight) Flight {
	var flight Flight

	flight.ID = int(flightDB.ID)
	if flightDB.Date.Valid {
		flight.Date = flightDB.Date.Time
	}
	flight.TakeoffLocation = flightDB.TakeoffLocation
	flight.IgcFilePath = flightDB.IgcFilePath
	flight.UserID = int(flightDB.UserID)
	flight.FlightStatisticsID = int(flightDB.FlightStatisticsID)

	return flight
}

func ConvertFlightStatisticDBToFlightStatistic(statDB storage.FlightStatistic) flightstats.FlightStatistic {
	var stat flightstats.FlightStatistic

	stat.ID = int(statDB.ID)
	stat.TotalThermicTime = statDB.TotalThermicTime
	stat.TotalFlightTime = statDB.TotalFlightTime
	stat.MaxClimb = int(statDB.MaxClimb)
	stat.MaxClimbRate = statDB.MaxClimbRate
	stat.TotalClimb = int(statDB.TotalClimb)
	stat.AverageClimbRate = statDB.AverageClimbRate
	stat.NumberOfThermals = int(statDB.NumberOfThermals)
	stat.PercentageThermic = statDB.PercentageThermic
	stat.MaxAltitude = int(statDB.MaxAltitude)

	return stat
}

func (r Repository) InsertFlight(
	ctx context.Context,
	flight Flight,
	flightStats flightstats.FlightStatistic,
	user user.User,
) error {
	util.Info().Str("user", user.Email).Msg("Inserting flight")

	insertFlightParams := storage.InsertFlightParams{
		Date:            pgtype.Timestamptz{Valid: true, Time: flight.Date},
		TakeoffLocation: flight.TakeoffLocation,
		UserID:          int32(user.ID),
		IgcFilePath:     "not yet",
	}

	transaction := func() error {
		flightStatID, err := r.insertFlightStats(ctx, flightStats)
		if err != nil {
			util.Error().Err(err).Str("user", user.Email).Msg("Failed to insert flight statistics")
			return err
		}
		insertFlightParams.FlightStatisticsID = int32(flightStatID)

		_, err = r.queries.InsertFlight(ctx, insertFlightParams)
		if err != nil {
			util.Error().Err(err).Str("user", user.Email).Msg("Failed to insert flight")
		}
		return err
	}
	return r.tm.ExecuteTransaction(ctx, transaction)
}

func (r Repository) GetFlights(ctx context.Context, user user.User) ([]Flight, error) {
	util.Info().Str("user", user.Email).Msg("Getting flights")

	flightsDB, err := r.queries.GetFlights(ctx, int32(user.ID))
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to get flights")
	}

	var flights []Flight
	for _, f := range flightsDB {
		flights = append(flights, ConvertFlightDBToFlight(f))
	}

	return flights, err
}

func (r Repository) insertFlightStats(ctx context.Context, flightStat flightstats.FlightStatistic) (int, error) {
	util.Info().Msg("Inserting flight statistics")

	param := storage.InsertFlightStatsParams{
		TotalThermicTime:  flightStat.TotalThermicTime,
		TotalFlightTime:   flightStat.TotalFlightTime,
		MaxClimb:          int32(flightStat.MaxClimb),
		MaxClimbRate:      flightStat.MaxClimbRate,
		TotalClimb:        int32(flightStat.TotalClimb),
		AverageClimbRate:  flightStat.AverageClimbRate,
		NumberOfThermals:  int32(flightStat.NumberOfThermals),
		PercentageThermic: flightStat.PercentageThermic,
		MaxAltitude:       int32(flightStat.MaxAltitude),
	}
	id, err := r.queries.InsertFlightStats(ctx, param)
	if err != nil {
		util.Error().Err(err).Msg("Failed to insert flight statistics")
		return 0, err
	}

	return int(id), nil
}

func (r Repository) GetTotalFlightTime(ctx context.Context, userID int) (time.Duration, error) {
	util.Info().Int("user_id", userID).Msg("Getting total flight time")

	flightTimeMicroseconds, err := r.queries.GetTotalFlightTime(ctx, int32(userID))
	if err != nil {
		util.Error().Err(err).Int("user_id", userID).Msg("Failed to get total flight time")
		return 0, err
	}
	duration, err := parseHHMMSSDuration(flightTimeMicroseconds)
	if err != nil {
		util.Error().Err(err).Msg("Failed to parse sum(interval)::text")
		return 0, err
	}
	return duration, nil
}

func parseHHMMSSDuration(durationStr string) (time.Duration, error) {
	parts := strings.Split(durationStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid duration format")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hours: %w", err)
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes: %w", err)
	}

	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("invalid seconds: %w", err)
	}

	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
}
