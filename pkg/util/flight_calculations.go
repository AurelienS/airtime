package util

import (
	"math"
)

func Average(numbers []float64) float64 {
	total := 0.0
	for _, number := range numbers {
		total += number
	}
	if len(numbers) == 0 {
		return 0
	}
	return total / float64(len(numbers))
}

func Calculatebearing(lat1, lat2, lng1, lng2 float64) float64 {
	lat1rad := Degtorad(lat1)
	lng1rad := Degtorad(lng1)
	lat2rad := Degtorad(lat2)
	lng2rad := Degtorad(lng2)

	deltaLon := lng2rad - lng1rad

	y := math.Sin(deltaLon) * math.Cos(lat2rad)
	x := math.Cos(lat1rad)*math.Sin(lat2rad) - math.Sin(lat1rad)*math.Cos(lat2rad)*math.Cos(deltaLon)

	initialBearing := Radtodeg(math.Atan2(y, x))

	// Normalize the bearing to be within [0, 360) degrees
	return math.Mod(initialBearing+360, 360)
}

// Degtorad converts degrees to radians.
func Degtorad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

// radToDeg converts radians to degrees.
func Radtodeg(rad float64) float64 {
	return rad * (180 / math.Pi)
}

func Updaterateofclimbhistory(history []float64, rateOfClimb float64, period int) []float64 {
	history = append(history, rateOfClimb)
	if len(history) > period {
		history = history[1:]
	}
	return history
}
