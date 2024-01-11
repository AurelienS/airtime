package flight

import (
	"image/color"
	"math"
	"time"

	"git.sr.ht/~sbinet/gg"
	"github.com/AurelienS/cigare/internal/user"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/ezgliding/goigc/pkg/igc"
	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func TrackToFlight(externalTrack igc.Track, user user.User) Flight {

	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		util.Warn().Msg("Error loading location Europe/Paris for")
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

	flight := Flight{
		Date:            combinedDateTime,
		TakeoffLocation: externalTrack.Site,
	}

	flight.UserID = user.ID
	flight.GliderID = user.DefaultGliderID

	return flight
}

var cruisingColor = color.Black
var circlingColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}

func Draw2DMap(withThermal bool, track igc.Track) {
	ctx := sm.NewContext()
	ctx.SetSize(600, 600)

	var flightPath []s2.LatLng
	for _, point := range track.Points {
		flightPath = append(flightPath, s2.LatLng{Lat: s1.Angle(point.Lat), Lng: s1.Angle(point.Lng)})
	}

	ctx.AddObject(sm.NewPath(flightPath, cruisingColor, 2))

	// if withThermal {
	// 	for _, thermal := range f.Thermals {
	// 		var path []s2.LatLng
	// 		for i := thermal.StartIndex; i <= thermal.EndIndex; i++ {
	// 			p := f.Points[i]
	// 			path = append(path, s2.LatLng{Lat: s1.Angle(p.Lat), Lng: s1.Angle(p.Lng)})
	// 		}

	// 		ctx.AddObject(sm.NewPath(path, circlingColor, 2))
	// 	}
	// }

	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	if err := gg.SavePNG("2DMap.png", img); err != nil {
		panic(err)
	}
}

func DrawElevation(track igc.Track) {
	p := plot.New()

	p.Title.Text = "Elevation Diagram"
	p.Y.Label.Text = "Elevation (m)"

	pts := make(plotter.XYs, len(track.Points))
	startTime := track.Points[0].Time
	for i, point := range track.Points {
		pts[i].X = float64(point.Time.Sub(startTime).Minutes())
		pts[i].Y = float64(point.GNSSAltitude)
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		util.Error().Msgf("Could not create line: %v", err)
	}
	line.Color = color.RGBA{B: 255, A: 255}

	p.Add(line)
	p.Y.Min = 0

	p.X.Tick.Marker = HourTicker{StartTime: startTime}

	if err := p.Save(600, 200, "elevationChart.png"); err != nil {
		util.Error().Msgf("Could not save elevationChart: %v", err)
	}
}

type HourTicker struct {
	StartTime time.Time
}

func (t HourTicker) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	step := math.Round((max - min) / 10)

	for x := min; x <= max; x += step {
		tickTime := t.StartTime.Add(time.Duration(x) * time.Minute)

		label := tickTime.Format("15:04")

		ticks = append(ticks, plot.Tick{Value: x, Label: label})
	}
	return ticks
}
