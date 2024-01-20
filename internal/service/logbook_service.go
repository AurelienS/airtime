package service

import (
	"archive/zip"
	"context"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"
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

	if isZipFile(file.Filename) {
		return s.processZipFile(ctx, src, file.Size, user)
	}
	return s.processSingleFile(ctx, src, file.Filename, user)
}

func (s *LogbookService) processZipFile(
	ctx context.Context,
	zipReader io.ReaderAt,
	size int64,
	user domain.User,
) error {
	zr, err := zip.NewReader(zipReader, size)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	flightChan, statsChan, errChan := s.setupChannels(len(zr.File), &wg)

	for _, f := range zr.File {
		if isIgcFile(f.Name) {
			wg.Add(1)
			go s.processIgcZipFile(f, flightChan, statsChan, errChan, &wg)
		}
	}

	return s.collectAndInsertFlights(ctx, flightChan, statsChan, errChan, user)
}

func (s *LogbookService) processIgcZipFile(
	file *zip.File,
	flightChan chan<- domain.Flight,
	statsChan chan<- domain.FlightStatistic,
	errChan chan<- error,
	wg *sync.WaitGroup,
) {
	rc, err := file.Open()
	if err != nil {
		errChan <- err
		return
	}

	byteContent, err := io.ReadAll(rc)
	if err != nil {
		errChan <- err
		return
	}

	err = rc.Close()
	if err != nil {
		errChan <- err
		return
	}

	content := string(byteContent)
	flight, stats, err := s.processIgcFile(content)
	if err != nil {
		errChan <- err
		return
	}
	flightChan <- flight
	statsChan <- stats

	wg.Done()
}

func (s *LogbookService) processSingleFile(
	ctx context.Context,
	reader io.Reader,
	filename string,
	user domain.User,
) error {
	byteContent, err := io.ReadAll(reader)
	if err != nil {
		util.Error().Err(err).Str("filename", filename).Msg("Failed to read IGC file")
		return err
	}
	content := string(byteContent)
	flight, stats, err := s.processIgcFile(content)
	if err != nil {
		return err
	}

	err = s.logbookRepo.InsertFlight(ctx, flight, stats, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *LogbookService) processIgcFile(content string) (domain.Flight, domain.FlightStatistic, error) {
	track, err := igc.Parse(content)
	if err != nil {
		return domain.Flight{}, domain.FlightStatistic{}, err
	}

	flight := TrackToFlight(track)
	stats := domain.NewFlightStatistics(track.Points)

	return flight, stats, nil
}

func (s *LogbookService) setupChannels(
	fileCount int,
	wg *sync.WaitGroup,
) (chan domain.Flight, chan domain.FlightStatistic, chan error) {
	util.Debug().Int("fileCount", fileCount).Msg("Setting up channels")
	flightChan := make(chan domain.Flight, fileCount)
	statsChan := make(chan domain.FlightStatistic, fileCount)
	errChan := make(chan error, fileCount)

	go func() {
		wg.Wait()
		close(flightChan)
		close(statsChan)
		close(errChan)
	}()

	return flightChan, statsChan, errChan
}

func (s *LogbookService) collectAndInsertFlights(
	ctx context.Context,
	flightChan <-chan domain.Flight,
	statsChan <-chan domain.FlightStatistic,
	errChan <-chan error,
	user domain.User,
) error {
	var flights []domain.Flight
	var flightStats []domain.FlightStatistic
	var errors []error

	for flightChan != nil || statsChan != nil || errChan != nil {
		select {
		case flight, ok := <-flightChan:
			if !ok {
				flightChan = nil
				continue
			}
			flights = append(flights, flight)
		case stat, ok := <-statsChan:
			if !ok {
				statsChan = nil
				continue
			}
			flightStats = append(flightStats, stat)
		case pErr, ok := <-errChan:
			if !ok {
				errChan = nil
				continue
			}
			errors = append(errors, pErr)
		}
	}

	if len(errors) > 0 {
		util.Warn().Str("user", user.Email).Errs("errors", errors).Msg("Errors during processing files")
	}

	return s.logbookRepo.InsertFlights(ctx, flights, flightStats, user)
}

func isZipFile(filename string) bool {
	return strings.HasSuffix(filename, ".zip")
}

func isIgcFile(filename string) bool {
	return strings.ToLower(filepath.Ext(filename)) == ".igc"
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

func (s LogbookService) GetFlights(
	ctx context.Context,
	start time.Time,
	end time.Time,
	user domain.User,
) ([]domain.Flight, error) {
	return s.logbookRepo.GetFlights(ctx, start, end, user)
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
