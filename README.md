# Data Lake Project

## Overview

Data Lake is a personal data engineering project that ingests air quality data from the GIOŚ API and stores them using a multi-layer architecture:

- Bronze (raw ingestion)
- Silver (normalized datasets)
- Gold (analytical layer)

The project is written in Go and follows a layered architecture inspired by modern data lake and lakehouse solutions.

<!-- ## Project Goals -->

## Architecture

```ASCII
GIOŚ API
    │
    ▼
HTTP Client
    │
    ▼
Services
    │
    ▼
Repositories
    │
    ▼
S3 Bronze
    │
    ▼
Transform
    │
    ▼
Domain Models
    │
    ▼
S3 Silver
    │
    ▼
Postgres Gold
```

## Data Flow

### Fetch Job

```ASCII
API
 ↓
DTO
 ↓
Bronze JSON.GZ
```

### Transform Job

```ASCII
Bronze JSON.GZ
 ↓
DTO
 ↓
Domain
 ↓
Silver Parquet
```

### Load Job

```ASCII
Silver Parquet
 ↓
Postgres
```

## Layers

### Bronze

**Goal**
Immutable raw data storage.

**Format**
`json.gz`

**Example**

```ASCII
bronze/
└── stations/
    └── dt=2026-06-11/
        ├── stations.json.gz
        ├── _SUCCESS
        └── _MANIFEST.json
```

### Silver

**Goal**
Normalized and analytics-ready datasets.

**Format**
`parquet`

**Example**

```ASCII
silver/
└── stations/
    └── dt=2026-06-11/
        ├── stations.parquet
        ├── _SUCCESS
        └── _MANIFEST.json
```

### Gold

**Goal**
Serving layer for reporting and analytics.

**Backend**
PostgreSQL

## Project Structure

```ASCII
cmd/
internal/
├── domain/
├── services/
├── repositories/
├── infrastructure/
│   ├── gios/
│   ├── http/
│   ├── s3/
│   └── postgres/
```

## Entities

### Stations

### Sensors

### AQ Indexes

### Measurements

## Pipeline Jobs

| Job                | Purpose               |
| ------------------ | --------------------- |
| fetch_stations     | Download station list |
| fetch_sensors      | Download sensors      |
| fetch_aq_indexes   | Download AQ indexes   |
| fetch_measurements | Download measurements |

### fetch_stations TBD

### fetch_sensors TBD

### fetch_aq_indexes TBD

### fetch_measurements TBD

### transform_stations TBD

### transform_sensors TBD

### transform_aq_indexes TBD

### transform_measurements TBD

## Storage Layout

Convention nameing: `{layer}/{entity}/dt={yyyy-mm-dd}/`

### Bronze Example

`bronze/stations/dt=2026-06-11/`

### Silver Example

`silver/sensors/dt=2026-06-11/`

### Gold Example - TBD

## Markers

### `_SUCCESS`

Batch completed successfully.

### `_FAILED`

Batch failed.

### `_INPROGRESS`

Batch currently running.

## Manifest

### Raw

```JSON
{
  "requests": [],
  "pages": 0,
  "records": 0,
  "created_at": "",
  "processed_time": 0
}
```

### Transform

```JSON
{
  "files": [],
  "records": 0,
  "created_at": "",
  "processed_time": 0
}
```

### Gold Ingestion

TBD

## Environment Variables

- `S3_ACCESS_KEY`
- `S3_SECRET_KEY`
- `S3_ENDPOINT`
- `S3_REGION`
- `S3_BUCKET`

## Running

Examples:

- `go run cmd/fetch_stations/main.go`
- `go run cmd/transform_stations/main.go`

### Local TBD

### Docker TBD

### Cron TBD

## Future Improvements

## Roadmap

- Bronze → Silver transformations
- Parquet writers
- Gold loading
- Sensor lookup datasets
- Historical backfill
- Data quality checks
- Metrics and monitoring
- Docker deployment
