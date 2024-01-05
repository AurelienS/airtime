package enhancer

import (
	"math"

	"github.com/AurelienS/cigare/flight"
)

func EnhanceWithBearing(flight *flight.Flight) {
	points := flight.Points

	for i := 0; i < len(points)-1; i++ {
		// Calculate bearing between two consecutive points
		bearing := calculateBearing(points[i], points[i+1])

		// Assign the calculated bearing to the current point
		points[i].Bearing = bearing
	}

	// Assign the last point's bearing as 0 or any default value since there is no next point
	if len(points) > 0 {
		points[len(points)-1].Bearing = 0
	}
}

func calculateBearing(p1 flight.Point, p2 flight.Point) float64 {
	lat1 := degToRad(p1.Lat)
	lon1 := degToRad(p1.Lng)
	lat2 := degToRad(p2.Lat)
	lon2 := degToRad(p2.Lng)

	deltaLon := lon2 - lon1

	y := math.Sin(deltaLon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - (math.Sin(lat1) * math.Cos(lat2) * math.Cos(deltaLon))

	initialBearing := math.Atan2(y, x)
	initialBearing = radToDeg(initialBearing)

	// Normalize the bearing to be in the range [0, 360)
	initialBearing = math.Mod((initialBearing + 360), 360)

	return initialBearing
}

// degToRad converts degrees to radians
func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

// radToDeg converts radians to degrees
func radToDeg(rad float64) float64 {
	return rad * (180 / math.Pi)
}
