package domain

import (
	"math"
	"strings"
	"time"

	"github.com/AurelienS/cigare/internal/util"
	"github.com/ezgliding/goigc/pkg/igc"
	"github.com/golang/geo/s2"
)

type Point struct {
	s2.LatLng
	Time             time.Time
	PressureAltitude int
	GNSSAltitude     int
	NumSatellites    int
	Description      string
}

type Flight struct {
	ID          int
	Date        time.Time
	Pilot       User
	Location    string
	Duration    time.Duration
	Distance    int
	AltitudeMax int
	IgcData     string
}

var earthRadius = 6371e3

func NewFlightFromIgc(igcData string) (Flight, error) {
	f := Flight{
		IgcData: igcData,
	}

	track, err := igc.Parse(igcData)
	if err != nil {
		return f, err
	}

	points := f.convertPoints(track.Points)

	var lastPoint *Point
	for i := range points {
		point := &points[i]
		if i == 0 {
			lastPoint = point
			continue
		}
		f.AltitudeMax = int(math.Max(float64(f.AltitudeMax), float64(point.GNSSAltitude)))

		if i > 0 {
			f.Distance += haversineDistance(
				lastPoint.Lat.Degrees(), lastPoint.Lng.Degrees(),
				point.Lat.Degrees(), point.Lng.Degrees(), earthRadius,
			)
		}
		lastPoint = point
	}

	endTime := points[len(points)-1].Time
	startTime := points[0].Time
	f.Duration = endTime.Sub(startTime)

	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		util.Warn().Msg("Error loading location Europe/Paris for")
		return f, err
	}

	correctDate := time.Date(
		track.Date.Year(),
		track.Date.Month(),
		track.Date.Day(),
		track.Points[0].Time.Hour(),
		track.Points[0].Time.Minute(),
		track.Points[0].Time.Second(),
		track.Points[0].Time.Nanosecond(),
		loc,
	)

	siteName := strings.Split(track.Site, "_")
	site := "Inconnu"

	if len(siteName) > 0 {
		if siteName[0] != "" {
			site = siteName[0]
		}
	}

	f.Location = site
	f.Date = correctDate

	return f, nil
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

func (f *Flight) convertPoints(points []igc.Point) []Point {
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

	return pts
}
