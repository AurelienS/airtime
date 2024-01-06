package model

import (
	"fmt"
	"time"
)

type FlightStatistics struct {
	TotalThermicTime  time.Duration
	TotalFlightTime   time.Duration
	MaxClimb          int
	MaxClimbRate      float64
	TotalClimb        int
	AverageClimbRate  float64
	NumberOfThermals  int
	PercentageThermic float64
	MaxAltitude       int
}

func (s *FlightStatistics) AddThermal(t Thermal, minThermalDuration time.Duration) {
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

func (s *FlightStatistics) Finalize(f *Flight) {
	endTime := f.Points[len(f.Points)-1].Time
	startTime := f.Points[0].Time
	s.TotalFlightTime = endTime.Sub(startTime)

	if s.NumberOfThermals > 0 {
		s.AverageClimbRate /= float64(s.NumberOfThermals)
	}
	s.PercentageThermic = float64(s.TotalThermicTime) / float64(s.TotalFlightTime) * 100
}

func (s FlightStatistics) PrettyPrint() {
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
