package adapters

import (
	"go.uber.org/zap"
	"main/internal/adapters/database/memory"
	"main/internal/adapters/database/psql"
	"main/internal/config"
	"main/internal/services"
)

func NewLinksRepository(c *config.Config, logger *zap.SugaredLogger) (services.LinksRepository, error) {
	var repository services.LinksRepository
	var err error

	if c.PostgresDNS == nil {
		repository, err = memory.NewInMemoryRepository(c.StorageFilePaths, logger)
		if err != nil {
			return nil, err
		}
	} else {
		repository, err = psql.NewPostgresRepository(c.PostgresDNS, logger)
		if err != nil {
			return nil, err
		}
	}
	return repository, nil
}
