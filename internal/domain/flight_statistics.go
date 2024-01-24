package domain

import (
	"time"
)

type MultipleFlightStats struct {
	Flights []Flight

	AltitudeMaxFlight Flight
	DistanceTotal     int
	DistanceMaxFlight Flight
	DurationMaxFlight Flight
	AverageDuration   time.Duration
	DurationTotal     time.Duration
}

func ComputeMultipleFlightStats(flights []Flight) MultipleFlightStats {
	var altitudeMaxFlight Flight
	var durationMaxFlight Flight
	averageDuration := time.Duration(0)
	var durationTotal time.Duration
	var totalDistance int
	var distanceMaxFlight Flight

	for _, f := range flights {
		if f.AltitudeMax > altitudeMaxFlight.AltitudeMax {
			altitudeMaxFlight = f
		}
		if f.Duration > durationMaxFlight.Duration {
			durationMaxFlight = f
		}

		if f.Distance > distanceMaxFlight.Distance {
			distanceMaxFlight = f
		}
		durationTotal += f.Duration
		totalDistance += f.Distance
	}

	flightCount := len(flights)
	if flightCount > 0 {
		averageDuration = durationTotal / time.Duration(flightCount)
	}

	aggregatedStats := MultipleFlightStats{
		AltitudeMaxFlight: altitudeMaxFlight,
		DurationMaxFlight: durationMaxFlight,
		DistanceMaxFlight: distanceMaxFlight,
		AverageDuration:   averageDuration,
		DurationTotal:     durationTotal,
		DistanceTotal:     totalDistance,
		Flights:           flights,
	}
	return aggregatedStats
}
