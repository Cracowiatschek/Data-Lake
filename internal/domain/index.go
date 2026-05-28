package domain

import "time"

type AqIndex struct {
	// StationID                       int
	GeneralIndexCalculationDatetime time.Time
	GeneralIndexValue               int
	GeneralIndexCategory            string
	// GeneralIndexSourceDatetime      time.Time
	SO2IndexCalculationDatetime time.Time
	SO2IndexValue               int
	SO2IndexCategory            string
	// SO2IndexSourceDatetime          time.Time
	NO2IndexCalculationDatetime time.Time
	NO2IndexValue               int
	NO2IndexCategory            string
	// NO2IndexSourceDatetime          time.Time
	PM10IndexCalculationDatetime time.Time
	PM10IndexValue               int
	PM10IndexCategory            string
	// PM10IndexSourceDatetime         time.Time
	PM25IndexCalculationDatetime time.Time
	PM25IndexValue               int
	PM25IndexCategory            string
	// PM25IndexSourceDatetime         time.Time
	O3IndexCalculationDatetime time.Time
	O3IndexValue               int
	O3IndexCategory            string
	// O3IndexSourceDatetime           time.Time
	Status bool
	Code   string
}
