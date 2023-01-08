package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
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

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func LoadWebSocetConfig() WebSocetConfig {
	key := goDotEnvVariable("KEY")
	return WebSocetConfig{
		CryptoURL:    "wss://stream.data.alpaca.markets/v1beta1/crypto",
		StockURL:     "wss://stream.data.alpaca.markets/v2/iex",
		CryptoQuotes: []string{"BTCUSD", "ETHUSD", "ADAUSD", "DOTUSD", "USDTUSD", "SOLUSD", "MATICUSD", "LINKUSD", "ATOMUSD", "BMBUSD", "LTCUSD"},
		StockQuotes:  []string{"AAPL", "AMZN"},
		Key:          key,
		Secret:       "IYUve7og11abpgu4Qv8MIGCoFTd8HdALaLg6aaZ5",
	}
}
