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

func (f *Flight) Initialize() {
	const (
		minClimbRate               = 0.2              // m/s, the threshold for considering it thermic activity
		climbRateIntegrationPeriod = 10               // Number of seconds to smooth the climbRate
		minThermalDuration         = 20 * time.Second // Minimum duration to consider a sustained climb as thermal
		allowedDownwardPoints      = 4                // Number of consecutive downward points allowed in a thermal
	)

	// demReader := NewDEMReader("30n000e_20101117_gmted_std300.tif")
	// elevations, _ := demReader.GetElevations(*f)
	// fmt.Println("file: flight.go ~ line 31 ~ elevations : ", elevations)

	// for i := range f.Points {
	// 	f.Points[i].GroundAltitude = int(math.Round(elevations[i]))
	// }

	f.calculateBearings()
	f.GenerateThermals(minClimbRate, allowedDownwardPoints, minThermalDuration, climbRateIntegrationPeriod)
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

func (f *Flight) calculateBearings() {
	for i := 0; i < len(f.Points)-1; i++ {
		f.Points[i].Bearing = calculateBearing(f.Points[i], f.Points[i+1])
	}

	if len(f.Points) > 0 {
		f.Points[len(f.Points)-1].Bearing = 0 // Assign a default value for the last point
	}
}

func (f *Flight) calculateRateOfClimb(i int) float64 {
	altitudeGain := f.Points[i].GNSSAltitude - f.Points[i-1].GNSSAltitude
	timeElapsed := f.Points[i].Time.Sub(f.Points[i-1].Time).Seconds()
	return float64(altitudeGain) / timeElapsed
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
