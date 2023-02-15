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

func goDotEnvVariable(key string) string {
	err := godotenv.Load("./env/.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
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
