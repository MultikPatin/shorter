package adapters

import (
	"go.uber.org/zap"
	"main/internal/adapters/database/memory"
	"main/internal/adapters/database/sql_db"
	"main/internal/config"
	"main/internal/services"
)

func GetDatabase(c *config.Config, logger *zap.SugaredLogger) (services.DataBase, error) {
	var database services.DataBase
	var err error

	if c.PostgresDNS == nil {
		database, err = memory.NewInMemoryDB(c.StorageFilePaths, logger)
		if err != nil {
			return nil, err
		}
	} else {
		database, err = sql_db.NewPostgresDB(c.PostgresDNS, logger)
		if err != nil {
			return nil, err
		}
	}
	return database, nil
}
