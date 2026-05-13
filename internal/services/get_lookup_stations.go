package services

import (
	"encoding/json"
	"fmt"
	"time"

	"DataLake/internal/infrastructure/s3"
	"DataLake/internal/repositories"
	"DataLake/internal/repositories/silver"
)

type GetLookupStationsService struct {
	s3Client *s3.Client
	repo     silver.Env
}

func NewGetLookupStationsService(dt string) GetLookupStationsService {
	return GetLookupStationsService{
		s3Client: s3.New(),
		repo:     silver.SetupReferencesStationIds(dt),
	}
}

func (s *GetLookupStationsService) Run() error {
	start := time.Now()
	manifest := repositories.NewManifestRepository(s.repo.Layer, s.repo.Entity, s.repo.Dt)

	err, goToExit := s.CleanUp()
	if err != nil {
		return fmt.Errorf("Fatal error during cleaning up past job")
	}

	if goToExit {
		// log something about skipped job
		return nil
	}

	manifest.MarkInProgress()
	rawRecords := 0
	breakCounter := 0

	for {
		rawData, err := s.s3Client.Get()

		if err != nil && breakCounter < 3 {
			breakCounter++
			// sometyhing to log
			continue
		} else if err != nil {
			manifest.MarkFailed()
			return fmt.Errorf("Fatal error during fetch %s, to layer %s!", s.repo.Entity, s.repo.Layer)
		}
	}

	dataDecompress, err := gzipDecompress(rawData)

	payload, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		manifest.MarkFailed()
		return fmt.Errorf("Error during converting DTO to bytes.")
	}

	path := repositories.BatchPathJSON(s.repo.Layer, s.repo.Entity, s.repo.Dt)

	return nil
}

func (s *GetLookupStationsService) CleanUp() (error, bool) {
	failedPath := repositories.FailedPath(s.repo.Layer, s.repo.Entity, s.repo.Dt)
	inProgressPath := repositories.InProgressPath(s.repo.Layer, s.repo.Entity, s.repo.Dt)
	batchPath := repositories.BatchPathJSON(s.repo.Layer, s.repo.Entity, s.repo.Dt)

	failedState, err := s.s3Client.Exists(failedPath)
	inProgressState, err := s.s3Client.Exists(inProgressPath)
	batchState, err := s.s3Client.Exists(batchPath)

	if err != nil {
		return err, true
	}

	if !failedState {
		return nil, true
	} else {
		err = s.s3Client.Delete(failedPath)
		if err != nil {
			return err, true
		}
	}

	if inProgressState {
		err = s.s3Client.Delete(inProgressPath)
		if err != nil {
			return err, true
		}
	}

	if batchState {
		err = s.s3Client.Delete(batchPath)
		if err != nil {
			return err, true
		}
	}

	return nil, false
}
