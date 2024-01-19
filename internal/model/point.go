package model

import (
	"time"

	"github.com/golang/geo/s2"
)

type Point struct {
	s2.LatLng
	Time             time.Time
	PressureAltitude int
	GNSSAltitude     int
	NumSatellites    int
	Description      string
}
