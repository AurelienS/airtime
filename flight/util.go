package flight

import "math"

func average(numbers []float64) float64 {
	total := 0.0
	for _, number := range numbers {
		total += number
	}
	if len(numbers) == 0 {
		return 0
	}
	return total / float64(len(numbers))
}

func calculateBearing(p1, p2 Point) float64 {
	lat1 := degToRad(p1.Lat)
	lon1 := degToRad(p1.Lng)
	lat2 := degToRad(p2.Lat)
	lon2 := degToRad(p2.Lng)

	deltaLon := lon2 - lon1

	y := math.Sin(deltaLon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(deltaLon)

	initialBearing := radToDeg(math.Atan2(y, x))

	// Normalize the bearing to be within [0, 360) degrees
	return math.Mod(initialBearing+360, 360)
}

// degToRad converts degrees to radians.
func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

// radToDeg converts radians to degrees.
func radToDeg(rad float64) float64 {
	return rad * (180 / math.Pi)
}

func updateRateOfClimbHistory(history []float64, rateOfClimb float64, period int) []float64 {
	history = append(history, rateOfClimb)
	if len(history) > period {
		history = history[1:]
	}
	return history
}
