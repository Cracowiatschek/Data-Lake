package repositories

import (
	"DataLake/internal/infrastructure/s3"
	"encoding/json"
)

type ManifestRepository struct {
	Client s3.Client
}

func (r *ManifestRepository) SaveManifest(layer, entity, dt string, manifest any) error {
	data, err := json.MarshalIndent(manifest, "", " ")
	if err != nil {
		return err
	}

	key := manifestPath(layer, entity, dt)

	return r.Client.Put(key, data)
}

func (r *ManifestRepository) MarkSuccess(layer, entity, dt string) error {
	key := successPath(layer, entity, dt)
	return r.Client.Put(key, []byte{})
}

func (r *ManifestRepository) MarkFailed(layer, entity, dt string) error {
	key := failedPath(layer, entity, dt)
	return r.Client.Put(key, []byte{})
}

func (r *ManifestRepository) MarkInProgress(layer, entity, dt string) error {
	key := inProgressPath(layer, entity, dt)
	return r.Client.Put(key, []byte{})
}

func (r *ManifestRepository) ClearInProgress(layer, entity, dt string) error {
	key := inProgressPath(layer, entity, dt)
	return r.Client.Delete(key)
}
