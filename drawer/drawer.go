// drawer package
package drawer

import (
	"image/color"

	"github.com/AurelienS/cigare/model"
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// Drawer draws a path on a map using flopp/go-staticmaps and saves it as an image
func Drawer(flight model.Flight) {
	ctx := sm.NewContext()

	var path []s2.LatLng

	for _, p := range flight.Points {
		path = append(path, s2.LatLng{Lat: s1.Angle(p.Lat), Lng: s1.Angle(p.Lng)})
	}
	ctx.SetSize(400, 300)
	ctx.AddObject(sm.NewPath(path, color.Black, 2))

	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	if err := gg.SavePNG("my-map.png", img); err != nil {
		panic(err)
	}
}

// lonLatToXY converts latitude and longitude to screen coordinates
func lonLatToXY(lon, lat float64, width, height float64) (float64, float64) {
	x := (lon + 180) * (width / 360)
	y := (90 - lat) * (height / 180)
	return x, y
}
