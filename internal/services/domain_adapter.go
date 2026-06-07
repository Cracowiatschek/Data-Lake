package services

import (
	"DataLake/internal/domain"
	"strconv"
	"strings"
	"regexp"
	"DataLake/internal/repositories"
	"DataLake/internal/repositories/silver/schemas"
	"DataLake/internal/repositories/silver"
	"fmt"
	"time"
	"sort"
	"encoding/json"
	"DataLake/internal/infrastructure/s3"
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
	s3Client *s3.Client
}

func (d *DomainAdapterService) MakeAdaptation() ([]domain.Station, error) {
	var stationsIds schemas.StationIds

	leatestStation := d.GetLeatestLookupStationDate()
	if leatestStation == "" {
		fmt.Println("Go to exit. Lookup is empty.")
		return nil,nil // to do 
	}
	// fmt.Println(leatestStation)
	lookupStation, err := d.GetLookupStations(leatestStation)
	if err != nil {
		return nil, fmt.Errorf("Service didn't find correct lookup stations dataset.")
	}
	if err := json.Unmarshal(lookupStation, &stationsIds); err != nil {
		return nil, fmt.Errorf("Error was happend during unpacking lookup station.")
	}
	
	var sensorIds schemas.SensorIds
	var 
	
	if d.Sensors || d.Measurements {
		leatestSensor := d.GetLeatestLookupSensorsDate()
		if leatestSensor == "" {
			fmt.Println("Go to exit. Lookup is empty.")
			return nil, nil
		}
		// fmt.Println(leatestSensor)
		lookupSensor, err := d.GetLookupSensors(leatestSensor)
		if err != nil {
			return nil, fmt.Errorf("Service didn't find correct lookup sensor dataset.")
		}
		if err := json.Unmarshal(lookupSensor, &sensorIds); err != nil {
			return nil, fmt.Errorf("Error was happend during unpacking lookup sensors.")
		}
	}

	switch d.SourceLayer {
	case "bronze":

		return nil,nil
	case "silver":
		return nil, nil
	default:
		return nil,nil
	}
}

func (d *DomainAdapterService) AdaptStationFromBronze(StationID int, BasicData dto.StationDTO, DetailsData dto.StationDetailsDTO, AqIndexes []domain.AqIndex, Sensors []domain.Sensor) domain.Station {
	latitude, err := strconv.ParseFloat(BasicData.Latitude, 32)
	if err != nil {
		fmt.Printf("Error occured during latitude parsing. StationID: %d, Value: %s", StationID, BasicData.Latitude)
	}
	longitude, err := strconv.ParseFloat(BasicData.Longitude, 32)
	if err != nil {
		fmt.Printf("Error occured during longitude parsing. StationID: %d, Value: %s", StationID, BasicData.Longitude)
	}

	startDate, err := time.Parse("12/31/2003", DetailsData.StartDate)
	if err != nil {
		fmt.Printf("Error occured during startDate parsing. StationID: %d, Value: %s", StationID, DetailsData.StartDate)
	}
	var endDate time.Time
	if len(DetailsData.EndDate) > 0 {
		endDate, err = time.Parse("12/31/2003", DetailsData.EndDate)
		if err != nil {
			fmt.Printf("Error occured during endDate parsing. StationID: %d, Value: %s", StationID, DetailsData.EndDate)
		}
	} else {
		endDate = time.Date(2999, time.December, 31, 23, 59, 59, 0, time.UTC)
	}
	return domain.Station{
		StationID:    StationID,
		Code:         BasicData.StationCode,
		Name:         BasicData.StationName,
		Latitude:     float32(latitude),
		Longitude:    float32(longitude),
		City:         BasicData.City,
		District:     BasicData.Distirct,
		Municipality: BasicData.Municipality,
		Voivodeship:  BasicData.Voivodeship,
		Street:       BasicData.Street,
		StartDate:    startDate,
		EndDate:      endDate,
		Type:         DetailsData.StationType,
		FieldType:    DetailsData.FieldType,
		Category:     DetailsData.StationCategory,
		Sensors:      Sensors,
		AqIndexes:    AqIndexes,
	}
}

func (d *DomainAdapterService) AdaptSensorFromBronze(StationID int, MapSesnorIdToCode map[int]string, BasicData []dto.SensorDTO, DetailsData map[string]dto.SensorDetailsDTO, MeasuermentsData []dto.MeasurementValueDTO) []domain.Sensor {
	var result []domain.Sensor
	errors := 0

	for _, basicRecord := range BasicData {
		var measurementsRecords []domain.Measurement
		var endDate time.Time
		sensorCode := MapSesnorIdToCode[basicRecord.SensorID]
		detailsRecords := DetailsData[sensorCode]
		if d.Measurements {
			measurementsRecords = d.AdaptMeasurementFromBronze(basicRecord.SensorID, MeasuermentsData)
		} else {
			measurementsRecords = nil
		}
		startDate, err := time.Parse("12/31/2003", detailsRecords.StartDate)
		if err != nil {
			fmt.Printf("Error occured during startDate parsing. SensorID: %d, Value: %s", basicRecord.SensorID, detailsRecords.StartDate)
		}
		if len(detailsRecords.EndDate) > 0 {
			endDate, err = time.Parse("12/31/2003", detailsRecords.EndDate)
			if err != nil {
				fmt.Printf("Error occured during endDate parsing. SensorID: %d, Value: %s", basicRecord.SensorID, detailsRecords.EndDate)
			}
		} else {
			endDate = time.Date(2999, time.December, 31, 23, 59, 59, 0, time.UTC)
		}

		if err != nil {
			errors++
			continue
		}
		record := domain.Sensor{
			SensorID:         basicRecord.SensorID,
			Indicator:        basicRecord.Indicator,
			IndicatorFormula: basicRecord.IndicatorFormula,
			IndicatorCode:    basicRecord.IndicatorFormula,
			Name:             sensorCode,
			AveragingTime:    detailsRecords.AveragingTime,
			MeasurementType:  detailsRecords.MeasurementType,
			StartDate:        startDate,
			EndDate:          endDate,
			Measurements:     measurementsRecords,
		}

		result = append(result, record)

	}
	fmt.Printf("Station %d has %d errors in Sensors records.", StationID, errors)

	if len(result) == 0 {
		return nil
	}

	return result
}

func (d *DomainAdapterService) AdaptMeasurementFromBronze(SensorID int, Records []dto.MeasurementValueDTO) []domain.Measurement {
	var result []domain.Measurement
	errors := 0

	for _, i := range Records {
		eventDatetime, err := time.Parse("2026-05-28 18:20:24", i.Date)
		if err != nil {
			fmt.Printf("Error occured during eventDatetime parsing. SensorID: %d, Value: %s", SensorID, i.Date)
		}
		if err != nil {
			errors++
			continue
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

func (s *DomainAdapterService) GetLookupStations(dt string) ([]byte, error) {
	env := silver.SetupReferencesStationIds(dt)
	path := repositories.PathJson(env.Layer, env.Entity, env.Dt, "stationsList")
	fmt.Println(path)
	data, err := s.s3Client.Get(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *DomainAdapterService) GetLeatestLookupStationDate() string {
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

func (s *DomainAdapterService) GetLookupSensors(dt string) ([]byte, error) {
	env := silver.SetupReferencesSensorIds(dt)
	path := repositories.PathJson(env.Layer, env.Entity, env.Dt, "sensorsList")
	fmt.Println(path)
	data, err := s.s3Client.Get(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *DomainAdapterService) GetLeatestLookupSensorsDate() string {
	env := silver.SetupReferencesSensorIds("")
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
