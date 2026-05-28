package domain

import "time"

type Measurement struct {
	// SensorID      int
	// StationID     int
	EventDatetime time.Time
	EventValue    float32
}
