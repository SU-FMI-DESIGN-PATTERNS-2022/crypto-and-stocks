package env

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type ServerConfig struct {
	Port int `envconfig:"SERVER_PORT" default:"8080"`
}

type PostgreDBConfig struct {
	Host     string `envconfig:"POSTGRE_HOST" required:"true"`
	Port     int    `envconfig:"POSTGRE_PORT" required:"true"`
	Name     string `envconfig:"POSTGRE_NAME" required:"true"`
	User     string `envconfig:"POSTGRE_USER" required:"true"`
	Password string `envconfig:"POSTGRE_PASSWORD" required:"true"`
}

func LoadServerConfig() (ServerConfig, error) {
	var serverConfig ServerConfig
	if err := envconfig.Process("", &serverConfig); err != nil {
		return ServerConfig{}, fmt.Errorf("failed to load server config from environment: %w", err)
	}

	return serverConfig, nil
}

func LoadPostgreDBConfig() (PostgreDBConfig, error) {
	var postgreDBConfig PostgreDBConfig
	if err := envconfig.Process("", &postgreDBConfig); err != nil {
		return PostgreDBConfig{}, fmt.Errorf("failed to load postgre database config from environment: %w", err)
	}

	return postgreDBConfig, nil
}
