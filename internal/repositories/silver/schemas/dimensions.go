package schemas

import "time"

type Station struct {
	StationId                int
	StationCode              string
	StationInternationalCode string
	StationOldCode           string
	StationName              string
	StationType              string
	StationFieldType         string
	StationCategory          string
	Latitude                 float64
	Longitude                float64
	CityId                   int
	CityName                 string
	Municipality             string
	District                 string
	Voivodeship              string
	Address                  string
	StartDate                time.Time
	EndDate                  time.Time
}
