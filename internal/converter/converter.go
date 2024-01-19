package converter

import (
	"time"

	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/internal/storage/ent"
)

func DBToDomainUser(userDB *ent.User) model.User {
	return model.User{
		ID:         userDB.ID,
		GoogleID:   userDB.GoogleID,
		Email:      userDB.Email,
		Name:       userDB.Name,
		PictureURL: userDB.PictureURL,
	}
}

func DBToDomainUsers(userDB []*ent.User) []model.User {
	var users []model.User
	for _, user := range userDB {
		users = append(users, DBToDomainUser(user))
	}
	return users
}

func DBToDomainFlight(flightDB *ent.Flight) model.Flight {
	return model.Flight{
		ID:              flightDB.ID,
		Date:            flightDB.Date,
		TakeoffLocation: flightDB.TakeoffLocation,
		IgcFilePath:     flightDB.IgcFilePath,
		Pilot:           DBToDomainUser(flightDB.Edges.Pilot),
		Statistic:       DBToDomainFlightStatistic(flightDB.Edges.Statistic),
	}
}

func DBToDomainFlightStatistic(statDB *ent.FlightStatistic) model.FlightStatistic {
	var stat model.FlightStatistic

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

func DBToDomainFlightStatistics(statsDB []*ent.FlightStatistic) []model.FlightStatistic {
	var stats []model.FlightStatistic
	for _, s := range statsDB {
		stats = append(stats, DBToDomainFlightStatistic(s))
	}
	return stats
}
