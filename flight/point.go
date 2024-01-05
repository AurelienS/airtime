package flight

import (
	"time"
)

type Point struct {
	Lat          float64
	Lng          float64
	Bearing      float64
	Time         time.Time
	GNSSAltitude int
	// GroundAltitude int
}
