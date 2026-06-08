package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	dto "DataLake/internal/infrastructure/gios/dto"
	"DataLake/internal/infrastructure/s3"
	"DataLake/internal/repositories"
	"DataLake/internal/repositories/bronze"
	"DataLake/internal/repositories/silver"
	"DataLake/internal/repositories/silver/schemas"
)

type LoaderService struct {
	Client s3.Client
}

func NewLoaderService() LoaderService {
	return LoaderService{
		Client: *s3.New(),
	}
}

func (l *LoaderService) LoadLeatestAqIndexesFromBronze(rangeInDays int) ([][]dto.AirQualityIndexesDTO, error) {
	env := bronze.SetupMeasurements("")
	leatestDate := l.SearchLeatestDate(env.Layer, env.Entity)
	if leatestDate == "" {
		return nil, fmt.Errorf("Leatest date of %s layer, %s entity wasn't found", env.Layer, env.Entity)
	}
	env.Dt = leatestDate

	var paths []string
	startDate, err := time.Parse("2006/01/02", env.Dt)

	if err != nil {
		return nil, fmt.Errorf("Error has occured during date parsing.")
	}

	for i := range rangeInDays {
		for j := range 24 {
			dt := startDate.AddDate(0, 0, -i).Format("2006/01/02") + "/" + strconv.Itoa(j)
			succesPath := repositories.SuccessPath(env.Layer, env.Entity, dt)
			pathIsExist, err := l.Client.Exists(succesPath)

			if err != nil {
				fmt.Println("Error has occured during checking file existing.")
			}

			if pathIsExist {
				paths = append(paths, repositories.BatchPathJSON(env.Layer, env.Entity, dt))
			}
		}
	}

	var data [][]dto.AirQualityIndexesDTO
	for _, sourcePath := range paths {
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

				var aqIndex []dto.AirQualityIndexesDTO
				if err := json.Unmarshal(dataDecompress, &aqIndex); err != nil {
					return nil, err
				}
				data = append(data, aqIndex)
				break
			}
		}
	}

	return data, nil
}

func (l *LoaderService) LoadLeatestMeasurementsFromBronze(rangeInDays int) ([][]dto.MeasurementDTO, error) {
	env := bronze.SetupMeasurements("")
	leatestDate := l.SearchLeatestDate(env.Layer, env.Entity)
	if leatestDate == "" {
		return nil, fmt.Errorf("Leatest date of %s layer, %s entity wasn't found", env.Layer, env.Entity)
	}
	env.Dt = leatestDate

	var paths []string
	startDate, err := time.Parse("2006/01/02", env.Dt)

	if err != nil {
		return nil, fmt.Errorf("Error has occured during date parsing.")
	}

	for i := range rangeInDays {
		dt := startDate.AddDate(0, 0, -i).Format("2006/01/02")
		succesPath := repositories.SuccessPath(env.Layer, env.Entity, dt)
		pathIsExist, err := l.Client.Exists(succesPath)

		if err != nil {
			fmt.Println("Error has occured during checking file existing.")
		}

		if pathIsExist {
			paths = append(paths, repositories.BatchPathJSON(env.Layer, env.Entity, dt))
		}
	}

	var data [][]dto.MeasurementDTO
	for _, sourcePath := range paths {
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

				var measurements []dto.MeasurementDTO
				if err := json.Unmarshal(dataDecompress, &measurements); err != nil {
					return nil, err
				}
				data = append(data, measurements)
				break
			}
		}
	}

	return data, nil
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

func (l *LoaderService) LoadLeatestLookupStation() (schemas.StationIds, error) {
	env := silver.SetupReferencesStationIds("")
	leatestDate := l.SearchLeatestDate(env.Layer, env.Entity)
	if leatestDate == "" {
		return schemas.StationIds{}, fmt.Errorf("Leatest date of %s layer, %s entity wasn't found", env.Layer, env.Entity)
	}
	env.Dt = leatestDate

	sourcePath := repositories.PathJson(env.Layer, env.Entity, env.Dt, "stationsList")
	breakCounter := 0

	for {
		rawData, err := l.Client.Get(sourcePath)
		if err != nil && breakCounter < 3 {
			breakCounter++
			// sometyhing to log
			continue
		} else if err != nil {
			return schemas.StationIds{}, fmt.Errorf("Leatest file of %s layer, %s entity wasn't found", env.Layer, env.Entity)
		} else {
			var stationRaw schemas.StationIds
			if err := json.Unmarshal(rawData, &stationRaw); err != nil {
				return schemas.StationIds{}, err
			}

			return stationRaw, nil
		}
	}
}

func (l *LoaderService) LoadLeatestLookupSensors() (schemas.SensorIds, error) {
	env := silver.SetupReferencesSensorIds("")
	leatestDate := l.SearchLeatestDate(env.Layer, env.Entity)
	if leatestDate == "" {
		return schemas.SensorIds{}, fmt.Errorf("Leatest date of %s layer, %s entity wasn't found", env.Layer, env.Entity)
	}
	env.Dt = leatestDate

	sourcePath := repositories.PathJson(env.Layer, env.Entity, env.Dt, "sensorsList")
	breakCounter := 0

	for {
		rawData, err := l.Client.Get(sourcePath)
		if err != nil && breakCounter < 3 {
			breakCounter++
			// sometyhing to log
			continue
		} else if err != nil {
			return schemas.SensorIds{}, fmt.Errorf("Leatest file of %s layer, %s entity wasn't found", env.Layer, env.Entity)
		} else {
			var sensorRaw schemas.SensorIds
			if err := json.Unmarshal(rawData, &sensorRaw); err != nil {
				return schemas.SensorIds{}, err
			}

			return sensorRaw, nil
		}
	}
}

func (l *LoaderService) LoadLeatestMapSensorIdSensorCode() (schemas.MapSensorIdSensorCode, error) {
	env := silver.SetupReferencesSensorIds("")
	leatestDate := l.SearchLeatestDate(env.Layer, env.Entity)
	if leatestDate == "" {
		return schemas.MapSensorIdSensorCode{}, fmt.Errorf("Leatest date of %s layer, %s entity wasn't found", env.Layer, env.Entity)
	}
	env.Dt = leatestDate

	sourcePath := repositories.PathJson(env.Layer, env.Entity, env.Dt, "mapSensorIdToSensorCode")
	breakCounter := 0

	for {
		rawData, err := l.Client.Get(sourcePath)
		if err != nil && breakCounter < 3 {
			breakCounter++
			// sometyhing to log
			continue
		} else if err != nil {
			return schemas.MapSensorIdSensorCode{}, fmt.Errorf("Leatest file of %s layer, %s entity wasn't found", env.Layer, env.Entity)
		} else {
			var mapSensorRaw schemas.MapSensorIdSensorCode
			if err := json.Unmarshal(rawData, &mapSensorRaw); err != nil {
				return schemas.MapSensorIdSensorCode{}, err
			}

			return mapSensorRaw, nil
		}
	}
}
