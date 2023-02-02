package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type WebSocetConfig struct {
	CryptoURL    string
	StockURL     string
	CryptoQuotes []string
	StockQuotes  []string
	Key          string
	Secret       string
}

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

func LoadWebSocetConfig() WebSocetConfig {
	key := goDotEnvVariable("KEY")
	secret := goDotEnvVariable("SECRET")
	return WebSocetConfig{
		CryptoURL:    "wss://stream.data.alpaca.markets/v1beta1/crypto",
		StockURL:     "wss://stream.data.alpaca.markets/v2/iex",
		CryptoQuotes: []string{"BTCUSD", "ETHUSD", "ADAUSD", "DOTUSD", "USDTUSD", "SOLUSD", "MATICUSD", "LINKUSD", "ATOMUSD", "BMBUSD", "LTCUSD"},
		StockQuotes:  []string{"AAPL", "AMZN"},
		Key:          key,
		Secret:       secret,
	}
}
