package services

import (
	dto "DataLake/internal/infrastructure/gios/dto"
	"DataLake/internal/infrastructure/s3"
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
	return nil, nil
}

func (l *LoaderService) LoadLeatestStationDetailsFromBronze() ([]dto.StationMetadataDTO, error) {
	return nil, nil
}

func (l *LoaderService) LoadLeatestSensorFromBronze() ([]dto.SensorByIdDTO, error) {
	return nil, nil
}

func (l *LoaderService) LoadLeatestSensorDetailsFromBronze() ([]dto.StationMetadataDTO, error) {
	return nil, nil
}
