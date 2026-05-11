package repositories

import (
	"DataLake/internal/infrastructure/s3"
	"encoding/json"
)

type ManifestRepository struct {
	Client s3.Client
	layer  string
	entity string
	dt     string
}

func NewManifestRepository(layer, entity, dt string) *ManifestRepository {
	return &ManifestRepository{
		Client: *s3.New(),
		layer:  layer,
		entity: entity,
		dt:     dt,
	}
}

func (r *ManifestRepository) SaveManifest(manifest any) error {
	data, err := json.MarshalIndent(manifest, "", " ")
	if err != nil {
		return err
	}

	key := ManifestPath(r.layer, r.entity, r.dt)

	return r.Client.Put(key, data)
}

func (r *ManifestRepository) MarkSuccess() error {
	key := SuccessPath(r.layer, r.entity, r.dt)
	return r.Client.Put(key, []byte{})
}

func (r *ManifestRepository) MarkFailed() error {
	key := FailedPath(r.layer, r.entity, r.dt)
	return r.Client.Put(key, []byte{})
}

func (r *ManifestRepository) MarkInProgress() error {
	key := InProgressPath(r.layer, r.entity, r.dt)
	return r.Client.Put(key, []byte{})
}

func (r *ManifestRepository) ClearInProgress() error {
	key := InProgressPath(r.layer, r.entity, r.dt)
	return r.Client.Delete(key)
}
