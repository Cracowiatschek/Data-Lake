package domain

import "time"

type Station struct {
	StationID    int
	Code         string
	Name         string
	Latitude     float32
	Longitude    float32
	City         string
	District     string
	Municipality string
	Voivodeship  string
	Street       string
	StartDate    time.Time
	EndDate      time.Time
	Type         string
	FieldType    string
	Category     string
	Sensors      []Sensor
	AqIndexes    []AqIndex
}
