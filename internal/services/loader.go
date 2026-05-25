package services

import (
	dto "DataLake/internal/infrastructure/gios/dto"
	"DataLake/internal/infrastructure/s3"
)

type LoaderService struct {
	Client s3.Client
}

func (l *LoaderService) LoadLeatestAqIndexes(rangeInDays int) ([]dto.AirQualityIndexesDTO, error) {
	return nil, nil
}

func (l *LoaderService) LoadLeatestMeasurements(rangeInDays int) ([]dto.MeasurementDTO, error) {
	return nil, nil
}

func (l *LoaderService) LoadLeatestStation() ([]dto.StationFindAllDTO, error) {
	return nil, nil
}

func (l *LoaderService) LoadLeatestStationDetails() ([]dto.StationMetadataDTO, error) {
	return nil, nil
}

func (l *LoaderService) LoadLeatestSensor() ([]dto.SensorByIdDTO, error) {
	return nil, nil
}

func (l *LoaderService) LoadLeatestSensorDetails() ([]dto.StationMetadataDTO, error) {
	return nil, nil
}
