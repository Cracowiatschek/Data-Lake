package bronze

type Env struct {
	Layer  string
	Entity string
	Dt     string
}

func SetupStations(dt string) Env {
	return Env{
		Layer:  "bronze",
		Entity: "stations",
		Dt:     dt,
	}
}

func SetupStationDetails(dt string) Env {
	return Env{
		Layer:  "bronze",
		Entity: "stationDetails",
		Dt:     dt,
	}
}

func SetupSensors(dt string) Env {
	return Env{
		Layer:  "bronze",
		Entity: "sensors",
		Dt:     dt,
	}
}

func SetupSensorDetails(dt string) Env {
	return Env{
		Layer:  "bronze",
		Entity: "sensorsDetails",
		Dt:     dt,
	}
}

func SetupAirQualityIndexes(dt string) Env {
	return Env{
		Layer:  "bronze",
		Entity: "airQualityIndexes",
		Dt:     dt,
	}
}

func SetupMeasurements(dt string) Env {
	return Env{
		Layer:  "bronze",
		Entity: "measurements",
		Dt:     dt,
	}
}
