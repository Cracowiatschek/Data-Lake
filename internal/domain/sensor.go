package domain

import "time"

type Sensor struct {
	SensorID int
	// StationID        int
	Indicator        string
	IndicatorFormula string
	IndicatorCode    string
	Name             int
	AveragingTime    string
	MeasurementType  string
	StartDate        time.Time
	EndDate          time.Time
	Measurements     []Measurement
}
