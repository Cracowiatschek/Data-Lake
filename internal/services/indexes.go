package services

import (
	dto "DataLake/internal/infrastructure/gios/dto"
	"DataLake/internal/repositories/silver/schemas"
)

func IndexBronzeStation(data dto.StationFindAllDTO) map[int]dto.StationDTO {
	result := make(map[int]dto.StationDTO)
	for _, row := range data.Stations {
		result[row.StationID] = row
	}
	return result
}

func IndexBronzeStationDetails(data dto.StationMetadataDTO) map[string]dto.StationDetailsDTO {
	result := make(map[string]dto.StationDetailsDTO)
	for _, row := range data.Metadata {
		result[row.StationCode] = row
	}
	return result
}

func IndexBronzeSensor(data dto.SensorByIdDTO) map[int][]dto.SensorDTO {
	result := make(map[int][]dto.SensorDTO)
	for _, row := range data.Sensors {
		result[row.StationID] = append(result[row.StationID], row)
	}
	return result
}

func IndexBronzeSensorDetails(data dto.SensorMetadataDTO) map[string]dto.SensorDetailsDTO {
	result := make(map[string]dto.SensorDetailsDTO)
	for _, row := range data.Metadata {
		result[row.SensorCode] = row
	}
	return result
}

func IndexBronzeAqIndex(data []dto.AirQualityIndexesDTO) map[int][]dto.AQIndexesDTO {
	result := make(map[int][]dto.AQIndexesDTO)
	for _, row := range data {
		result[row.Indexes.StationID] = append(result[row.Indexes.StationID], row.Indexes)
	}
	return result
}

func IndexBronzeMeasurements(data []dto.MeasurementDTO, mapping schemas.MapSensorIdSensorCode) map[int][]dto.MeasurementValueDTO {
	result := make(map[int][]dto.MeasurementValueDTO)

	// Stwórz mapę sensorCode -> sensorId dla szybszego wyszukiwania
	sensorCodeToId := make(map[string]int)
	for _, item := range mapping.Mapping {
		sensorCodeToId[item.SensorCode] = item.SensorId
	}

	for _, measurementSet := range data {
		for _, measurement := range measurementSet.Measurements {
			// Znajdź sensorId na podstawie sensorCode
			if sensorId, exists := sensorCodeToId[measurement.SensorCode]; exists {
				// Dodaj pomiar do odpowiedniej listy dla danego sensorId
				result[sensorId] = append(result[sensorId], measurement)
			}
		}
	}

	return result
}
