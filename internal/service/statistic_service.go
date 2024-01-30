package service

import (
	"context"
	"sort"
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

func (s StatisticService) GetMonthlyStatisticsByYear(
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

type FlightCount struct {
	Date  time.Time
	Count int
}

type FlightDuration struct {
	Date     time.Time
	Duration time.Duration
}

func initializeFlightData(flights []domain.Flight) (map[time.Time]int, time.Time) {
	if len(flights) == 0 {
		return nil, time.Time{}
	}

	sort.Slice(flights, func(i, j int) bool {
		return flights[i].Date.Before(flights[j].Date)
	})

	firstFlight := flights[0].Date
	lastFlight := flights[len(flights)-1].Date

	flightData := make(map[time.Time]int)
	for d := firstFlight; !d.After(lastFlight); d = d.AddDate(0, 1, 0) {
		flightData[time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, time.UTC)] = 0
	}

	return flightData, lastFlight
}

// GetCumulativeFlightCount function with refactored logic.
func (s StatisticService) GetCumulativeFlightCount(
	ctx context.Context,
	user domain.User,
	start time.Time,
	end time.Time,
) ([]FlightCount, error) {
	flights, err := s.flightService.GetFlights(ctx, start, end, user)
	if err != nil {
		return nil, err
	}

	flightCounts, lastFlight := initializeFlightData(flights)
	if flightCounts == nil {
		return []FlightCount{}, nil
	}

	totalCount := 0
	for _, flight := range flights {
		totalCount++
		for d := flight.Date; !d.After(lastFlight); d = d.AddDate(0, 1, 0) {
			flightCounts[time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, time.UTC)] = totalCount
		}
	}

	yearFlightCounts := make([]FlightCount, 0, len(flightCounts))
	for date, count := range flightCounts {
		yearFlightCounts = append(yearFlightCounts, FlightCount{Date: date, Count: count})
	}

	sort.Slice(yearFlightCounts, func(i, j int) bool {
		return yearFlightCounts[i].Date.Before(yearFlightCounts[j].Date)
	})

	return yearFlightCounts, nil
}

// GetCumulativeFlightDuration function with refactored logic.
func (s StatisticService) GetCumulativeFlightDuration(
	ctx context.Context,
	user domain.User,
	start time.Time,
	end time.Time,
) ([]FlightDuration, error) {
	flights, err := s.flightService.GetFlights(ctx, start, end, user)
	if err != nil {
		return nil, err
	}

	flightDurationsMap, lastFlight := initializeFlightData(flights)
	if flightDurationsMap == nil {
		return []FlightDuration{}, nil
	}

	var totalDuration time.Duration
	for _, flight := range flights {
		totalDuration += flight.Duration
		for d := flight.Date; !d.After(lastFlight); d = d.AddDate(0, 1, 0) {
			flightDurationsMap[time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, time.UTC)] = int(totalDuration)
		}
	}

	yearFlightDurations := make([]FlightDuration, 0, len(flightDurationsMap))
	for date, duration := range flightDurationsMap {
		yearFlightDurations = append(yearFlightDurations, FlightDuration{Date: date, Duration: time.Duration(duration)})
	}

	sort.Slice(yearFlightDurations, func(i, j int) bool {
		return yearFlightDurations[i].Date.Before(yearFlightDurations[j].Date)
	})

	return yearFlightDurations, nil
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
