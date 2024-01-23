package domain

import (
	"encoding/json"
	"math"
	"time"

	"github.com/AurelienS/cigare/internal/util"
	"github.com/ezgliding/goigc/pkg/igc"
)

type FlightStatistic struct {
	ID                int
	TotalThermicTime  time.Duration
	TotalFlightTime   time.Duration
	MaxClimb          int
	MaxClimbRate      float64
	TotalClimb        int
	TotalDistance     int
	AverageClimbRate  float64
	NumberOfThermals  int
	PercentageThermic float64
	MaxAltitude       int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Thermals          []*Thermal
	Points            []Point
	GeoJSON           string
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

	geoJSON, err := GenerateGeoJSON(points)
	if err != nil {
		util.Warn().Err(err).Msg("Cannot generate GeoJSON")
	}
	stat := FlightStatistic{
		Points:  pts,
		GeoJSON: geoJSON,
	}
	stat.Compute()
	return stat
}

type GeoJSON struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string      `json:"type"`
	Geometry   Geometry    `json:"geometry"`
	Properties interface{} `json:"properties"`
}

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

func GenerateGeoJSON(points []igc.Point) (string, error) {
	var coordinates [][]float64
	for _, ll := range points {
		coordinates = append(coordinates, []float64{ll.Lng.Degrees(), ll.Lat.Degrees()})
	}

	geoJSON := GeoJSON{
		Type: "FeatureCollection",
		Features: []Feature{
			{
				Type: "Feature",
				Geometry: Geometry{
					Type:        "LineString",
					Coordinates: coordinates,
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(geoJSON)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (fs *FlightStatistic) Compute() {
	const (
		minClimbRate               = 0.2              // m/s, the threshold for considering it thermic activity
		climbRateIntegrationPeriod = 10               // Number of seconds to smooth the climbRate
		minThermalDuration         = 20 * time.Second // Minimum duration to consider a sustained climb as thermal
		allowedDownwardPoints      = 4                // Number of consecutive downward Points allowed in a thermal
		earthRadius                = 6371e3           // Earth's radius in meters
	)

	var current *Thermal
	var rateOfClimbHistory []float64
	var lastPoint *Point

	for i := range fs.Points {
		point := &fs.Points[i]
		if i == 0 {
			lastPoint = point
			continue
		}

		// flight related
		fs.MaxAltitude = int(math.Max(float64(fs.MaxAltitude), float64(point.GNSSAltitude)))

		// distance calculation
		if i > 0 {
			fs.TotalDistance += haversineDistance(
				lastPoint.Lat.Degrees(), lastPoint.Lng.Degrees(),
				point.Lat.Degrees(), point.Lng.Degrees(), earthRadius,
			)
		}
		lastPoint = point

		// thermal related
		smoothedRateOfClimb := fs.calculateSmoothedRateOfClimb(i, climbRateIntegrationPeriod, rateOfClimbHistory)
		if current == nil {
			current = fs.maybeStartNewThermal(smoothedRateOfClimb, minClimbRate, climbRateIntegrationPeriod, *point, i)
		} else {
			current.Update(*point, smoothedRateOfClimb)
			current = fs.checkAndFinalizeThermal(current, allowedDownwardPoints, minThermalDuration, i)
		}
	}

	fs.finalizeLastThermal(current, minThermalDuration)
	fs.Finalize()
}

func haversineDistance(lat1, lon1, lat2, lon2, radius float64) int {
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	lat1 = lat1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return int(radius * c)
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

type StatsAggregated struct {
	FlightCount           int
	MaxAltitudeFlight     Flight
	MaxClimbFlight        Flight
	TotalClimb            int
	TotalDistance         int
	MaxDistanceFlight     Flight
	TotalNumberOfThermals int
	MaxClimbRateFlight    Flight
	MaxDurationFLight     Flight
	MinFlightLength       time.Duration
	AverageFlightLength   time.Duration
	TotalFlightTime       time.Duration
	TotalThermicTime      time.Duration
	FirstFlight           Flight
	LastFlight            Flight
}

func ComputeAggregateStatistics(flights []Flight) StatsAggregated {
	var maxAltitude Flight
	var maxVario Flight
	var firstFlight Flight
	var lastFlight Flight
	var maxFlightLength Flight
	minFlightLength := time.Duration(0)
	averageFlightLength := time.Duration(0)
	var totalFlightTime time.Duration
	flightCount := len(flights)
	if flightCount > 0 {
		firstFlight = flights[flightCount-1]
	}
	if flightCount > 1 {
		lastFlight = flights[0]
	}

	var totalThermicTime time.Duration
	var maxClimb Flight
	var totalClimb int
	var totalNumberOfThermals int
	var totalDistance int
	var maxDistance Flight

	for _, f := range flights {
		if f.Statistic.MaxAltitude > maxAltitude.Statistic.MaxAltitude {
			maxAltitude = f
		}
		if f.Statistic.MaxClimbRate > maxVario.Statistic.MaxClimbRate {
			maxVario = f
		}
		if f.Statistic.TotalFlightTime > maxFlightLength.Statistic.TotalFlightTime {
			maxFlightLength = f
		}
		if f.Statistic.TotalFlightTime < minFlightLength {
			minFlightLength = f.Statistic.TotalFlightTime
		}
		if f.Statistic.MaxClimb > maxClimb.Statistic.MaxClimb {
			maxClimb = f
		}
		if f.Statistic.TotalDistance > maxDistance.Statistic.TotalDistance {
			maxDistance = f
		}
		totalClimb += f.Statistic.TotalClimb
		totalNumberOfThermals += f.Statistic.NumberOfThermals
		totalThermicTime += f.Statistic.TotalThermicTime
		totalFlightTime += f.Statistic.TotalFlightTime
		totalDistance += int(f.Statistic.TotalDistance)
	}

	if flightCount > 0 {
		averageFlightLength = totalFlightTime / time.Duration(flightCount)
	}

	aggregatedStats := StatsAggregated{
		FlightCount:           flightCount,
		MaxAltitudeFlight:     maxAltitude,
		MaxClimbFlight:        maxClimb,
		TotalClimb:            totalClimb,
		TotalNumberOfThermals: totalNumberOfThermals,
		MaxClimbRateFlight:    maxVario,
		MaxDurationFLight:     maxFlightLength,
		MinFlightLength:       minFlightLength,
		AverageFlightLength:   averageFlightLength,
		TotalFlightTime:       totalFlightTime,
		TotalThermicTime:      totalThermicTime,
		TotalDistance:         totalDistance,
		MaxDistanceFlight:     maxDistance,
		FirstFlight:           firstFlight,
		LastFlight:            lastFlight,
	}
	return aggregatedStats
}
