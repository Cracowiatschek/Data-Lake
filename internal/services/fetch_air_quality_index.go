package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"DataLake/internal/infrastructure/gios"
	dto "DataLake/internal/infrastructure/gios/dto"
	httpclient "DataLake/internal/infrastructure/http"
	"DataLake/internal/infrastructure/s3"
	"DataLake/internal/repositories"
	"DataLake/internal/repositories/bronze"
	"DataLake/internal/repositories/silver"
	"DataLake/internal/repositories/silver/schemas"
)

type FetchAirQualityIndexesService struct {
	s3Client *s3.Client
	gios     *gios.Client
	repo     bronze.Env
}

func NewAirQualityIndexesService(dt string) FetchAirQualityIndexesService {
	return FetchAirQualityIndexesService{
		s3Client: s3.New(),
		gios:     gios.New(httpclient.New(), 2000, 3, 500),
		repo:     bronze.SetupAirQualityIndexes(dt),
	}
}

func (s *FetchAirQualityIndexesService) Run() error {
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

	var stationsIds schemas.StationIds

	leatestDate := s.GetLeatestLookupStationDate()
	if leatestDate == "" {
		fmt.Println("Go to exit. Lookup is empty.")
		return nil
	}
	fmt.Println(leatestDate)
	lookupData, err := s.GetLookupStations(leatestDate)

	if err := json.Unmarshal(lookupData, &stationsIds); err != nil {
		manifest.MarkFailed()
		return err
	}

	var data []dto.AirQualityIndexesDTO
	var requests []string
	records := 0
	breakCounter := 0

	for _, station := range stationsIds.StationId {
		d, err := s.gios.FetchAirQualityIndexes(station)

		requests = append(requests, d.Links.Self)

		if err != nil && breakCounter < 3 {
			time.Sleep(time.Duration(time.Second * 2))
			breakCounter++
			// sometyhing to log
			continue
		} else if err != nil {
			manifest.MarkFailed()
			return fmt.Errorf("Fatal error during fetch %s, to layer %s!", s.repo.Entity, s.repo.Layer)
		}

		data = append(data, d)

		records++
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
		Pages:    len(requests),
		Dt:       s.repo.Dt,
		Endpoint: "https://api.gios.gov.pl/pjp-api/v1/rest/aqindex/getIndex/",
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

func (s *FetchAirQualityIndexesService) CleanUp() (error, bool) {
	failedPath := repositories.FailedPath(s.repo.Layer, s.repo.Entity, s.repo.Dt)
	inProgressPath := repositories.InProgressPath(s.repo.Layer, s.repo.Entity, s.repo.Dt)
	batchPath := repositories.BatchPathJSON(s.repo.Layer, s.repo.Entity, s.repo.Dt)

	failedState, err := s.s3Client.Exists(failedPath)
	inProgressState, err := s.s3Client.Exists(inProgressPath)
	batchState, err := s.s3Client.Exists(batchPath)

	if err != nil {
		fmt.Println(err)
		return err, true
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

func (s *FetchAirQualityIndexesService) GetLookupStations(dt string) ([]byte, error) {
	env := silver.SetupReferencesStationIds(dt)
	path := repositories.PathJson(env.Layer, env.Entity, env.Dt, "stationsList")
	fmt.Println(path)
	data, err := s.s3Client.Get(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *FetchAirQualityIndexesService) GetLeatestLookupStationDate() string {
	env := silver.SetupReferencesStationIds("")
	prefix := fmt.Sprintf("%s/%s/", env.Layer, env.Entity)

	keys, err := s.s3Client.List(prefix)
	if err != nil {
		return ""
	}

	re := regexp.MustCompile(`dt=([0-9]{4}/[0-9]{2}/[0-9]{2})`)

	var dates []string

	for _, key := range keys {

		if strings.HasSuffix(key, "_SUCCESS") {
			match := re.FindStringSubmatch(key)

			if len(match) > 1 {
				fmt.Println(key)
				dates = append(dates, match[1])
			}
		}
	}

	if len(dates) == 0 {
		return ""
	}

	sort.Strings(dates)

	return dates[len(dates)-1]
}
