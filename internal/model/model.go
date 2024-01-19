package model

import (
	"time"
)

type Flight struct {
	ID              int
	Date            time.Time
	TakeoffLocation string
	IgcFilePath     string
	Pilot           User
	Statistic       FlightStatistic
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Squad struct {
	ID        int
	Name      string
	CreatedAt time.Time
	Members   []User
}

type User struct {
	ID         int
	GoogleID   string
	Email      string
	Name       string
	PictureURL string
}

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

type (
	Year           = int
	StatsYearMonth = map[Year]map[time.Month]StatsAggregated
)
