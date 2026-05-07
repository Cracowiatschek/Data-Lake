package bronze

type Env struct {
	layer  string
	entity string
	dt     string
	page   string
}

func SetupStations(dt, page string) Env {
	return Env{
		layer:  "bronze",
		entity: "stations",
		dt:     dt,
		page:   page,
	}
}

func SetupStationDetails(dt, page string) Env {
	return Env{
		layer:  "bronze",
		entity: "stationDetails",
		dt:     dt,
		page:   page,
	}
}

func SetupSensors(dt, page string) Env {
	return Env{
		layer:  "bronze",
		entity: "sensors",
		dt:     dt,
		page:   page,
	}
}

func SetupSensorDetails(dt, page string) Env {
	return Env{
		layer:  "bronze",
		entity: "sensorsDetails",
		dt:     dt,
		page:   page,
	}
}

func SetupAirQualityIndexes(dt, page string) Env {
	return Env{
		layer:  "bronze",
		entity: "airQualityIndexes",
		dt:     dt,
		page:   page,
	}
}

func SetupMeasurements(dt, page string) Env {
	return Env{
		layer:  "bronze",
		entity: "measurements",
		dt:     dt,
		page:   page,
	}
}
