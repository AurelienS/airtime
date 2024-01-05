package drawer

import (
	"image/color"
	"log"

	"github.com/AurelienS/cigare/flight"
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

// Define colors for different phases
var cruisingColor = color.Black                            // Green for cruising
var circlingColor = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red for circling

// Draw2DMap draws the flight path on a map using different colors for cruising and circling phases
func Draw2DMap(flight flight.Flight, withThermal bool) {
	ctx := sm.NewContext()
	ctx.SetSize(600, 600)

	var flightPath []s2.LatLng
	for _, point := range flight.Points {
		flightPath = append(flightPath, s2.LatLng{Lat: s1.Angle(point.Lat), Lng: s1.Angle(point.Lng)})
	}

	ctx.AddObject(sm.NewPath(flightPath, cruisingColor, 2))

	if withThermal {
		for _, thermal := range flight.Thermals {
			var path []s2.LatLng
			for i := thermal.StartIndex; i <= thermal.EndIndex; i++ {
				p := flight.Points[i]
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

func DrawElevation(f flight.Flight) {
	p := plot.New()

	p.Title.Text = "Elevation Diagram"
	p.Y.Label.Text = "Elevation (m)"

	pts := make(plotter.XYs, len(f.Points))
	startTime := f.Points[0].Time // assuming the slice is not empty and is sorted by time
	for i, point := range f.Points {
		pts[i].X = float64(point.Time.Sub(startTime).Minutes()) // X-axis in minutes since start
		pts[i].Y = float64(point.GNSSAltitude)
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatalf("Could not create line: %v", err)
	}
	line.Color = color.RGBA{B: 255, A: 255}

	p.Add(line)
	p.Y.Min = 0

	// Set the custom ticker
	p.X.Tick.Marker = CustomTicker{StartTime: startTime}

	// Save the plot to a PNG file.
	if err := p.Save(600, 200, "elevationChart.png"); err != nil {
		log.Fatalf("Could not save elevationChart: %v", err)
	}
}
