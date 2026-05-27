package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	dto "DataLake/internal/infrastructure/gios/dto"
	"DataLake/internal/infrastructure/s3"
	"DataLake/internal/repositories"
	"DataLake/internal/repositories/bronze"
)

type LoaderService struct {
	Client s3.Client
}

func (l *LoaderService) LoadLeatestAqIndexesFromBronze(rangeInDays int) ([]dto.AirQualityIndexesDTO, error) {
	return nil, nil
}

func (l *LoaderService) LoadLeatestMeasurementsFromBronze(rangeInDays int) ([]dto.MeasurementDTO, error) {
	return nil, nil
}

func (l *LoaderService) LoadLeatestStationFromBronze() ([]dto.StationFindAllDTO, error) {
	env := bronze.SetupStations("")
	leatestDate := l.SearchLeatestDate(env.Layer, env.Entity)
	if leatestDate == "" {
		return nil, fmt.Errorf("Leatest date of %s layer, %s entity wasn't found", env.Layer, env.Entity)
	}
	env.Dt = leatestDate

	sourcePath := repositories.BatchPathJSON(env.Layer, env.Entity, env.Dt)
	breakCounter := 0

	for {
		rawData, err := l.Client.Get(sourcePath)
		if err != nil && breakCounter < 3 {
			breakCounter++
			// sometyhing to log
			continue
		} else if err != nil {
			return nil, fmt.Errorf("Leatest file of %s layer, %s entity wasn't found", env.Layer, env.Entity)
		} else {
			dataDecompress, err := gzipDecompress(rawData)

			if err != nil {
				return nil, fmt.Errorf("Error during decompressing data.")
			}

			var stationsRaw []dto.StationFindAllDTO
			if err := json.Unmarshal(dataDecompress, &stationsRaw); err != nil {
				return nil, err
			}

			return stationsRaw, nil
		}
	}
}

func (l *LoaderService) LoadLeatestStationDetailsFromBronze() ([]dto.StationMetadataDTO, error) {
	env := bronze.SetupStationDetails("")
	leatestDate := l.SearchLeatestDate(env.Layer, env.Entity)
	if leatestDate == "" {
		return nil, fmt.Errorf("Leatest date of %s layer, %s entity wasn't found", env.Layer, env.Entity)
	}
	env.Dt = leatestDate

	sourcePath := repositories.BatchPathJSON(env.Layer, env.Entity, env.Dt)
	breakCounter := 0

	for {
		rawData, err := l.Client.Get(sourcePath)
		if err != nil && breakCounter < 3 {
			breakCounter++
			// sometyhing to log
			continue
		} else if err != nil {
			return nil, fmt.Errorf("Leatest file of %s layer, %s entity wasn't found", env.Layer, env.Entity)
		} else {
			dataDecompress, err := gzipDecompress(rawData)

			if err != nil {
				return nil, fmt.Errorf("Error during decompressing data.")
			}

			var stationsRaw []dto.StationMetadataDTO
			if err := json.Unmarshal(dataDecompress, &stationsRaw); err != nil {
				return nil, err
			}

			return stationsRaw, nil
		}
	}
}

func (l *LoaderService) LoadLeatestSensorFromBronze() ([]dto.SensorByIdDTO, error) {
	env := bronze.SetupSensors("")
	leatestDate := l.SearchLeatestDate(env.Layer, env.Entity)
	if leatestDate == "" {
		return nil, fmt.Errorf("Leatest date of %s layer, %s entity wasn't found", env.Layer, env.Entity)
	}
	env.Dt = leatestDate

	sourcePath := repositories.BatchPathJSON(env.Layer, env.Entity, env.Dt)
	breakCounter := 0

	for {
		rawData, err := l.Client.Get(sourcePath)
		if err != nil && breakCounter < 3 {
			breakCounter++
			// sometyhing to log
			continue
		} else if err != nil {
			return nil, fmt.Errorf("Leatest file of %s layer, %s entity wasn't found", env.Layer, env.Entity)
		} else {
			dataDecompress, err := gzipDecompress(rawData)

			if err != nil {
				return nil, fmt.Errorf("Error during decompressing data.")
			}

			var sensorsRaw []dto.SensorByIdDTO
			if err := json.Unmarshal(dataDecompress, &sensorsRaw); err != nil {
				return nil, err
			}

			return sensorsRaw, nil
		}
	}
}

func (l *LoaderService) LoadLeatestSensorDetailsFromBronze() ([]dto.SensorMetadataDTO, error) {
	env := bronze.SetupSensorDetails("")
	leatestDate := l.SearchLeatestDate(env.Layer, env.Entity)
	if leatestDate == "" {
		return nil, fmt.Errorf("Leatest date of %s layer, %s entity wasn't found", env.Layer, env.Entity)
	}
	env.Dt = leatestDate

	sourcePath := repositories.BatchPathJSON(env.Layer, env.Entity, env.Dt)
	breakCounter := 0

	for {
		rawData, err := l.Client.Get(sourcePath)
		if err != nil && breakCounter < 3 {
			breakCounter++
			// sometyhing to log
			continue
		} else if err != nil {
			return nil, fmt.Errorf("Leatest file of %s layer, %s entity wasn't found", env.Layer, env.Entity)
		} else {
			dataDecompress, err := gzipDecompress(rawData)

			if err != nil {
				return nil, fmt.Errorf("Error during decompressing data.")
			}

			var sensorsRaw []dto.SensorMetadataDTO
			if err := json.Unmarshal(dataDecompress, &sensorsRaw); err != nil {
				return nil, err
			}

			return sensorsRaw, nil
		}
	}
}

func (l *LoaderService) SearchLeatestDate(layer, entity string) string {
	prefix := fmt.Sprintf("%s/%s/", layer, entity)

	keys, err := l.Client.List(prefix)
	if err != nil {
		return ""
	}

	re := regexp.MustCompile(`dt=([0-9]{4}/[0-9]{2}/[0-9]{2})`)

	var dates []string

	for _, key := range keys {

		if strings.HasSuffix(key, "_SUCCESS") {
			match := re.FindStringSubmatch(key)

			if len(match) > 1 {
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
