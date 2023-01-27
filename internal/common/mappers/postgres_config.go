package mappers

import (
	"nnchat/internal/config"
	"nnchat/pkg/database/postgres"
)

func MapPostgresConfig(pgConfig config.PostgresConfig) postgres.Config {
	return postgres.Config{
		Host:     pgConfig.Host,
		Port:     pgConfig.Port,
		DBName:   pgConfig.DBName,
		User:     pgConfig.User,
		Password: pgConfig.Password,
	}
}
