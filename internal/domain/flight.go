package domain

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

type Flight struct {
	ID              int
	Date            time.Time
	TakeoffLocation string
	IgcFilePath     string
	Pilot           User
	Statistic       FlightStatistic
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Thermal struct {
	Start              time.Time
	End                time.Time
	StartAltitude      int
	MaxThermalAltitude int
	DownwardTolerance  int
	MaxClimbRate       float64
	AverageClimbRate   float64
	StartIndex         int
	EndIndex           int
}

func NewThermal(startTime time.Time, startAltitude, startIndex int) *Thermal {
	return &Thermal{
		Start:              startTime,
		StartAltitude:      startAltitude,
		StartIndex:         startIndex,
		MaxThermalAltitude: startAltitude,
	}
}

func (t *Thermal) Update(point Point, integratedClimbRate float64) {
	altitudeGain := point.GNSSAltitude - t.MaxThermalAltitude
	if altitudeGain < 0 {
		t.DownwardTolerance++
	} else {
		t.DownwardTolerance = 0
		if point.GNSSAltitude > t.MaxThermalAltitude {
			t.MaxThermalAltitude = point.GNSSAltitude
		}
		if integratedClimbRate > t.MaxClimbRate {
			t.MaxClimbRate = integratedClimbRate
		}
	}
	t.End = point.Time
}

func (t *Thermal) ShouldEnd(maxDownwardTolerance int) bool {
	return t.DownwardTolerance > maxDownwardTolerance
}

func (t *Thermal) Duration() time.Duration {
	if t.End.IsZero() {
		return 0
	}
	return t.End.Sub(t.Start)
}

func (t *Thermal) Climb() int {
	return t.MaxThermalAltitude - t.StartAltitude
}
