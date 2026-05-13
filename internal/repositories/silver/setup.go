package silver

type Env struct {
	Layer  string
	Entity string
	Dt     string
}

func SetupReferencesStationIds(dt string) Env {
	return Env{
		Layer:  "silver/references",
		Entity: "stationIds",
		Dt:     dt,
	}
}

func SetupReferencesSensorIds(dt string) Env {
	return Env{
		Layer:  "silver/references",
		Entity: "sensorIds",
		Dt:     dt,
	}
}
func SetupReferencesMapSensorIdToSensorCode(dt string) Env {
	return Env{
		Layer:  "silver/references",
		Entity: "mapSensorIdsSensorCode",
		Dt:     dt,
	}
}
