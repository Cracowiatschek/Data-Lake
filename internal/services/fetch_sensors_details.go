package services

import (
	"encoding/json"
	"fmt"

	// "regexp"
	// "sort"
	// "strings"
	"strconv"
	"time"

	"DataLake/internal/infrastructure/gios"
	dto "DataLake/internal/infrastructure/gios/dto"
	httpclient "DataLake/internal/infrastructure/http"
	"DataLake/internal/infrastructure/s3"
	"DataLake/internal/repositories"
	"DataLake/internal/repositories/bronze"
	// "DataLake/internal/repositories/silver"
	// "DataLake/internal/repositories/silver/schemas"
)

type FetchSensorsDetailsService struct {
	s3Client *s3.Client
	gios     *gios.Client
	repo     bronze.Env
}

func NewFetchSensorsDetailsService(dt string) FetchSensorsDetailsService {
	return FetchSensorsDetailsService{
		s3Client: s3.New(),
		gios:     gios.New(httpclient.New(), 35000, 3, 500),
		repo:     bronze.SetupSensorDetails(dt),
	}
}

func (s *FetchSensorsDetailsService) Run() error {
	start := time.Now()
	manifest := repositories.NewManifestRepository(s.repo.Layer, s.repo.Entity, s.repo.Dt)

	err, goToExit := s.CleanUp()
	if err != nil {
		return fmt.Errorf("Fatal error during cleaning up past job")
	}
	fmt.Println("After CleanUp")
	if goToExit {
		// log something about skipped job
		fmt.Println("Go to exit")
		return nil
	}

	manifest.MarkInProgress()
	page := 0
	records := 0
	breakCounter := 0

	var data []dto.SensorMetadataDTO
	var requests []string

	for {
		d, err := s.gios.FetchSensorDetails(page)

		requests = append(requests, d.Links.Self)

		if (err != nil || len(d.Metadata) == 0) && breakCounter < 3 {
			breakCounter++
			time.Sleep(time.Duration(time.Second * 120))
			// sometyhing to log
			continue
		} else if err != nil || len(d.Metadata) == 0 {
			manifest.MarkFailed()
			return fmt.Errorf("Fatal error during fetch %s, to layer %s!", s.repo.Entity, s.repo.Layer)
		}

		data = append(data, d)

		nextPage := getPageFromAPILink(d.Links.Next)
		selfPage := getPageFromAPILink(d.Links.Self)
		lastPage := getPageFromAPILink(d.Links.Last)
		records += len(d.Metadata)

		if nextPage != "0" && nextPage != selfPage && selfPage != lastPage {
			conversionPage, err := strconv.Atoi(nextPage)

			if err != nil {
				manifest.MarkFailed()
				return fmt.Errorf("Fatal error during fetch %s, to layer %s! Conversion page error!", s.repo.Entity, s.repo.Layer)
			}
			page = conversionPage
			time.Sleep(time.Duration(time.Second * 120))
		} else {
			break
		}
	}

	payload, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		manifest.MarkFailed()
		return fmt.Errorf("Error during converting DTO to bytes.")
	}

	path := repositories.BatchPathJSON(s.repo.Layer, s.repo.Entity, s.repo.Dt)
	payloadAsGzip, err := gzipCompress(payload)

	if err != nil {
		manifest.MarkFailed()
		return fmt.Errorf("Error during Gzip compressioning.")
	}

	err = s.s3Client.Put(path, payloadAsGzip)

	if err != nil {
		manifest.MarkFailed()
		return fmt.Errorf("Error during saving %s.", path)
	}

	bronzeManifest := repositories.ManifestBronze{
		Requests: requests,
		Pages:    page,
		Dt:       s.repo.Dt,
		Endpoint: "https://api.gios.gov.pl/pjp-api/v1/rest/metadata/sensors",
		Manifest: repositories.Manifest{
			Records:       records,
			Layer:         s.repo.Layer,
			CreatedAt:     time.Now().UTC().Format(time.RFC3339),
			ProcessedTime: int(time.Since(start).Milliseconds()),
			Source:        "Źródło danych: GIOŚ - EKOINFONET",
		},
	}

	manifest.SaveManifest(bronzeManifest)
	err = manifest.ClearInProgress()
	if err != nil {
		manifest.MarkFailed()
		return fmt.Errorf("Error occurence when in progress mark was deleting.")
	}
	manifest.MarkSuccess()

	return nil
}

func (s *FetchSensorsDetailsService) CleanUp() (error, bool) {
	failedPath := repositories.FailedPath(s.repo.Layer, s.repo.Entity, s.repo.Dt)
	successPath := repositories.SuccessPath(s.repo.Layer, s.repo.Entity, s.repo.Dt)
	inProgressPath := repositories.InProgressPath(s.repo.Layer, s.repo.Entity, s.repo.Dt)
	batchPath := repositories.BatchPathJSON(s.repo.Layer, s.repo.Entity, s.repo.Dt)

	failedState, err := s.s3Client.Exists(failedPath)
	successState, err := s.s3Client.Exists(successPath)
	inProgressState, err := s.s3Client.Exists(inProgressPath)
	batchState, err := s.s3Client.Exists(batchPath)

	if err != nil {
		fmt.Println(err)
		return err, true
	}

	if successState {
		return nil, false
	}

	if inProgressState && !failedState {
		return nil, false
	}

	if !failedState {
		return nil, false
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
