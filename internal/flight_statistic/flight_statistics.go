package flightstats

import (
	"fmt"
	"math"
	"strings"
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

func NewFlightStatistics(points []igc.Point) FlightStatistic {
	var pts []Point
	for _, p := range points {
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

func (fs *FlightStatistic) Compute() {
	const (
		minClimbRate               = 0.2              // m/s, the threshold for considering it thermic activity
		climbRateIntegrationPeriod = 10               // Number of seconds to smooth the climbRate
		minThermalDuration         = 20 * time.Second // Minimum duration to consider a sustained climb as thermal
		allowedDownwardPoints      = 4                // Number of consecutive downward Points allowed in a thermal
	)

	var current *Thermal
	var rateOfClimbHistory []float64

	for i, point := range fs.Points {
		if i == 0 {
			continue
		}

		// flight related
		fs.MaxAltitude = int(math.Max(float64(fs.MaxAltitude), float64(point.GNSSAltitude)))
		fmt.Println("file: flight_statistics.go ~ line 67 ~ func ~ point.GNSSAltitude : ", point.GNSSAltitude)
		fmt.Println("file: flight_statistics.go ~ line 67 ~ func ~ fs.MaxAltitude : ", fs.MaxAltitude)

		// thermal related
		smoothedRateOfClimb := fs.calculateSmoothedRateOfClimb(i, climbRateIntegrationPeriod, rateOfClimbHistory)
		if current == nil {
			current = fs.maybeStartNewThermal(smoothedRateOfClimb, minClimbRate, climbRateIntegrationPeriod, point, i)
		} else {
			current.Update(point, smoothedRateOfClimb)
			current = fs.checkAndFinalizeThermal(current, allowedDownwardPoints, minThermalDuration, i)
		}
	}

	fs.finalizeLastThermal(current, minThermalDuration)
	fs.Finalize()
}

func (fs *FlightStatistic) calculateSmoothedRateOfClimb(i, period int, rateOfClimbHistory []float64) float64 {
	rateOfClimb := fs.calculateRateOfClimb(i)
	rateOfClimbHistory = UpdateRateOfClimbHistory(rateOfClimbHistory, rateOfClimb, period)
	return Average(rateOfClimbHistory)
}

func (fs *FlightStatistic) calculateRateOfClimb(i int) float64 {
	altitudeGain := fs.Points[i].GNSSAltitude - fs.Points[i-1].GNSSAltitude
	timeElapsed := fs.Points[i].Time.Sub(fs.Points[i-1].Time).Seconds()
	return float64(altitudeGain) / timeElapsed
}

func (fs *FlightStatistic) checkAndFinalizeThermal(
	current *Thermal,
	tolerance int,
	duration time.Duration,
	index int,
) *Thermal {
	if current.ShouldEnd(tolerance) {
		if current.Duration() >= duration {
			current.EndIndex = index
			current.AverageClimbRate = float64(current.Climb()) / current.Duration().Seconds()
			fs.Thermals = append(fs.Thermals, current)
			fs.AddThermal(*current, duration)
		}
		return nil
	}
	return current
}

func (fs *FlightStatistic) maybeStartNewThermal(
	smoothedRate float64,
	minRate float64,
	period int,
	point Point,
	index int,
) *Thermal {
	if smoothedRate >= minRate && len(fs.Points) >= period {
		return NewThermal(point.Time, point.GNSSAltitude, index)
	}
	return nil
}

func (fs *FlightStatistic) finalizeLastThermal(current *Thermal, duration time.Duration) {
	if current != nil && current.Duration() >= duration {
		current.EndIndex = len(fs.Points) - 1
		fs.Thermals = append(fs.Thermals, current)
		fs.AddThermal(*current, duration)
	}
}

func (fs *FlightStatistic) AddThermal(t Thermal, minThermalDuration time.Duration) {
	duration := t.Duration()
	if duration < minThermalDuration || t.Climb() <= 0 {
		return
	}

	fs.NumberOfThermals++
	fs.TotalThermicTime += duration

	climb := t.Climb()
	fs.TotalClimb += climb
	fs.MaxClimb = int(math.Max(float64(fs.MaxClimb), float64(climb)))
	fs.MaxClimbRate = math.Max(fs.MaxClimbRate, t.MaxClimbRate)
	fs.AverageClimbRate += t.AverageClimbRate
}

func (fs *FlightStatistic) Finalize() {
	endTime := fs.Points[len(fs.Points)-1].Time
	startTime := fs.Points[0].Time
	fs.TotalFlightTime = endTime.Sub(startTime)

	if fs.NumberOfThermals > 0 {
		fs.AverageClimbRate /= float64(fs.NumberOfThermals)
	}
	fs.PercentageThermic = float64(fs.TotalThermicTime) / float64(fs.TotalFlightTime) * 100
}

func (fs FlightStatistic) PrettyPrint() string {
	var sb strings.Builder

	sb.WriteString("Thermal Statistics:\n")
	sb.WriteString(fmt.Sprintf("Total Thermic Time: %v\n", fs.TotalThermicTime))
	sb.WriteString(fmt.Sprintf("Total Flight Time: %v\n", fs.TotalFlightTime))
	sb.WriteString(fmt.Sprintf("Total Climb: %v\n", fs.TotalClimb))
	sb.WriteString(fmt.Sprintf("Max Climb in a Single Thermal: %dm\n", fs.MaxClimb))
	sb.WriteString(fmt.Sprintf("Max Altitude in Thermal: %dm\n", fs.MaxAltitude))
	sb.WriteString(fmt.Sprintf("Average Climb Rate in Thermals: %.2f m/s\n", fs.AverageClimbRate))
	sb.WriteString(fmt.Sprintf("Max Climb Rate in Thermals: %.2f m/s\n", fs.MaxClimbRate))
	sb.WriteString(fmt.Sprintf("Number of Thermals Encountered: %d\n", fs.NumberOfThermals))
	sb.WriteString(fmt.Sprintf("Percentage of Flight in Thermals: %.2f%%\n", fs.PercentageThermic))

	return sb.String()
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
