package services

import (
	"DataLake/internal/domain"
	// "DataLake/internal/infrastructure/gios"
	// dto "DataLake/internal/infrastructure/gios/dto"
	// httpclient "DataLake/internal/infrastructure/http"
	"DataLake/internal/infrastructure/s3"
)

type DomainAdapterService struct {
	SourceLayer       string //available "bronze" and "silver"
	Stations          bool
	Sensors           bool
	AqIndex           bool
	AqIndexRange      int
	Measurements      bool
	MeasurementsRange int
}

func (d *DomainAdapterService) AdaptStationFromBronze(Client s3.Client) []domain.Station {
	return nil
}

func (d *DomainAdapterService) AdaptSensorFromBronze(StationID int, Client s3.Client) []domain.Sensor {
	return nil
}

func (d *DomainAdapterService) AdaptMeasurementFromBronze(SensorID int, Client s3.Client) []domain.Measurement {
	return nil
}

func (d *DomainAdapterService) AdaptAqIndexFromBronze(StationID int, Client s3.Client) []domain.AqIndex {

	return nil
}
