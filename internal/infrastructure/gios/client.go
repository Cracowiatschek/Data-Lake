// internal/infrastructure/gios/client.go
package gios

import (
	"encoding/json"
	"fmt"
	"time"

	dto "DataLake/internal/infrastructure/gios/dto"
	httpclient "DataLake/internal/infrastructure/http"
)

type Client struct {
	http  *httpclient.Client
	base  string
	retry int
	timer int
	size  int
}

func New(http *httpclient.Client, timerMs, retryTimes int, pageSize int) *Client {
	return &Client{
		http:  http,
		base:  "https://api.gios.gov.pl/pjp-api/v1/rest",
		retry: retryTimes,
		timer: timerMs,
		size:  pageSize,
	}
}

func (c *Client) FetchStations(page int) (dto.StationFindAllDTO, error) {
	url := fmt.Sprintf("%s/station/findAll?page=%d&size=%d", c.base, page, c.size)

	var body []byte
	var err error
	for range c.retry {
		body, err = c.http.Get(url)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(c.timer) * time.Millisecond)
		}
	}

	if err != nil {
		return dto.StationFindAllDTO{}, err
	}

	var result dto.StationFindAllDTO
	if err := json.Unmarshal(body, &result); err != nil {
		return dto.StationFindAllDTO{}, err
	}

	return result, nil
}

func (c *Client) FetchStationsDetails(page int) (dto.StationMetadataDTO, error) {
	url := fmt.Sprintf("%s/metadata/stations?page=%d&size=%d", c.base, page, c.size)

	var body []byte
	var err error
	for range c.retry {
		body, err = c.http.Get(url)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(c.timer) * time.Millisecond)
		}
	}

	if err != nil {
		return dto.StationMetadataDTO{}, err
	}

	var result dto.StationMetadataDTO
	if err := json.Unmarshal(body, &result); err != nil {
		return dto.StationMetadataDTO{}, err
	}

	return result, nil
}

func (c *Client) FetchSensor(stationId int) (dto.SensorByIdDTO, error) {
	url := fmt.Sprintf("%s/station/sensors/%d?size=%d", c.base, stationId, c.size)

	var body []byte
	var err error
	for range c.retry {
		body, err = c.http.Get(url)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(c.timer) * time.Millisecond)
		}
	}

	if err != nil {
		return dto.SensorByIdDTO{}, err
	}

	var result dto.SensorByIdDTO
	if err := json.Unmarshal(body, &result); err != nil {
		return dto.SensorByIdDTO{}, err
	}

	return result, nil
}

func (c *Client) FetchSensorDetails(page int) (dto.SensorMetadataDTO, error) {
	url := fmt.Sprintf("%s/metadata/sensors?page=%d&size=%d", c.base, page, c.size)

	var body []byte
	var err error
	for range c.retry {
		body, err = c.http.Get(url)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(c.timer) * time.Millisecond)
		}
	}

	if err != nil {
		return dto.SensorMetadataDTO{}, err
	}

	var result dto.SensorMetadataDTO
	if err := json.Unmarshal(body, &result); err != nil {
		return dto.SensorMetadataDTO{}, err
	}

	return result, nil
}

func (c *Client) FetchAirQualityIndexes(stationId int) (dto.AirQualityIndexesDTO, error) {
	url := fmt.Sprintf("%s/aqindex/getIndex/%d?size=50", c.base, stationId)
	// fmt.Println(url)
	var body []byte
	var err error
	for range c.retry {
		body, err = c.http.Get(url)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(c.timer) * time.Millisecond)
		}
	}

	if err != nil {
		return dto.AirQualityIndexesDTO{}, err
	}

	var result dto.AirQualityIndexesDTO
	if err := json.Unmarshal(body, &result); err != nil {
		return dto.AirQualityIndexesDTO{}, err
	}

	return result, nil
}

func (c *Client) FetchMeasurement(sensorId int, page int) (dto.MeasurementDTO, error) {
	url := fmt.Sprintf("%s/data/getData/%d?page=%d&size=%d", c.base, sensorId, page, c.size)

	var body []byte
	var err error
	for range c.retry {
		body, err = c.http.Get(url)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(c.timer) * time.Millisecond)
		}
	}

	if err != nil {
		return dto.MeasurementDTO{}, err
	}

	var result dto.MeasurementDTO
	if err := json.Unmarshal(body, &result); err != nil {
		return dto.MeasurementDTO{}, err
	}

	return result, nil
}
