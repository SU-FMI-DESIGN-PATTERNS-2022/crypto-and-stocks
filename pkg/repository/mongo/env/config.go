package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type MongoConfig struct {
	LocalDriver  string
	RemoteDriver string
	Host         string
	Port         string
	Database     string
	User         string
	Password     string
	Options      string
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load("../../pkg/repository/mongo/env/.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func LoadMongoConfig() MongoConfig {
	host := goDotEnvVariable("MONGO_HOST")
	port := goDotEnvVariable("MONGO_PORT")
	localDriver := goDotEnvVariable("MONGO_LOCAL_DRIVER")
	remoteDriver := goDotEnvVariable("MONGO_REMOTE_DRIVER")
	user := goDotEnvVariable("MONGO_USER")
	database := goDotEnvVariable("MONGO_DATABASE")
	password := goDotEnvVariable("MONGO_PASSWORD")
	options := goDotEnvVariable("MONGO_OPTIONS")

	return MongoConfig{
		LocalDriver:  localDriver,
		RemoteDriver: remoteDriver,
		Host:         host,
		Port:         port,
		User:         user,
		Database:     database,
		Password:     password,
		Options:      options,
	}
}
