package bronze

import (
	s3Service "DataLake/internal/infrastructure/s3"
	repo "DataLake/internal/repositories"
)

func saveTempFile(bytes []byte, env Env, s3 *s3Service.Client, page int) error {
	layer := "tmp_" + env.Layer
	path := repo.PageBasedPathJSON(layer, env.Entity, env.Dt, page)

	err := s3.Put(path, bytes)

	if err != nil {
		return err
	}
	return nil
}

func saveFile(bytes []byte, env Env, s3 *s3Service.Client) error {
	path := repo.BatchPathJSON(env.Layer, env.Entity, env.Dt)

	err := s3.Put(path, bytes)

	if err != nil {
		return err
	}
	return nil
}
