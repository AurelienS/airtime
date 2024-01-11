package flightstats

import (
	"time"
)

type Thermal struct {
	Start             time.Time
	End               time.Time
	StartAltitude     int
	MaxAltitude       int
	DownwardTolerance int
	MaxClimbRate      float64
	AverageClimbRate  float64
	StartIndex        int
	EndIndex          int
}

func NewThermal(startTime time.Time, startAltitude int, startIndex int) *Thermal {
	return &Thermal{
		Start:         startTime,
		StartAltitude: startAltitude,
		StartIndex:    startIndex,
		MaxAltitude:   startAltitude,
	}
}

func (t *Thermal) Update(point Point, integratedClimbRate float64) {
	altitudeGain := point.GNSSAltitude - t.MaxAltitude
	if altitudeGain < 0 {
		t.DownwardTolerance++
	} else {
		t.DownwardTolerance = 0
		if point.GNSSAltitude > t.MaxAltitude {
			t.MaxAltitude = point.GNSSAltitude
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
	return t.MaxAltitude - t.StartAltitude
}
