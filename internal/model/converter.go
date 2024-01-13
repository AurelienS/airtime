package model

import (
	"time"

	flightstats "github.com/AurelienS/cigare/internal/flight_statistic"
	"github.com/AurelienS/cigare/internal/storage/ent"
)

func DBToDomainUser(userDB *ent.User) User {
	return User{
		ID:         userDB.ID,
		GoogleID:   userDB.GoogleID,
		Email:      userDB.Email,
		Name:       userDB.Name,
		PictureURL: userDB.PictureURL,
	}
}

func DBToDomainUsers(userDB []*ent.User) []User {
	var users []User
	for _, user := range userDB {
		users = append(users, DBToDomainUser(user))
	}
	return users
}

func DBToDomainSquad(squadDB *ent.Squad) Squad {
	return Squad{
		ID:        squadDB.ID,
		Name:      squadDB.Name,
		CreatedAt: squadDB.CreatedAt,
		Members:   DBToDomainUsers(squadDB.Edges.Members),
	}
}

func DBToDomainSquads(squadsDB []*ent.Squad) []Squad {
	var squads []Squad
	for _, squad := range squadsDB {
		squads = append(squads, DBToDomainSquad(squad))
	}
	return squads
}

func DBToDomainFlight(flightDB *ent.Flight) Flight {
	return Flight{
		ID:              flightDB.ID,
		Date:            flightDB.Date,
		TakeoffLocation: flightDB.TakeoffLocation,
		IgcFilePath:     flightDB.IgcFilePath,
		Pilot:           DBToDomainUser(flightDB.Edges.Pilot),
		Statistic:       DBToDomainFlightStatistic(flightDB.Edges.Statistic),
	}
}

func DBToDomainFlightStatistic(statDB *ent.FlightStatistic) flightstats.FlightStatistic {
	var stat flightstats.FlightStatistic

	stat.ID = statDB.ID
	stat.TotalThermicTime = time.Duration(statDB.TotalThermicTime) * time.Second
	stat.TotalFlightTime = time.Duration(statDB.TotalFlightTime) * time.Second
	stat.MaxClimb = statDB.MaxClimb
	stat.MaxClimbRate = statDB.MaxClimbRate
	stat.TotalClimb = statDB.TotalClimb
	stat.AverageClimbRate = statDB.AverageClimbRate
	stat.NumberOfThermals = statDB.NumberOfThermals
	stat.PercentageThermic = statDB.PercentageThermic
	stat.MaxAltitude = statDB.MaxAltitude

	return stat
}
