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
	Thermals       []*Thermal
	Stats          FlightStatistics
}

func (f *Flight) GenerateThermals(minRateOfClimb float64, maxDownwardTolerance int, minThermalDuration time.Duration, climbRateIntegrationPeriod int) {
	var current *Thermal
	var rateOfClimbHistory []float64

	for i, point := range f.Points {
		if i == 0 {
			continue
		}

		rateOfClimb := f.calculateRateOfClimb(i)
		rateOfClimbHistory = updateRateOfClimbHistory(rateOfClimbHistory, rateOfClimb, climbRateIntegrationPeriod)
		smoothedRateOfClimb := average(rateOfClimbHistory)

		if current != nil {
			current.Update(point, smoothedRateOfClimb)
			current = f.checkAndFinalizeThermal(current, maxDownwardTolerance, minThermalDuration, i)
		} else {
			current = f.maybeStartNewThermal(smoothedRateOfClimb, minRateOfClimb, climbRateIntegrationPeriod, point, i)
		}
	}

	f.finalizeLastThermal(current, minThermalDuration)
	f.Stats.Finalize(f)
}

func (f *Flight) calculateRateOfClimb(i int) float64 {
	altitudeGain := f.Points[i].GNSSAltitude - f.Points[i-1].GNSSAltitude
	timeElapsed := f.Points[i].Time.Sub(f.Points[i-1].Time).Seconds()
	return float64(altitudeGain) / timeElapsed
}

func updateRateOfClimbHistory(history []float64, rateOfClimb float64, period int) []float64 {
	history = append(history, rateOfClimb)
	if len(history) > period {
		history = history[1:]
	}
	return history
}

func (f *Flight) checkAndFinalizeThermal(current *Thermal, tolerance int, duration time.Duration, index int) *Thermal {
	if current.ShouldEnd(tolerance) {
		if current.Duration() >= duration {
			current.EndIndex = index
			current.AverageClimbRate = float64(current.Climb()) / current.Duration().Seconds()
			f.Thermals = append(f.Thermals, current)
			f.Stats.AddThermal(*current, duration)
		}
		return nil
	}
	return current
}

func (f *Flight) maybeStartNewThermal(smoothedRate float64, minRate float64, period int, point Point, index int) *Thermal {
	if smoothedRate >= minRate && len(f.Points) >= period {
		return NewThermal(point.Time, point.GNSSAltitude, index)
	}
	return nil
}

func (f *Flight) finalizeLastThermal(current *Thermal, duration time.Duration) {
	if current != nil && current.Duration() >= duration {
		current.EndIndex = len(f.Points) - 1
		f.Thermals = append(f.Thermals, current)
		f.Stats.AddThermal(*current, duration)
	}
}

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
