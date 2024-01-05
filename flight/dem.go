package flight

// import (
// 	"fmt"
// 	"log"

// 	"github.com/lukeroth/gdal"
// )

// type DEMReader struct {
// 	DemPath string
// }

// func NewDEMReader(demPath string) *DEMReader {
// 	return &DEMReader{DemPath: demPath}
// }

// func (d *DEMReader) GetElevations(flight Flight) ([]float64, error) {
// 	dataset, err := gdal.Open(d.DemPath, gdal.Access(1))
// 	if err != nil {
// 		return nil, fmt.Errorf("error opening DEM: %v", err)
// 	}
// 	defer dataset.Close()

// 	geoTransform := dataset.GeoTransform()
// 	inverseTransform := gdal.InvGeoTransform(geoTransform)

// 	elevations := make([]float64, len(flight.Points))
// 	for i, point := range flight.Points {
// 		// Manually apply the inverse geotransform to get the pixel coordinates
// 		px := inverseTransform[0] + inverseTransform[1]*point.Lng + inverseTransform[2]*point.Lat
// 		py := inverseTransform[3] + inverseTransform[4]*point.Lng + inverseTransform[5]*point.Lat
// 		x, y := int(px+0.5), int(py+0.5)

// 		raster := dataset.RasterBand(1)
// 		elevation := make([]int16, 1)
// 		if err := raster.IO(gdal.RWFlag(1), x, y, 1, 1, elevation, 1, 1, 0, 0); err != nil {
// 			log.Printf("Error reading elevation for point %d: %v", i, err)
// 			continue
// 		}
// 		elevations[i] = float64(elevation[0])
// 	}

// 	return elevations, nil
// }
