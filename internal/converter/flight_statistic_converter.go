package converter

import (
	"time"

	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/internal/storage/ent"
)

func DBToDomainFlightStatistic(statDB *ent.FlightStatistic) domain.FlightStatistic {
	var stat domain.FlightStatistic

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
	stat.TotalDistance = statDB.TotalDistance
	stat.GeoJSON = statDB.GeoJSON

	return stat
}

func DBToDomainFlightStatistics(statsDB []*ent.FlightStatistic) []domain.FlightStatistic {
	var stats []domain.FlightStatistic
	for _, s := range statsDB {
		stats = append(stats, DBToDomainFlightStatistic(s))
	}
	return stats
}
