package drawer

import (
	"image/color"

	"github.com/AurelienS/cigare/flight"
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// Define colors for different phases
var cruisingColor = color.RGBA{R: 0, G: 255, B: 0, A: 255} // Green for cruising
var circlingColor = color.RGBA{R: 0, G: 0, B: 0, A: 255}   // Red for circling

// Drawer draws the flight path on a map using different colors for cruising and circling phases
func Drawer(flight flight.Flight) {
	ctx := sm.NewContext()
	ctx.SetSize(600, 600)

	var flightPath []s2.LatLng
	for _, point := range flight.Points {
		flightPath = append(flightPath, s2.LatLng{Lat: s1.Angle(point.Lat), Lng: s1.Angle(point.Lng)})
	}

	ctx.AddObject(sm.NewPath(flightPath, cruisingColor, 2))

	for _, thermal := range flight.Thermals {
		var path []s2.LatLng
		for i := thermal.StartIndex; i <= thermal.EndIndex; i++ {
			p := flight.Points[i]
			path = append(path, s2.LatLng{Lat: s1.Angle(p.Lat), Lng: s1.Angle(p.Lng)})
		}

		ctx.AddObject(sm.NewPath(path, circlingColor, 2))
	}

	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	if err := gg.SavePNG("flight-phases-map.png", img); err != nil {
		panic(err)
	}
}
