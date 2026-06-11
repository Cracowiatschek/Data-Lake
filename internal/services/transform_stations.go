package services

import (
	// "DataLake/internal/domain"
	"DataLake/internal/domain"
	"DataLake/internal/infrastructure/s3"
	"fmt"
)

type TransformStationsService struct {
	Client  *s3.Client
	Adapter DomainAdapterService
}

func NewTransformStationsService() TransformStationsService {
	return TransformStationsService{
		Client: s3.New(),
		Adapter: DomainAdapterService{
			SourceLayer:       "bronze",
			Stations:          true,
			Sensors:           false,
			AqIndex:           false,
			AqIndexRange:      0,
			Measurements:      false,
			MeasurementsRange: 0,
		},
	}
}

func (t *TransformStationsService) TransformStations() error {
	domainData, err := t.Adapter.MakeAdaptation()
	if err != nil {
		fmt.Println("Error occures during adapting data for Station.")
		return err
	}

	return nil
}

func (t *TransformStationsService) BuildParquet(data []domain.Station) error {
	return nil
}
