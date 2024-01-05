package flight

import (
	"time"
)

type Flight struct {
	Manufacturer   string
	UniqueID       string
	AdditionalData string
	Date           time.Time
	Site           string
	Pilot          string
	ID             string
	Points         []Point
	Thermals       []Thermal
	Stats          ThermicStats
}

const (
	minRateOfClimb        = 2                // m/s, the threshold for considering it thermic activity
	minThermalDuration    = 40 * time.Second // Minimum duration to consider a sustained climb as thermal
	allowedDownwardPoints = 8                // Number of consecutive downward points allowed in a thermal
)

func (f *Flight) GenerateThermals() {
	f.generateThermals()
}

// generateThermals analyzes the flight points and identifies thermal segments.
func (f *Flight) generateThermals() {
	var inThermal bool
	var currentThermal Thermal

	for i := 1; i < len(f.Points); i++ {
		altitudeGain := f.Points[i].GNSSAltitude - f.Points[i-1].GNSSAltitude
		timeElapsed := f.Points[i].Time.Sub(f.Points[i-1].Time).Seconds()
		rateOfClimb := float64(altitudeGain) / timeElapsed

		if inThermal {
			currentThermal.Update(altitudeGain, rateOfClimb, f.Points[i].GNSSAltitude)
			if currentThermal.ShouldEnd() {
				inThermal = false
				currentThermal.End = f.Points[i].Time
				currentThermal.EndIndex = i
				f.Thermals = append(f.Thermals, currentThermal)
				f.Stats.AddThermal(currentThermal)
			}
		} else if rateOfClimb >= minRateOfClimb {
			inThermal = true
			currentThermal = NewThermal(f.Points[i].Time, f.Points[i].GNSSAltitude, i)
		}
	}

	if inThermal {
		f.Thermals = append(f.Thermals, currentThermal)
	}
	// f.Stats.Finalize(f.)
}

func (f *Flight) endT(thermal Thermal) {

}
