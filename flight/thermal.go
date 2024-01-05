package flight

import "time"

// Thermal represents a period of thermic activity
type Thermal struct {
	Start             time.Time
	End               time.Time
	StartAltitude     int
	MaxAltitude       int
	downwardTolerance int
	ClimbRateTotal    float64
	ClimbRateCount    int
	StartIndex        int
	EndIndex          int
}

func NewThermal(startTime time.Time, startAltitude, startIndex int) Thermal {
	return Thermal{
		Start:         startTime,
		StartAltitude: startAltitude,
		MaxAltitude:   startAltitude,
		StartIndex:    startIndex,
	}
}

func (t *Thermal) Update(altitudeGain int, rateOfClimb float64, currentAltitude int) {
	if altitudeGain < 0 {
		t.downwardTolerance++
	} else {
		t.downwardTolerance = 0
		t.ClimbRateTotal += rateOfClimb
		t.ClimbRateCount++
	}

	if currentAltitude > t.MaxAltitude {
		t.MaxAltitude = currentAltitude
	}
}

func (t *Thermal) ShouldEnd() bool {
	return t.downwardTolerance > allowedDownwardPoints
}

func (t *Thermal) Duration() time.Duration {
	if t.End.IsZero() {
		// If the end time is not set, return 0 duration
		return 0
	}
	return t.End.Sub(t.Start)
}

func (t *Thermal) Climb() int {
	return t.MaxAltitude - t.StartAltitude
}

func (t *Thermal) AverageClimbRate() float64 {
	if t.ClimbRateCount == 0 {
		return 0
	}
	return t.ClimbRateTotal / float64(t.ClimbRateCount)
}
