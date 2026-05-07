package bronze

import (
	s3Service "DataLake/internal/infrastructure/s3"
	repo "DataLake/internal/repositories"
)

func saveTempFile(bytes []byte, env Env, s3 *s3Service.Client) error {
	layer := "tmp_" + env.layer
	path := repo.PageBasedPathJSON(layer, env.entity, env.dt, env.page)

	err := s3.Put(path, bytes)

	if err != nil {
		return err
	}
	return nil
}

func saveFile(bytes []byte, env Env, s3 *s3Service.Client) error {
	path := repo.BatchPathJSON(env.layer, env.entity, env.dt)

	err := s3.Put(path, bytes)

	if err != nil {
		return err
	}
	return nil
}
