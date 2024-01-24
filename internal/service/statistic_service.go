package service

import (
	"context"
	"time"

	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/internal/repository"
)

type StatisticService struct {
	flightRepo    repository.FlightRepository
	flightService FlightService
}

func NewStatisticService(
	flightRepo repository.FlightRepository,
	flightService FlightService,
) StatisticService {
	return StatisticService{
		flightRepo:    flightRepo,
		flightService: flightService,
	}
}

func (s StatisticService) GetStatisticsByYearAndMonth(
	ctx context.Context,
	user domain.User,
) (StatsYearMonth, error) {
	statsYearMonth := StatsYearMonth{}

	flights, err := s.flightService.GetFlights(ctx, time.Time{}, time.Now(), user)
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
			statsYearMonth[year] = make(map[time.Month]domain.MultipleFlightStats)
		}
		for month, stats := range monthStats {
			statsYearMonth[year][month] = domain.ComputeMultipleFlightStats(stats)
		}
	}

	return statsYearMonth, err
}

func (s StatisticService) GetFlyingYears(ctx context.Context, user domain.User) ([]int, error) {
	return s.flightRepo.GetFlyingYears(ctx, user)
}

type (
	Year           = int
	StatsYearMonth = map[Year]map[time.Month]domain.MultipleFlightStats
)

func (s StatisticService) GetFlightStats(
	ctx context.Context,
	user domain.User,
	start time.Time,
	end time.Time,
) (domain.MultipleFlightStats, error) {
	flights, err := s.flightService.GetFlights(ctx, start, end, user)
	var stats domain.MultipleFlightStats
	if err != nil {
		return stats, err
	}
	stats = domain.ComputeMultipleFlightStats(flights)
	return stats, nil
}
