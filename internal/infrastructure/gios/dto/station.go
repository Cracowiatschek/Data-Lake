package gios

type StationFindAllDTO struct {
	Stations   []StationDTO `json:"Lista stacji pomiarowych"`
	Links      LinksDTO     `json:"links"`
	TotalPages int          `json:"totalPages"`
}

type StationDTO struct {
	StationID    int    `json:"Identyfikator stacji"`
	StationCode  string `json:"Kod stacji"`
	StationName  string `json:"Nazwa stacji"`
	Latitude     string `json:"WGS84 φ N"`
	Longitude    string `json:"WGS84 λ E"`
	CityID       int    `json:"Identyfikator miasta"`
	City         string `json:"Nazwa miasta"`
	Distirct     string `json:"Gmina"`
	Municipality string `json:"Powiat"`
	Voivodeship  string `json:"Województwo"`
	Street       string `json:"Ulica"`
}

type StationMetadataDTO struct {
	Metadata   []StationDetailsDTO `json:"Lista stacji pomiarowych"`
	Links      LinksDTO            `json:"links"`
	TotalPages int                 `json:"totalPages"`
}

type StationDetailsDTO struct {
	StationNumber           int    `json:"Nr"`
	StationCode             string `json:"Kod stacji"`
	InternationaStationCode string `json:"Kod międzynarodowy"`
	StationName             string `json:"Nazwa stacji"`
	OldStationCode          string `json:"Stary Kod stacji"`
	StartDate               string `json:"Data uruchomienia"`
	EndDate                 string `json:"Data zamknięcia"`
	StationType             string `json:"Typ stacji"`
	FieldType               string `json:"Typ obszaru"`
	StationCategory         string `json:"Rodzaj stacji"`
	Latitude                string `json:"WGS84 φ N"`
	Longitude               string `json:"WGS84 λ E"`
	City                    string `json:"Miejscowość"`
	Voivodeship             string `json:"Województwo"`
	Address                 string `json:"Adres"`
}
