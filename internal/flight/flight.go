package flight

import (
	"image/color"
	"time"

	"git.sr.ht/~sbinet/gg"
	"github.com/AurelienS/cigare/internal/log"
	"github.com/AurelienS/cigare/internal/storage"
	"github.com/ezgliding/goigc/pkg/igc"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"

	sm "github.com/flopp/go-staticmaps"
)

type Flight struct {
	igc.Track
	Thermals []*Thermal
	Stats    FlightStatistics
}

func ConvertToMyFlight(externalTrack igc.Track) storage.Flight {

	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Warn().Msg("Error loading location Europe/Paris for")
	}

	combinedDateTime := time.Date(
		externalTrack.Date.Year(),
		externalTrack.Date.Month(),
		externalTrack.Date.Day(),
		externalTrack.Points[0].Time.Hour(),
		externalTrack.Points[0].Time.Minute(),
		externalTrack.Points[0].Time.Second(),
		externalTrack.Points[0].Time.Nanosecond(),
		loc,
	)

	return storage.Flight{
		Date:            combinedDateTime.Local(),
		TakeoffLocation: externalTrack.Site,
	}
}

func (f *Flight) Initialize() {
	const (
		minClimbRate               = 0.2              // m/s, the threshold for considering it thermic activity
		climbRateIntegrationPeriod = 10               // Number of seconds to smooth the climbRate
		minThermalDuration         = 20 * time.Second // Minimum duration to consider a sustained climb as thermal
		allowedDownwardPoints      = 4                // Number of consecutive downward points allowed in a thermal
	)

	f.GenerateThermals(minClimbRate, allowedDownwardPoints, minThermalDuration, climbRateIntegrationPeriod)
}

func (f *Flight) GenerateThermals(minRateOfClimb float64, maxDownwardTolerance int, minThermalDuration time.Duration, climbRateIntegrationPeriod int) {
	var current *Thermal
	var rateOfClimbHistory []float64

	for i, point := range f.Points {
		if i == 0 {
			continue
		}

		smoothedRateOfClimb := f.calculateSmoothedRateOfClimb(i, climbRateIntegrationPeriod, rateOfClimbHistory)

		if current == nil {
			current = f.maybeStartNewThermal(smoothedRateOfClimb, minRateOfClimb, climbRateIntegrationPeriod, point, i)
		} else {
			current.Update(point, smoothedRateOfClimb)
			current = f.checkAndFinalizeThermal(current, maxDownwardTolerance, minThermalDuration, i)
		}
	}

	f.finalizeLastThermal(current, minThermalDuration)
	f.Stats.Finalize(f)
}

func (f *Flight) calculateSmoothedRateOfClimb(i, period int, rateOfClimbHistory []float64) float64 {
	rateOfClimb := f.calculateRateOfClimb(i)
	rateOfClimbHistory = UpdateRateOfClimbHistory(rateOfClimbHistory, rateOfClimb, period)
	return Average(rateOfClimbHistory)
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

func (f *Flight) maybeStartNewThermal(smoothedRate float64, minRate float64, period int, point igc.Point, index int) *Thermal {
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

var cruisingColor = color.Black
var circlingColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}

func (f Flight) Draw2DMap(withThermal bool) {
	ctx := sm.NewContext()
	ctx.SetSize(600, 600)

	var flightPath []s2.LatLng
	for _, point := range f.Points {
		flightPath = append(flightPath, s2.LatLng{Lat: s1.Angle(point.Lat), Lng: s1.Angle(point.Lng)})
	}

	ctx.AddObject(sm.NewPath(flightPath, cruisingColor, 2))

	if withThermal {
		for _, thermal := range f.Thermals {
			var path []s2.LatLng
			for i := thermal.StartIndex; i <= thermal.EndIndex; i++ {
				p := f.Points[i]
				path = append(path, s2.LatLng{Lat: s1.Angle(p.Lat), Lng: s1.Angle(p.Lng)})
			}

			ctx.AddObject(sm.NewPath(path, circlingColor, 2))
		}
	}

	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	if err := gg.SavePNG("2DMap.png", img); err != nil {
		panic(err)
	}
}

func (f Flight) DrawElevation() {
	p := plot.New()

	p.Title.Text = "Elevation Diagram"
	p.Y.Label.Text = "Elevation (m)"

	pts := make(plotter.XYs, len(f.Points))
	startTime := f.Points[0].Time
	for i, point := range f.Points {
		pts[i].X = float64(point.Time.Sub(startTime).Minutes())
		pts[i].Y = float64(point.GNSSAltitude)
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		log.Error().Msgf("Could not create line: %v", err)
	}
	line.Color = color.RGBA{B: 255, A: 255}

	p.Add(line)
	p.Y.Min = 0

	p.X.Tick.Marker = HourTicker{StartTime: startTime}

	if err := p.Save(600, 200, "elevationChart.png"); err != nil {
		log.Error().Msgf("Could not save elevationChart: %v", err)
	}
}
