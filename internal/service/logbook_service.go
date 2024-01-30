package service

import (
	"archive/zip"
	"context"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"

	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/internal/repository"
	"github.com/AurelienS/cigare/internal/util"
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

func (s *LogbookService) AddIGCFlight(ctx context.Context, file *multipart.FileHeader, user domain.User) error {
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
	flightChan, errChan := s.setupChannels(len(zr.File))

	for _, f := range zr.File {
		if isIgcFile(f.Name) {
			wg.Add(1)
			go s.processIgcZipFile(f, flightChan, errChan, &wg)
		}
	}

	wg.Wait()
	close(flightChan)
	close(errChan)

	return s.collectAndInsertFlights(ctx, flightChan, errChan, user)
}

func (s *LogbookService) processIgcZipFile(
	file *zip.File,
	flightChan chan<- domain.Flight,
	errChan chan<- error,
	wg *sync.WaitGroup,
) {
	util.Debug().Str("filename", file.Name).Msg("Processing IGC file")
	defer func() {
		util.Debug().Str("filename", file.Name).Msg("Done processing IGC file")
		wg.Done()
	}()
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
	flight, err := domain.NewFlightFromIgc(content)
	if err != nil {
		errChan <- err
		return
	}
	flightChan <- flight
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
	flight, err := domain.NewFlightFromIgc(content)
	if err != nil {
		return err
	}

	err = s.logbookRepo.InsertFlight(ctx, flight, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *LogbookService) setupChannels(
	fileCount int,
) (chan domain.Flight, chan error) {
	util.Debug().Int("fileCount", fileCount).Msg("Setting up channels")
	flightChan := make(chan domain.Flight, fileCount)
	errChan := make(chan error, fileCount)

	return flightChan, errChan
}

func (s *LogbookService) collectAndInsertFlights(
	ctx context.Context,
	flightChan <-chan domain.Flight,
	errChan <-chan error,
	user domain.User,
) error {
	var flights []domain.Flight
	var errors []error

	for flightChan != nil || errChan != nil {
		select {
		case flight, ok := <-flightChan:
			if !ok {
				flightChan = nil
				continue
			}
			flights = append(flights, flight)
		case pErr, ok := <-errChan:
			if !ok {
				errChan = nil
				continue
			}
			errors = append(errors, pErr)
		}
	}

	if len(errors) > 0 {
		util.Warn().
			Str("user", user.Email).
			Int("count", len(errors)).
			Errs("errors", errors).
			Msg("Errors during processing files")
	}

	err := s.logbookRepo.InsertFlights(ctx, flights, user)
	util.Debug().Str("user", user.Email).Msg("Finished processing files")
	return err
}

func isZipFile(filename string) bool {
	return strings.HasSuffix(filename, ".zip")
}

func isIgcFile(filename string) bool {
	return strings.ToLower(filepath.Ext(filename)) == ".igc"
}
