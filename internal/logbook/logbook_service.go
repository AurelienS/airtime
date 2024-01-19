package logbook

import (
	"archive/zip"
	"context"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

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
	stats := model.NewFlightStatistics(track.Points)
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

func (s Service) GetStatisticsByYearAndMonth(ctx context.Context, user model.User) (model.StatsYearMonth, error) {
	statsYearMonth := model.StatsYearMonth{}

	flights, err := s.logbookRepo.GetFlights(ctx, time.Time{}, time.Now(), user)
	if err != nil {
		return statsYearMonth, err
	}

	// Prepare map to hold statistics by year and month
	flightsStatisticsByYearMonth := make(map[int]map[time.Month][]model.FlightStatistic)
	for _, flight := range flights {
		stat := flight.Statistic
		year, month, _ := flight.Date.Date()

		// Initialize year and month if not already present
		if flightsStatisticsByYearMonth[year] == nil {
			flightsStatisticsByYearMonth[year] = make(map[time.Month][]model.FlightStatistic)
			for m := time.January; m <= time.December; m++ {
				flightsStatisticsByYearMonth[year][m] = []model.FlightStatistic{}
			}
		}
		flightsStatisticsByYearMonth[year][month] = append(flightsStatisticsByYearMonth[year][month], stat)
	}

	// Flatten the YearMonth stats to aggregated stats
	for year, monthStats := range flightsStatisticsByYearMonth {
		if statsYearMonth[year] == nil {
			statsYearMonth[year] = make(map[time.Month]model.StatsAggregated)
		}
		for month, stats := range monthStats {
			statsYearMonth[year][month] = model.ComputeAggregateStatistics(stats)
		}
	}

	return statsYearMonth, err
}

func (s Service) GetStatistics(
	ctx context.Context,
	startDate, endDate time.Time,
	user model.User,
) (model.StatsAggregated, error) {
	logStats := model.StatsAggregated{}

	stats, err := s.logbookRepo.GetStatistics(ctx, startDate, endDate, user)
	if err != nil {
		return logStats, err
	}

	aggregatedStats := model.ComputeAggregateStatistics(stats)

	return aggregatedStats, nil
}

func (s Service) GetFlights(ctx context.Context, year int, user model.User) ([]model.Flight, error) {
	startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	endOfYear := time.Date(year, time.December, 31, 23, 59, 59, 999999999, time.UTC)

	return s.logbookRepo.GetFlights(ctx, startOfYear, endOfYear, user)
}

func (s Service) GetFlyingYears(ctx context.Context, user model.User) ([]int, error) {
	return s.logbookRepo.GetFlyingYears(ctx, user)
}

func (s Service) GetLastFlight(ctx context.Context, user model.User) (*model.Flight, error) {
	return s.logbookRepo.GetLastFlight(ctx, user)
}
