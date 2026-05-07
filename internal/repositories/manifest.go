package repositories

type Manifest struct {
	Records       int    `json:"records"`
	Layer         string `json:"layer"`
	CreatedAt     string `json:"created_at"`
	ProcessedTime int    `json:"processed_time"`
	Source        string `json:"source"`
}

type ManifestBronze struct {
	Manifest
	Requests []string `json:"requests"`
	Pages    int      `json:"pages"`
	Dt       string   `json:"dt"`
	Endpoint string   `json:"endpoint"`
}

type ManifestSilver struct {
	Manifest
	Files         []string `json:"files"`
	SourceRecords int      `json:"source_records"`
}

type ManifestGold struct {
	Manifest
	Table         string `json:"table"`
	RowsInserted  int    `json:"rows_inserted"`
	RowsUpdated   int    `json:"rows_updated"`
	SourceRecords int    `json:"source_records"`
}
