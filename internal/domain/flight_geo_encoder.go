package domain

import (
	"encoding/json"

	"github.com/AurelienS/cigare/internal/util"
	"github.com/ezgliding/goigc/pkg/igc"
)

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

func (f Flight) GenerateGeoJSON() (string, error) {
	track, err := igc.Parse(f.IgcData)
	if err != nil {
		util.Error().Err(err).Msg("Error parsing IGC data")
	}

	var coordinates [][]float64
	for _, ll := range track.Points {
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
