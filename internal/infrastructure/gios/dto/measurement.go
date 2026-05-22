package gios

type AirQualityIndexesDTO struct {
	Indexes    AQIndexesDTO `json:"AqIndex"`
	Links      LinksDTO     `json:"links"`
	TotalPages int          `json:"totalPages"`
}

type AQIndexesDTO struct {
	SensorID             int    `json:"Identyfikator stacji pomiarowej"`
	IndexCalculationDate string `json:"Data wykonania obliczeń indeksu"`
	IndexValue           int    `json:"Wartość indeksu"`
	IndexCategory        string `json:"Nazwa kategorii indeksu"`
	// IndexSourceDate          string `json:"Data danych źródłowych, z których policzono wartość indeksu dla wskaźnika st"`
	SO2IndexCalculationDate string `json:"Data wykonania obliczeń indeksu dla wskaźnika SO2"`
	SO2IndexValue           int    `json:"Wartość indeksu dla wskaźnika SO2"`
	SO2IndexCategory        string `json:"Nazwa kategorii indeksu dla wskażnika SO2"`
	// SO2IndexSourceDate       string `json:"Data danych źródłowych, z których policzono wartość indeksu dla wskaźnika SO2"`
	NO2IndexCalculationDate string `json:"Data wykonania obliczeń indeksu dla wskaźnika NO2"`
	NO2IndexValue           int    `json:"Wartość indeksu dla wskaźnika NO2"`
	NO2IndexCategory        string `json:"Nazwa kategorii indeksu dla wskażnika NO2"`
	// NO2IndexSourceDate       string `json:"Data danych źródłowych, z których policzono wartość indeksu dla wskaźnika NO2"`
	PM10IndexCalculationDate string `json:"Data wykonania obliczeń indeksu dla wskaźnika PM10"`
	PM10IndexValue           int    `json:"Wartość indeksu dla wskaźnika PM10"`
	PM10IndexCategory        string `json:"Nazwa kategorii indeksu dla wskażnika PM10"`
	// PM10IndexSourceDate      string `json:"Data danych źródłowych, z których policzono wartość indeksu dla wskaźnika PM10"`
	PM25IndexCalculationDate string `json:"Data wykonania obliczeń indeksu dla wskaźnika PM2.5"`
	PM25IndexValue           int    `json:"Wartość indeksu dla wskaźnika PM2.5"`
	PM25IndexCategory        string `json:"Nazwa kategorii indeksu dla wskażnika PM2.5"`
	// PM25IndexSourceDate      string `json:"Data danych źródłowych, z których policzono wartość indeksu dla wskaźnika PM2.5"`
	O3IndexCalculationDate string `json:"Data wykonania obliczeń indeksu dla wskaźnika O3"`
	O3IndexValue           int    `json:"Wartość indeksu dla wskaźnika O3"`
	O3IndexCategory        string `json:"Nazwa kategorii indeksu dla wskażnika O3"`
	// O3IndexSourceDate        string `json:"Data danych źródłowych, z których policzono wartość indeksu dla wskaźnika O3"`
	IndexStatus  bool   `json:"Status indeksu ogólnego dla stacji pomiarowej"`
	PolutionCode string `json:"Kod zanieczyszczenia krytycznego"`
}

type MeasurementDTO struct {
	Measurements []MeasurementValueDTO `json:"Lista danych pomiarowych"`
	Links        LinksDTO              `json:"links"`
	TotalPages   int                   `json:"totalPages"`
}

type MeasurementValueDTO struct {
	SensorCode string `json:"Kod stanowiska"`
	Date       string `json:"Data"`
	Value      int    `json:"Wartość"`
}
