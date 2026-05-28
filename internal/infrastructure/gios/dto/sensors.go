package gios

type SensorByIdDTO struct {
	Sensors    []SensorDTO `json:"Lista stanowisk pomiarowych dla podanej stacji"`
	Links      LinksDTO    `json:"links"`
	TotalPages int         `json:"totalPages"`
}

type SensorDTO struct {
	SensorID         int    `json:"Identyfikator stanowiska"`
	StationID        int    `json:"Identyfikator stacji"`
	Indicator        string `json:"Wskaźnik"`
	IndicatorFormula string `json:"Wskaźnik - wzór"`
	IndicatorCode    string `json:"Wskaźnik - kod"`
	IndicatorID      int    `json:"Id wskaźnika"`
}

type SensorMetadataDTO struct {
	Metadata   []SensorDetailsDTO `json:"Lista metadanych stanowisk pomiarowych"`
	Links      LinksDTO           `json:"links"`
	TotalPages int                `json:"totalPages"`
}

type SensorDetailsDTO struct {
	SensorNumber    int    `json:"Nr"`
	SensorCode      string `json:"Kod stanowiska"`
	StationName     string `json:"Kod stacji"`
	OldStationName  string `json:"Stary Kod stacji"`
	IndictorCode    string `json:"Wskaźnik - kod"`
	Indicator       string `json:"Wskaźnik"`
	AveragingTime   string `json:"Czas uśredniania"`
	MeasurementType string `json:"Typ pomiaru"`
	StartDate       string `json:"Data uruchomienia"`
	EndDate         string `json:"Data zamknięcia"`
	Voivodeship     string `json:"Województwo"`
	FieldName       string `json:"Nazwa strefy"`
}
