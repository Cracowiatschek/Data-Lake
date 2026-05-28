package services

import (
	"DataLake/internal/domain"
	// "DataLake/internal/repositories"
	// "DataLake/internal/repositories/silver/schemas"
	"fmt"
	"time"

	// "DataLake/internal/infrastructure/gios"
	// dto "DataLake/internal/infrastructure/gios/dto"
	// httpclient "DataLake/internal/infrastructure/http"
	dto "DataLake/internal/infrastructure/gios/dto"
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

func (d *DomainAdapterService) AdaptStationFromBronze(StationID int, BasicData map[int][]dto.StationDTO, DetailsData map[int][]dto.StationDetailsDTO, AqIndexes []domain.AqIndex, Sensors []domain.Sensor) domain.Station {

	return domain.Station{}
}

func (d *DomainAdapterService) AdaptSensorFromBronze(StationID int, MapSesnorIdToCode map[int]string, BasicData map[int][]dto.SensorDTO, DetailsData map[string][]dto.SensorDetailsDTO, MeasuermentsData map[int][]dto.MeasurementValueDTO) []domain.Sensor {
	basicSourceRecords := BasicData[StationID]

	var result []domain.Sensor
	errors := 0

	for _, basicRecord := range basicSourceRecords {
		sensorCode := MapSesnorIdToCode[basicRecord.SensorID]
		detailsRecords := DetailsData[sensorCode]
		measurements := MeasuermentsData[basicRecord.SensorID]
		measurementsRecords := d.AdaptMeasurementFromBronze(basicRecord.SensorID, measurements)

		record := domain.Sensor{
			SensorID: basicRecord.SensorID,
		}

	}
	return nil
}

func (d *DomainAdapterService) AdaptMeasurementFromBronze(SensorID int, Records []dto.MeasurementValueDTO) []domain.Measurement {
	var result []domain.Measurement
	errors := 0

	for _, i := range Records {
		eventDatetime, err := time.Parse("2026-05-28 18:20:24", i.Date)
		if err != nil {
			fmt.Printf("Error occured during eventDatetime parsing. SensorID: %d, Value: %s", SensorID, i.Date)
		}

		record := domain.Measurement{
			EventValue:    i.Value,
			EventDatetime: eventDatetime,
		}
		result = append(result, record)
	}
	fmt.Printf("Sensor %d has %d errors in Measurements records.", SensorID, errors)

	if len(result) == 0 {
		return nil
	}

	return result
}

func (d *DomainAdapterService) AdaptAqIndexFromBronze(StationID int, Data map[int][]dto.AQIndexesDTO) []domain.AqIndex {
	sourceRecords := Data[StationID]

	var result []domain.AqIndex
	errors := 0

	for _, i := range sourceRecords {
		generalIndexCalculationDatetime, err := time.Parse("2026-05-28 18:20:24", i.IndexCalculationDate)
		if err != nil {
			fmt.Printf("Error occured during generalIndexCalculationDatetime parsing. StationID: %d, Value: %s", i.StationID, i.IndexCalculationDate)
		}

		so2IndexCalculationDatetime, err := time.Parse("2026-05-28 18:20:24", i.SO2IndexCalculationDate)
		if err != nil {
			fmt.Printf("Error occured during so2IndexCalculationDatetime parsing. StationID: %d, Value: %s", i.StationID, i.SO2IndexCalculationDate)
		}

		no2IndexCalculationDatetime, err := time.Parse("2026-05-28 18:20:24", i.NO2IndexCalculationDate)
		if err != nil {
			fmt.Printf("Error occured during no2IndexCalculationDatetime parsing. StationID: %d, Value: %s", i.StationID, i.NO2IndexCalculationDate)
		}

		pm10IndexCalculationDatetime, err := time.Parse("2026-05-28 18:20:24", i.PM10IndexCalculationDate)
		if err != nil {
			fmt.Printf("Error occured during pm10IndexCalculationDatetime parsing. StationID: %d, Value: %s", i.StationID, i.PM10IndexCalculationDate)
		}

		pm25IndexCalculationDatetime, err := time.Parse("2026-05-28 18:20:24", i.PM25IndexCalculationDate)
		if err != nil {
			fmt.Printf("Error occured during pm25IndexCalculationDatetime parsing. StationID: %d, Value: %s", i.StationID, i.PM25IndexCalculationDate)
		}

		o3IndexCalculationDatetime, err := time.Parse("2026-05-28 18:20:24", i.O3IndexCalculationDate)
		if err != nil {
			fmt.Printf("Error occured during o3IndexCalculationDatetime parsing. StationID: %d, Value: %s", i.StationID, i.O3IndexCalculationDate)
		}

		if err != nil {
			errors++
			continue
		}
		record := domain.AqIndex{
			// StationID:                       i.StationID,
			GeneralIndexCalculationDatetime: generalIndexCalculationDatetime,
			GeneralIndexValue:               i.IndexValue,
			GeneralIndexCategory:            i.IndexCategory,
			SO2IndexCalculationDatetime:     so2IndexCalculationDatetime,
			SO2IndexValue:                   i.SO2IndexValue,
			SO2IndexCategory:                i.SO2IndexCategory,
			NO2IndexCalculationDatetime:     no2IndexCalculationDatetime,
			NO2IndexValue:                   i.NO2IndexValue,
			NO2IndexCategory:                i.NO2IndexCategory,
			PM10IndexCalculationDatetime:    pm10IndexCalculationDatetime,
			PM10IndexValue:                  i.PM10IndexValue,
			PM10IndexCategory:               i.PM10IndexCategory,
			PM25IndexCalculationDatetime:    pm25IndexCalculationDatetime,
			PM25IndexValue:                  i.PM25IndexValue,
			PM25IndexCategory:               i.PM25IndexCategory,
			O3IndexCalculationDatetime:      o3IndexCalculationDatetime,
			O3IndexValue:                    i.O3IndexValue,
			O3IndexCategory:                 i.O3IndexCategory,
			Status:                          i.IndexStatus,
			Code:                            i.PolutionCode,
		}
		result = append(result, record)
	}

	fmt.Printf("Station %d has %d errors in AqIndex records.", StationID, errors)

	if len(result) == 0 {
		return nil
	}

	return result
}
