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

type StatisticType int

const (
	MonthlyCount StatisticType = iota
	MonthlyDuration
	YearlyCount
	YearlyDuration
	CumulativeMonthlyCount
	CumulativeMonthlyDuration
)

type Statistics struct {
	MonthlyCount              map[int]map[time.Month]int
	MonthlyDuration           map[int]map[time.Month]time.Duration
	CumulativeMonthlyCount    map[int]map[time.Month]int
	CumulativeMonthlyDuration map[int]map[time.Month]time.Duration
	YearlyCount               map[int]int
	YearlyDuration            map[int]time.Duration
}

func initializeYearMonthStats(start, end time.Time) *Statistics {
	stats := &Statistics{
		MonthlyCount:              make(map[int]map[time.Month]int),
		MonthlyDuration:           make(map[int]map[time.Month]time.Duration),
		CumulativeMonthlyCount:    make(map[int]map[time.Month]int),
		CumulativeMonthlyDuration: make(map[int]map[time.Month]time.Duration),
		YearlyCount:               make(map[int]int),
		YearlyDuration:            make(map[int]time.Duration),
	}

	for y := start.Year(); y <= end.Year(); y++ {
		stats.MonthlyCount[y] = make(map[time.Month]int)
		stats.MonthlyDuration[y] = make(map[time.Month]time.Duration)
		stats.CumulativeMonthlyCount[y] = make(map[time.Month]int)
		stats.CumulativeMonthlyDuration[y] = make(map[time.Month]time.Duration)
		for m := time.January; m <= time.December; m++ {
			if y == start.Year() && m < start.Month() {
				continue
			}
			if y == end.Year() && m > end.Month() {
				break
			}
			stats.MonthlyCount[y][m] = 0
			stats.MonthlyDuration[y][m] = 0
			stats.CumulativeMonthlyCount[y][m] = 0
			stats.CumulativeMonthlyDuration[y][m] = 0
		}
		stats.YearlyCount[y] = 0
		stats.YearlyDuration[y] = 0
	}

	return stats
}

func (s StatisticService) ComputeStatistics(
	ctx context.Context,
	user domain.User,
	statTypes []StatisticType,
) (*Statistics, error) {
	flights, err := s.flightService.GetFlights(ctx, time.Time{}, time.Now(), user)
	if err != nil {
		return nil, err
	}

	sort.Slice(flights, func(i, j int) bool {
		return flights[i].Date.Before(flights[j].Date)
	})

	var firstDate, lastDate time.Time
	if len(flights) > 0 {
		firstDate = flights[0].Date
		lastDate = flights[len(flights)-1].Date
	} else {
		// Use the current year if there are no flights
		currentYear := time.Now().Year()
		firstDate = time.Date(currentYear, time.January, 1, 0, 0, 0, 0, time.UTC)
		lastDate = time.Date(currentYear, time.December, 31, 0, 0, 0, 0, time.UTC)
	}

	stats := initializeYearMonthStats(firstDate, lastDate)

	cumulativeCount := 0
	cumulativeDuration := time.Duration(0)

	for _, flight := range flights {
		year, month, _ := flight.Date.Date()

		// Update normal statistics
		stats.MonthlyCount[year][month]++
		stats.MonthlyDuration[year][month] += flight.Duration
		stats.YearlyCount[year]++
		stats.YearlyDuration[year] += flight.Duration

		// Update cumulative statistics
		cumulativeCount++
		cumulativeDuration += flight.Duration
		for yr := year; yr <= lastDate.Year(); yr++ {
			startMonth := time.January
			if yr == year {
				startMonth = month
			}
			for mth := startMonth; mth <= time.December; mth++ {
				stats.CumulativeMonthlyCount[yr][mth] = cumulativeCount
				stats.CumulativeMonthlyDuration[yr][mth] = cumulativeDuration
			}
		}
	}

	return stats, nil
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
