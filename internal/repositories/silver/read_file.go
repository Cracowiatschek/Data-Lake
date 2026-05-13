package silver

import (
	s3Service "DataLake/internal/infrastructure/s3"
	repo "DataLake/internal/repositories"
)

func readJSONFile(env Env, filename string, s3 *s3Service.Client) ([]byte, error) {
	path := repo.PathJson(env.Layer, env.Entity, env.Dt, filename)

	bytes, err := s3.Get(path)

	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func readParquetFile(env Env, filename string, s3 *s3Service.Client) ([]byte, error) {
	path := repo.PathParquet(env.Layer, env.Entity, env.Dt, filename)

	bytes, err := s3.Get(path)

	if err != nil {
		return nil, err
	}
	return bytes, nil
}
