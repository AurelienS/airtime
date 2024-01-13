package model

import (
	"time"

	flightstats "github.com/AurelienS/cigare/internal/flight_statistic"
)

type Flight struct {
	ID              int
	Date            time.Time
	TakeoffLocation string
	IgcFilePath     string
	Pilot           User
	Statistic       flightstats.FlightStatistic
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
