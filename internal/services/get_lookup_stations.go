package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	dto "DataLake/internal/infrastructure/gios/dto"
	"DataLake/internal/infrastructure/s3"
	"DataLake/internal/repositories"
	"DataLake/internal/repositories/bronze"
	"DataLake/internal/repositories/silver"
	"DataLake/internal/repositories/silver/schemas"
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
	sourceEnv := bronze.SetupStations(s.repo.Dt)
	sourceDt := s.SearchLeatestExistingSourceDate(sourceEnv.Layer, sourceEnv.Entity)
	if sourceDt != "" {
		sourceEnv.Dt = sourceDt
	}
	sourcePath := repositories.BatchPathJSON(sourceEnv.Layer, sourceEnv.Entity, sourceEnv.Dt)

	err, goToExit := s.CleanUp()
	if err != nil {
		return fmt.Errorf("Fatal error during cleaning up past job")
	}

	if goToExit {
		// log something about skipped job
		return nil
	}

	manifest.MarkInProgress()

	breakCounter := 0
	var rawData []byte

	for {
		rawData, err = s.s3Client.Get(sourcePath)

		if err != nil && breakCounter < 3 {
			breakCounter++
			// sometyhing to log
			continue
		} else if err != nil {
			manifest.MarkFailed()
			return fmt.Errorf("Fatal error during fetch %s, to layer %s!", s.repo.Entity, s.repo.Layer)
		}
		if err == nil {
			break
		}
	}

	dataDecompress, err := gzipDecompress(rawData)

	if err != nil {
		manifest.MarkFailed()
		return fmt.Errorf("Error during decompressing data.")
	}

	var stationsRaw []dto.StationFindAllDTO
	if err := json.Unmarshal(dataDecompress, &stationsRaw); err != nil {
		manifest.MarkFailed()
		return err
	}

	var stationIds []int

	for _, record := range stationsRaw {
		for _, station := range record.Stations {
			stationIds = append(stationIds, station.StationID)
		}
	}
	rawRecords := len(stationsRaw)
	computedRecords := len(stationIds)
	stationsResult := schemas.StationIds{
		StationId: stationIds,
	}

	path := repositories.PathJson(s.repo.Layer, s.repo.Entity, s.repo.Dt, "stationsList.json")
	payload, err := json.MarshalIndent(stationsResult, "", " ")
	err = s.s3Client.Put(path, payload)

	if err != nil {
		manifest.MarkFailed()
		return fmt.Errorf("Error during saving %s.", path)
	}
	var files []string
	files = append(files, sourcePath)

	silverManifest := repositories.ManifestSilver{
		Files:         files,
		SourceRecords: rawRecords,
		Manifest: repositories.Manifest{
			Records:       computedRecords,
			Layer:         s.repo.Layer,
			CreatedAt:     time.Now().UTC().Format(time.RFC3339),
			ProcessedTime: int(time.Since(start).Milliseconds()),
			Source:        "Źródło danych: GIOŚ - EKOINFONET",
		},
	}

	manifest.SaveManifest(silverManifest)
	err = manifest.ClearInProgress()
	if err != nil {
		manifest.MarkFailed()
		return fmt.Errorf("Error occurence when in progress mark was deleting.")
	}
	manifest.MarkSuccess()

	return nil
}

func (s *GetLookupStationsService) CleanUp() (error, bool) {
	failedPath := repositories.FailedPath(s.repo.Layer, s.repo.Entity, s.repo.Dt)
	inProgressPath := repositories.InProgressPath(s.repo.Layer, s.repo.Entity, s.repo.Dt)
	batchPath := repositories.PathJson(s.repo.Layer, s.repo.Entity, s.repo.Dt, "stationsList.json")

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

func (s *GetLookupStationsService) SearchLeatestExistingSourceDate(layer, entity string) string {
	prefix := fmt.Sprintf("%s/%s/", layer, entity)

	keys, err := s.s3Client.List(prefix)
	if err != nil {
		return ""
	}

	re := regexp.MustCompile(`dt=([0-9]{4}-[0-9]{2}-[0-9]{2})`)

	var dates []string

	for _, key := range keys {
		if !strings.HasSuffix(key, "_SUCCESS") {
			continue
		}

		match := re.FindStringSubmatch(key)

		if len(match) > 1 {
			dates = append(dates, match[1])
		}
	}

	if len(dates) == 0 {
		return ""
	}

	sort.Strings(dates)

	return dates[len(dates)-1]
}
