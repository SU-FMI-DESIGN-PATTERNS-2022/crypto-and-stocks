package env

import (
	"errors"
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type MongoDBConfig struct {
	Host         string `envconfig:"MONGO_HOST" required:"true"`
	Port         int    `envconfig:"MONGO_PORT" required:"true"`
	LocalDriver  string `envconfig:"MONGO_LOCAL_DRIVER"`
	RemoteDriver string `envconfig:"MONGO_REMOTE_DRIVER"`
	User         string `envconfig:"MONGO_USER" required:"true"`
	Database     string `envconfig:"MONGO_DATABASE" required:"true"`
	Password     string `envconfig:"MONGO_PASSWORD" required:"true"`
	Options      string `envconfig:"MONGO_OPTIONS" required:"true"`
}

func LoadMongoDBConfig() (MongoDBConfig, error) {
	var mongoDBConfig MongoDBConfig
	if err := envconfig.Process("", &mongoDBConfig); err != nil {
		return MongoDBConfig{}, fmt.Errorf("failed to load mongo database config from environment: %w", err)
	}

	if mongoDBConfig.LocalDriver == "" && mongoDBConfig.RemoteDriver == "" {
		return MongoDBConfig{}, errors.New("at least one of the drevers(local or remote) should be set")
	}

	return mongoDBConfig, nil
}
