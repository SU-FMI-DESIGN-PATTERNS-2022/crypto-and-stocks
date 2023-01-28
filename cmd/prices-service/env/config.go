package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
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
	Driver   string
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func LoadMongoConfig() MongoConfig {
	host := goDotEnvVariable("MONGO_HOST")
	port, err := strconv.Atoi(goDotEnvVariable("MONGO_PORT"))
	if err != nil {
		panic(err)
	}
	driver := goDotEnvVariable("MONGO_DRIVER")
	user := goDotEnvVariable("MONGO_USER")
	database := goDotEnvVariable("MONGO_DATABASE")
	password := goDotEnvVariable("MONGO_PASSWORD")

	return MongoConfig{
		Driver:   driver,
		Host:     host,
		Port:     port,
		User:     user,
		Database: database,
		Password: password,
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
