package service

import (
	"context"
	"time"

	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/internal/repository"
)

type StatisticService struct {
	flightRepo repository.FlightRepository
}

func NewStatisticService(
	flightRepo repository.FlightRepository,
) StatisticService {
	return StatisticService{
		flightRepo: flightRepo,
	}
}

func (s StatisticService) GetStatisticsByYearAndMonth(
	ctx context.Context,
	user domain.User,
) (StatsYearMonth, error) {
	statsYearMonth := StatsYearMonth{}

	flights, err := s.flightRepo.GetFlights(ctx, time.Time{}, time.Now(), user)
	if err != nil {
		return statsYearMonth, err
	}

	// Prepare map to hold statistics by year and month
	flightsStatisticsByYearMonth := make(map[int]map[time.Month][]domain.Flight)
	for _, flight := range flights {
		year, month, _ := flight.Date.Date()

		// Initialize year and month if not already present
		if flightsStatisticsByYearMonth[year] == nil {
			flightsStatisticsByYearMonth[year] = make(map[time.Month][]domain.Flight)
			for m := time.January; m <= time.December; m++ {
				flightsStatisticsByYearMonth[year][m] = []domain.Flight{}
			}
		}
		flightsStatisticsByYearMonth[year][month] = append(flightsStatisticsByYearMonth[year][month], flight)
	}

	// Flatten the YearMonth stats to aggregated stats
	for year, monthStats := range flightsStatisticsByYearMonth {
		if statsYearMonth[year] == nil {
			statsYearMonth[year] = make(map[time.Month]StatsAggregated)
		}
		for month, stats := range monthStats {
			statsYearMonth[year][month] = s.ComputeAggregateStatistics(stats)
		}
	}

	return statsYearMonth, err
}

func (s StatisticService) GetFlyingYears(ctx context.Context, user domain.User) ([]int, error) {
	return s.flightRepo.GetFlyingYears(ctx, user)
}

type (
	Year           = int
	StatsYearMonth = map[Year]map[time.Month]StatsAggregated
)

type StatsAggregated struct {
	FlightCount           int
	MaxAltitude           int
	MaxClimb              int
	TotalClimb            int
	TotalNumberOfThermals int
	MaxClimbRate          float64
	MaxFlightLength       time.Duration
	MinFlightLength       time.Duration
	AverageFlightLength   time.Duration
	TotalFlightTime       time.Duration
	TotalThermicTime      time.Duration
}

func (s StatisticService) ComputeAggregateStatistics(flights []domain.Flight) StatsAggregated {
	var maxAltitude int
	var maxVario float64
	var maxFlightLength time.Duration
	minFlightLength := time.Duration(0)
	averageFlightLength := time.Duration(0)
	var totalFlightTime time.Duration

	var totalThermicTime time.Duration
	var maxClimb int
	var totalClimb int
	var totalNumberOfThermals int

	flightCount := len(flights)

	for _, f := range flights {
		if f.Statistic.MaxAltitude > maxAltitude {
			maxAltitude = f.Statistic.MaxAltitude
		}
		if f.Statistic.MaxClimbRate > maxVario {
			maxVario = f.Statistic.MaxClimbRate
		}
		if f.Statistic.TotalFlightTime > maxFlightLength {
			maxFlightLength = f.Statistic.TotalFlightTime
		}
		if f.Statistic.TotalFlightTime < minFlightLength {
			minFlightLength = f.Statistic.TotalFlightTime
		}
		if f.Statistic.MaxClimb > maxClimb {
			maxClimb = f.Statistic.MaxClimb
		}
		totalClimb += f.Statistic.TotalClimb
		totalNumberOfThermals += f.Statistic.NumberOfThermals
		totalThermicTime += f.Statistic.TotalThermicTime
		totalFlightTime += f.Statistic.TotalFlightTime
	}

	if flightCount > 0 {
		averageFlightLength = totalFlightTime / time.Duration(flightCount)
	}

	aggregatedStats := StatsAggregated{
		FlightCount:           flightCount,
		MaxAltitude:           maxAltitude,
		MaxClimb:              maxClimb,
		TotalClimb:            totalClimb,
		TotalNumberOfThermals: totalNumberOfThermals,
		MaxClimbRate:          maxVario,
		MaxFlightLength:       maxFlightLength,
		MinFlightLength:       minFlightLength,
		AverageFlightLength:   averageFlightLength,
		TotalFlightTime:       totalFlightTime,
		TotalThermicTime:      totalThermicTime,
	}
	return aggregatedStats
}
