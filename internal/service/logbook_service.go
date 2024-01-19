package service

import (
	"archive/zip"
	"context"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/internal/repository"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/ezgliding/goigc/pkg/igc"
)

type LogbookService struct {
	logbookRepo repository.FlightRepository
}

func NewLogbookService(
	logbookRepo repository.FlightRepository,
) LogbookService {
	return LogbookService{
		logbookRepo: logbookRepo,
	}
}

func (s *LogbookService) ProcessAndAddFlight(ctx context.Context, file *multipart.FileHeader, user domain.User) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if strings.HasSuffix(file.Filename, ".zip") {
		e := s.processZipFile(ctx, src, file.Size, user)
		return e
	}
	e := s.processSingleFile(ctx, src, file.Filename, user)

	return e
}

func (s *LogbookService) processSingleFile(ctx context.Context, reader io.Reader, filename string,
	user domain.User,
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
	stats := domain.NewFlightStatistics(track.Points)
	err = s.logbookRepo.InsertFlight(ctx, flight, stats, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *LogbookService) processZipFile(
	ctx context.Context,
	zipReader io.ReaderAt,
	size int64,
	user domain.User,
) error {
	var flights []domain.Flight
	var flightStats []domain.FlightStatistic

	zr, err := zip.NewReader(zipReader, size) // 'size' should be the size of the zip file
	if err != nil {
		return err
	}

	dates := map[time.Time]struct{}{}
	addedCount := 0
	errorCount := 0
	for _, f := range zr.File {
		if strings.ToLower(filepath.Ext(f.Name)) == ".igc" {
			rc, err1 := f.Open()
			if err1 != nil {
				errorCount++
				continue
			}

			byteContent, err1 := io.ReadAll(rc)
			rc.Close() // Ensure file is closed after reading
			if err1 != nil {
				errorCount++
				continue
			}

			track, err1 := igc.Parse(string(byteContent))
			if err1 != nil {
				errorCount++
				continue
			}

			flight := TrackToFlight(track)
			if _, exists := dates[flight.Date]; exists {
				errorCount++
				continue
			}
			dates[flight.Date] = struct{}{}
			stats := domain.NewFlightStatistics(track.Points)

			flights = append(flights, flight)
			flightStats = append(flightStats, stats)

			addedCount++
		}
	}

	err = s.logbookRepo.InsertFlights(ctx, flights, flightStats, user)
	if err != nil {
		return err
	}

	util.
		Info().
		Str("user", user.Email).
		Int("added", addedCount).
		Int("errors", errorCount).
		Msg("File processed and flight record created successfully")

	return nil
}

func (s LogbookService) GetStatisticsByYearAndMonth(
	ctx context.Context,
	user domain.User,
) (domain.StatsYearMonth, error) {
	statsYearMonth := domain.StatsYearMonth{}

	flights, err := s.logbookRepo.GetFlights(ctx, time.Time{}, time.Now(), user)
	if err != nil {
		return statsYearMonth, err
	}

	// Prepare map to hold statistics by year and month
	flightsStatisticsByYearMonth := make(map[int]map[time.Month][]domain.Flight)
	for _, flight := range flights {
		year, month, _ := flight.Date.Date()

		// Initialize year and month if not already present
		if flightsStatisticsByYearMonth[year] == nil {
			flightsStatisticsByYearMonth[year] = make(map[time.Month][]domain.Flight)
			for m := time.January; m <= time.December; m++ {
				flightsStatisticsByYearMonth[year][m] = []domain.Flight{}
			}
		}
		flightsStatisticsByYearMonth[year][month] = append(flightsStatisticsByYearMonth[year][month], flight)
	}

	// Flatten the YearMonth stats to aggregated stats
	for year, monthStats := range flightsStatisticsByYearMonth {
		if statsYearMonth[year] == nil {
			statsYearMonth[year] = make(map[time.Month]domain.StatsAggregated)
		}
		for month, stats := range monthStats {
			statsYearMonth[year][month] = domain.ComputeAggregateStatistics(stats)
		}
	}

	return statsYearMonth, err
}

func (s LogbookService) GetStatistics(
	ctx context.Context,
	startDate, endDate time.Time,
	user domain.User,
) (domain.StatsAggregated, error) {
	logStats := domain.StatsAggregated{}

	flights, err := s.logbookRepo.GetFlights(ctx, startDate, endDate, user)
	if err != nil {
		return logStats, err
	}

	aggregatedStats := domain.ComputeAggregateStatistics(flights)

	return aggregatedStats, nil
}

func (s LogbookService) GetFlights(ctx context.Context, year int, user domain.User) ([]domain.Flight, error) {
	startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	endOfYear := time.Date(year, time.December, 31, 23, 59, 59, 999999999, time.UTC)

	return s.logbookRepo.GetFlights(ctx, startOfYear, endOfYear, user)
}

func (s LogbookService) GetFlight(ctx context.Context, flightID int, user domain.User) (domain.Flight, error) {
	return s.logbookRepo.GetFlight(ctx, flightID, user)
}

func (s LogbookService) GetFlyingYears(ctx context.Context, user domain.User) ([]int, error) {
	return s.logbookRepo.GetFlyingYears(ctx, user)
}

func (s LogbookService) GetLastFlight(ctx context.Context, user domain.User) (*domain.Flight, error) {
	return s.logbookRepo.GetLastFlight(ctx, user)
}

func TrackToFlight(externalTrack igc.Track) domain.Flight {
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

	siteName := strings.Split(externalTrack.Site, "_")
	site := "Inconnu"

	if len(siteName) > 0 {
		site = siteName[0]
	}

	flight := domain.Flight{
		Date:            combinedDateTime,
		TakeoffLocation: site,
	}

	return flight
}
