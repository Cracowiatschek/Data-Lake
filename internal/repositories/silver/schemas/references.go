package schemas

type StationIds struct {
	StationId []int `json:"stationId"`
}

type SensorIds struct {
	SensorId []int `json:"sensorId"`
}

type MapSensorIdSensorCode struct {
	Mapping []Sensor
}

type Sensor struct {
	SensorId   int    `json:"sensorId"`
	SensorCode string `json:"sensorCode"`
}
