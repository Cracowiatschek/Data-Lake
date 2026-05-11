package bronze

import (
	s3Service "DataLake/internal/infrastructure/s3"
	repo "DataLake/internal/repositories"
)

func readFile(env Env, s3 *s3Service.Client) ([]byte, error) {
	path := repo.BatchPathJSON(env.Layer, env.Entity, env.Dt)

	bytes, err := s3.Get(path)

	if err != nil {
		return nil, err
	}
	return bytes, nil
}
