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

type FlightCount struct {
	Date  time.Time
	Count int
}

type FlightDuration struct {
	Date     time.Time
	Duration time.Duration
}

func reverse(s []domain.Flight) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func (s StatisticService) GetCumulativeFlightCount(
	ctx context.Context,
	user domain.User,
	start time.Time,
	end time.Time,
) ([]FlightCount, error) {
	flights, err := s.flightService.GetFlights(ctx, start, end, user)
	if err != nil {
		return []FlightCount{}, err
	}

	reverse(flights)

	type Year = int
	type MonthCount = map[time.Month]int
	type YearMonthCount map[Year]MonthCount
	flightCounts := make(YearMonthCount)
	var totalCount int
	for _, f := range flights {
		year := f.Date.Year()
		month := f.Date.Month()

		if _, ok := flightCounts[year]; !ok {
			flightCounts[year] = make(MonthCount)
		}

		totalCount++
		flightCounts[year][month] = totalCount
	}

	yearFlightCounts := make([]FlightCount, 0, len(flightCounts))
	for year, monthCount := range flightCounts {
		for month, count := range monthCount {
			yearFlightCounts = append(yearFlightCounts, FlightCount{
				Date:  time.Date(year, month, 1, 0, 0, 0, 0, time.UTC),
				Count: count,
			})
		}
	}

	sort.Slice(yearFlightCounts, func(i, j int) bool {
		return yearFlightCounts[i].Date.Before(yearFlightCounts[j].Date)
	})

	return yearFlightCounts, nil
}

func (s StatisticService) GetCumulativeFlightDuration(
	ctx context.Context,
	user domain.User,
	start time.Time,
	end time.Time,
) ([]FlightDuration, error) {
	flights, err := s.flightService.GetFlights(ctx, start, end, user)
	if err != nil {
		return []FlightDuration{}, err
	}

	// Assuming reverse function is defined elsewhere
	reverse(flights)

	type Year = int
	type MonthDuration = map[time.Month]time.Duration
	type YearMonthDuration map[Year]MonthDuration
	flightDurations := make(YearMonthDuration)
	var totalDuration time.Duration
	for _, f := range flights {
		year := f.Date.Year()
		month := f.Date.Month()

		if _, ok := flightDurations[year]; !ok {
			flightDurations[year] = make(MonthDuration)
		}

		totalDuration += f.Duration
		flightDurations[year][month] = totalDuration
	}

	yearFlightDurations := make([]FlightDuration, 0, len(flightDurations))
	for year, monthDuration := range flightDurations {
		for month, duration := range monthDuration {
			yearFlightDurations = append(yearFlightDurations, FlightDuration{
				Date:     time.Date(year, month, 1, 0, 0, 0, 0, time.UTC),
				Duration: duration,
			})
		}
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
