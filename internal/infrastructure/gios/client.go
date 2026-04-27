// internal/infrastructure/gios/client.go
package gios

import (
	"encoding/json"
	"fmt"

	// "DataLake/internal/domain"
	dto "DataLake/internal/infrastructure/gios/dto"
	httpclient "DataLake/internal/infrastructure/http"
)

type Client struct {
	http *httpclient.Client
	base string
}

func New(http *httpclient.Client) *Client {
	return &Client{
		http: http,
		base: "https://api.gios.gov.pl/pjp-api/v1/rest",
	}
}

func (c *Client) FetchStations() (dto.StationFindAllDTO, error) {
	url := fmt.Sprintf("%s/station/findAll", c.base)

	body, err := c.http.Get(url)
	if err != nil {
		return dto.StationFindAllDTO{}, err
	}

	var result dto.StationFindAllDTO
	if err := json.Unmarshal(body, &result); err != nil {
		return dto.StationFindAllDTO{}, err
	}

	return result, nil
}
