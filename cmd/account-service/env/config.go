package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load("./env/.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func LoadDBConfig() DBConfig {
	host := goDotEnvVariable("HOST")
	port, err := strconv.Atoi(goDotEnvVariable("PORT"))
	if err != nil {
		panic(err)
	}
	dbuser := goDotEnvVariable("DBUSER")
	password := goDotEnvVariable("PASSWORD")
	dbname := goDotEnvVariable("DBNAME")
	return DBConfig{
		Host:     host,
		Port:     port,
		User:     dbuser,
		Password: password,
		DBName:   dbname,
	}
}
