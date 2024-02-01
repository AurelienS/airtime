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

func (s StatisticService) ComputeStatistics(
	ctx context.Context,
	user domain.User,
) (*domain.Statistics, error) {
	flights, err := s.flightService.GetFlights(ctx, time.Time{}, time.Now(), user)
	if err != nil {
		return nil, err
	}

	var firstDate, lastDate time.Time
	if len(flights) > 0 {
		firstDate = flights[0].Date
		lastDate = flights[len(flights)-1].Date
	} else {
		currentYear := time.Now().Year()
		firstDate = time.Date(currentYear, time.January, 1, 0, 0, 0, 0, time.UTC)
		lastDate = time.Date(currentYear, time.December, 31, 0, 0, 0, 0, time.UTC)
	}

	stats := addEveryDateInBetween(firstDate, lastDate)

	flightIndex := 0
	cumulativeCount := 0
	cumulativeDuration := time.Duration(0)

	for i := range stats.MonthlyCount {
		for flightIndex < len(flights) && flights[flightIndex].Date.Before(stats.MonthlyCount[i].Date.AddDate(0, 1, 0)) {
			flight := flights[flightIndex]
			stats.MonthlyCount[i].Count++
			stats.MonthlyDuration[i].Duration += flight.Duration
			cumulativeCount++
			cumulativeDuration += flight.Duration

			yearIndex := flight.Date.Year() - firstDate.Year()
			stats.YearlyCount[yearIndex].Count++
			stats.YearlyDuration[yearIndex].Duration += flight.Duration

			flightIndex++
		}

		stats.CumulativeMonthlyCount = append(
			stats.CumulativeMonthlyCount,
			domain.DateCount{Date: stats.MonthlyCount[i].Date, Count: cumulativeCount},
		)
		stats.CumulativeMonthlyDuration = append(
			stats.CumulativeMonthlyDuration,
			domain.DateDuration{Date: stats.MonthlyDuration[i].Date, Duration: cumulativeDuration},
		)
	}

	return stats, nil
}

func addEveryDateInBetween(start, end time.Time) *domain.Statistics {
	stats := &domain.Statistics{}

	for year := start.Year(); year <= end.Year(); year++ {
		stats.YearlyCount = append(
			stats.YearlyCount,
			domain.DateCount{Date: time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC), Count: 0},
		)
		stats.YearlyDuration = append(
			stats.YearlyDuration,
			domain.DateDuration{Date: time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC), Duration: 0},
		)
	}

	startDate := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC)
	for date := startDate; !date.After(end); date = date.AddDate(0, 1, 0) {
		stats.MonthlyCount = append(stats.MonthlyCount, domain.DateCount{Date: date, Count: 0})
		stats.MonthlyDuration = append(stats.MonthlyDuration, domain.DateDuration{Date: date, Duration: 0})
	}

	return stats
}

func (s StatisticService) GetFlyingYears(ctx context.Context, user domain.User) ([]int, error) {
	return s.flightRepo.GetFlyingYears(ctx, user)
}

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
