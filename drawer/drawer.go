package drawer

import (
	"image/color"

	"github.com/AurelienS/cigare/flight"
)

// Define colors for different phases
var cruisingColor = color.RGBA{R: 0, G: 255, B: 0, A: 255} // Green for cruising
var circlingColor = color.RGBA{R: 0, G: 0, B: 0, A: 255}   // Red for circling

// Drawer draws the flight path on a map using different colors for cruising and circling phases
func Drawer(flight flight.Flight) {
	// ctx := sm.NewContext()
	// ctx.SetSize(600, 600)

	// for _, phase := range flight.Phases {
	// 	var path []s2.LatLng
	// 	for i := phase.StartIndex; i <= phase.EndIndex; i++ {
	// 		p := flight.Points[i]
	// 		path = append(path, s2.LatLng{Lat: s1.Angle(p.Lat), Lng: s1.Angle(p.Lng)})
	// 	}

	// 	// Choose color based on phase type
	// 	var phaseColor color.Color
	// 	if phase.Type == model.Cruising {
	// 		phaseColor = cruisingColor
	// 	} else if phase.Type == model.Circling {
	// 		phaseColor = circlingColor
	// 	}

	// 	ctx.AddObject(sm.NewPath(path, phaseColor, 2))
	// }

	// img, err := ctx.Render()
	// if err != nil {
	// 	panic(err)
	// }

	// if err := gg.SavePNG("flight-phases-map.png", img); err != nil {
	// 	panic(err)
	// }
}
