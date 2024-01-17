package logbook

import (
	"archive/zip"
	"context"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	flightstats "github.com/AurelienS/cigare/internal/flight_statistic"
	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/ezgliding/goigc/pkg/igc"
)

type Service struct {
	logbookRepo Repository
}

func NewService(
	logbookRepo Repository,
) Service {
	return Service{
		logbookRepo: logbookRepo,
	}
}

func (s *Service) ProcessAndAddFlight(ctx context.Context, file *multipart.FileHeader, user model.User) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if strings.HasSuffix(file.Filename, ".zip") {
		return s.processZipFile(ctx, src, file.Size, user)
	}

	return s.processSingleFile(ctx, src, file.Filename, user)
}

func (s *Service) processSingleFile(ctx context.Context, reader io.Reader, filename string,
	user model.User,
) error {
	byteContent, err := io.ReadAll(reader)
	if err != nil {
		util.Error().Err(err).Str("filename", filename).Msg("Failed to read IGC file")
		return err
	}

	track, err := igc.Parse(string(byteContent))
	if err != nil {
		return err
	}

	flight := TrackToFlight(track)
	stats := flightstats.NewFlightStatistics(track.Points)
	err = s.logbookRepo.InsertFlight(ctx, flight, stats, user)
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to insert flight into database")
		return err
	}

	util.Info().Str("user", user.Email).Str("filename", filename).
		Msg("File processed and flight record created successfully")

	return nil
}

func (s *Service) processZipFile(ctx context.Context, zipReader io.ReaderAt, size int64, user model.User) error {
	zr, err := zip.NewReader(zipReader, size) // 'size' should be the size of the zip file
	if err != nil {
		return err
	}

	addedCount := 0
	errorCount := 0

	for _, f := range zr.File {
		if strings.ToLower(filepath.Ext(f.Name)) == ".igc" {
			file, err := f.Open()
			if err != nil {
				continue // or handle the error
			}

			err = s.processSingleFile(ctx, file, f.Name, user)
			file.Close()

			if err != nil {
				errorCount++
			} else {
				addedCount++
			}
		}
	}

	util.
		Info().
		Str("user", user.Email).
		Int("added", addedCount).
		Int("errors", errorCount).
		Msg("File processed and flight record created successfully")

	return nil
}

type Stats struct {
	FlightCount           int
	MaxAltitude           int
	MaxClimb              int
	TotalClimb            int
	TotalNumberOfThermals int
	MaxClimbRate          float64
	MaxFlightLength       time.Duration
	MinFlightLength       time.Duration
	AverageFlightLength   time.Duration
	TotalFlightTime       time.Duration
	TotalThermicTime      time.Duration
}

func (s Service) GetStatistics(ctx context.Context, startDate, endDate time.Time, user model.User) (Stats, error) {
	logStats := Stats{}

	stats, err := s.logbookRepo.GetStatistics(ctx, startDate, endDate, user)
	if err != nil {
		return logStats, err
	}

	var maxAltitude int
	var maxVario float64
	var maxFlightLength time.Duration
	minFlightLength := time.Duration(0)
	averageFlightLength := time.Duration(0)
	var totalFlightTime time.Duration

	var totalThermicTime time.Duration
	var maxClimb int
	var totalClimb int
	var totalNumberOfThermals int

	flightCount := len(stats)

	for _, stat := range stats {
		if stat.MaxAltitude > maxAltitude {
			maxAltitude = stat.MaxAltitude
		}
		if stat.MaxClimbRate > maxVario {
			maxVario = stat.MaxClimbRate
		}
		if stat.TotalFlightTime > maxFlightLength {
			maxFlightLength = stat.TotalFlightTime
		}
		if stat.TotalFlightTime < minFlightLength {
			minFlightLength = stat.TotalFlightTime
		}
		if stat.MaxClimb > maxClimb {
			maxClimb = stat.MaxClimb
		}
		totalClimb += stat.TotalClimb
		totalNumberOfThermals += stat.NumberOfThermals
		totalThermicTime += stat.TotalThermicTime
		totalFlightTime += stat.TotalFlightTime
	}

	if flightCount > 0 {
		averageFlightLength = totalFlightTime / time.Duration(flightCount)
	}

	logStats = Stats{
		MaxAltitude:           maxAltitude,
		MaxClimbRate:          maxVario,
		MaxFlightLength:       maxFlightLength,
		MinFlightLength:       minFlightLength,
		AverageFlightLength:   averageFlightLength,
		TotalFlightTime:       totalFlightTime,
		FlightCount:           flightCount,
		MaxClimb:              maxClimb,
		TotalClimb:            totalClimb,
		TotalNumberOfThermals: totalNumberOfThermals,
		TotalThermicTime:      totalThermicTime,
	}

	return logStats, nil
}

func (s Service) GetFlights(ctx context.Context, year int, user model.User) ([]model.Flight, error) {
	return s.logbookRepo.GetFlights(ctx, year, user)
}

func (s Service) GetFlyingYears(ctx context.Context, user model.User) ([]int, error) {
	return s.logbookRepo.GetFlyingYears(ctx, user)
}

func (s Service) GetLastFlight(ctx context.Context, user model.User) (*model.Flight, error) {
	return s.logbookRepo.GetLastFlight(ctx, user)
}
