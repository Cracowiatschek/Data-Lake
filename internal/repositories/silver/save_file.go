package silver

import (
	s3Service "DataLake/internal/infrastructure/s3"
	repo "DataLake/internal/repositories"
)

func saveJSONFile(bytes []byte, env Env, s3 *s3Service.Client, filename string) error {
	path := repo.PathJson(env.Layer, env.Entity, env.Dt, filename)

	err := s3.Put(path, bytes)

	if err != nil {
		return err
	}
	return nil
}

func saveParquetFile(bytes []byte, env Env, s3 *s3Service.Client, filename string) error {
	path := repo.PathParquet(env.Layer, env.Entity, env.Dt, filename)

	err := s3.Put(path, bytes)

	if err != nil {
		return err
	}
	return nil
}
