package flightstats

import (
	"fmt"
	"time"

	"github.com/ezgliding/goigc/pkg/igc"
)

type FlightStatistic struct {
	ID                int
	TotalThermicTime  time.Duration
	TotalFlightTime   time.Duration
	MaxClimb          int
	MaxClimbRate      float64
	TotalClimb        int
	AverageClimbRate  float64
	NumberOfThermals  int
	PercentageThermic float64
	MaxAltitude       int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Thermals          []*Thermal
	Points            []Point
}

func NewFlightStatistics(Points []igc.Point) FlightStatistic {

	var pts []Point
	for _, p := range Points {
		pts = append(pts, Point{
			LatLng:           p.LatLng,
			Time:             p.Time,
			PressureAltitude: int(p.PressureAltitude),
			GNSSAltitude:     int(p.GNSSAltitude),
			NumSatellites:    p.NumSatellites,
			Description:      p.Description,
		})
	}

	stat := FlightStatistic{
		Points: pts,
	}
	stat.Compute()
	return stat
}

func (f *FlightStatistic) Compute() {
	const (
		minClimbRate               = 0.2              // m/s, the threshold for considering it thermic activity
		climbRateIntegrationPeriod = 10               // Number of seconds to smooth the climbRate
		minThermalDuration         = 20 * time.Second // Minimum duration to consider a sustained climb as thermal
		allowedDownwardPoints      = 4                // Number of consecutive downward Points allowed in a thermal
	)

	var current *Thermal
	var rateOfClimbHistory []float64

	for i, point := range f.Points {
		if i == 0 {
			continue
		}

		smoothedRateOfClimb := f.calculateSmoothedRateOfClimb(i, climbRateIntegrationPeriod, rateOfClimbHistory)

		if current == nil {
			current = f.maybeStartNewThermal(smoothedRateOfClimb, minClimbRate, climbRateIntegrationPeriod, point, i)
		} else {
			current.Update(point, smoothedRateOfClimb)
			current = f.checkAndFinalizeThermal(current, allowedDownwardPoints, minThermalDuration, i)
		}
	}

	f.finalizeLastThermal(current, minThermalDuration)
	f.Finalize()
}

func (f *FlightStatistic) calculateSmoothedRateOfClimb(i, period int, rateOfClimbHistory []float64) float64 {
	rateOfClimb := f.calculateRateOfClimb(i)
	rateOfClimbHistory = UpdateRateOfClimbHistory(rateOfClimbHistory, rateOfClimb, period)
	return Average(rateOfClimbHistory)
}

func (f *FlightStatistic) calculateRateOfClimb(i int) float64 {
	altitudeGain := f.Points[i].GNSSAltitude - f.Points[i-1].GNSSAltitude
	timeElapsed := f.Points[i].Time.Sub(f.Points[i-1].Time).Seconds()
	return float64(altitudeGain) / timeElapsed
}

func (f *FlightStatistic) checkAndFinalizeThermal(current *Thermal, tolerance int, duration time.Duration, index int) *Thermal {
	if current.ShouldEnd(tolerance) {
		if current.Duration() >= duration {
			current.EndIndex = index
			current.AverageClimbRate = float64(current.Climb()) / current.Duration().Seconds()
			f.Thermals = append(f.Thermals, current)
			f.AddThermal(*current, duration)
		}
		return nil
	}
	return current
}

func (f *FlightStatistic) maybeStartNewThermal(smoothedRate float64, minRate float64, period int, point Point, index int) *Thermal {
	if smoothedRate >= minRate && len(f.Points) >= period {
		return NewThermal(point.Time, point.GNSSAltitude, index)
	}
	return nil
}

func (f *FlightStatistic) finalizeLastThermal(current *Thermal, duration time.Duration) {
	if current != nil && current.Duration() >= duration {
		current.EndIndex = len(f.Points) - 1
		f.Thermals = append(f.Thermals, current)
		f.AddThermal(*current, duration)
	}
}

func (s *FlightStatistic) AddThermal(t Thermal, minThermalDuration time.Duration) {
	duration := t.Duration()
	if duration < minThermalDuration || t.Climb() <= 0 {
		return
	}

	s.NumberOfThermals++
	s.TotalThermicTime += duration

	climb := t.Climb()
	s.TotalClimb += climb
	if climb > s.MaxClimb {
		s.MaxClimb = climb
	}
	if t.MaxAltitude > s.MaxAltitude {
		s.MaxAltitude = t.MaxAltitude
	}
	if t.MaxClimbRate > s.MaxClimbRate {
		s.MaxClimbRate = t.MaxClimbRate
	}
	s.AverageClimbRate += t.AverageClimbRate
}

func (s *FlightStatistic) Finalize() {
	endTime := s.Points[len(s.Points)-1].Time
	startTime := s.Points[0].Time
	s.TotalFlightTime = endTime.Sub(startTime)

	if s.NumberOfThermals > 0 {
		s.AverageClimbRate /= float64(s.NumberOfThermals)
	}
	s.PercentageThermic = float64(s.TotalThermicTime) / float64(s.TotalFlightTime) * 100
}

func (s FlightStatistic) PrettyPrint() {
	fmt.Println("Thermal Statistics:")
	fmt.Printf("Total Thermic Time: %v\n", s.TotalThermicTime)
	fmt.Printf("Total Flight Time: %v\n", s.TotalFlightTime)
	fmt.Printf("Total Climb: %v\n", s.TotalClimb)
	fmt.Printf("Max Climb in a Single Thermal: %dm\n", s.MaxClimb)
	fmt.Printf("Max Altitude in Thermal: %dm\n", s.MaxAltitude)
	fmt.Printf("Average Climb Rate in Thermals: %.2f m/s\n", s.AverageClimbRate)
	fmt.Printf("Max Climb Rate in Thermals: %.2f m/s\n", s.MaxClimbRate)
	fmt.Printf("Number of Thermals Encountered: %d\n", s.NumberOfThermals)
	fmt.Printf("Percentage of Flight in Thermals: %.2f%%\n", s.PercentageThermic)
}

func Average(numbers []float64) float64 {
	total := 0.0
	for _, number := range numbers {
		total += number
	}
	if len(numbers) == 0 {
		return 0
	}
	return total / float64(len(numbers))
}

func UpdateRateOfClimbHistory(history []float64, rateOfClimb float64, period int) []float64 {
	history = append(history, rateOfClimb)
	if len(history) > period {
		history = history[1:]
	}
	return history
}
