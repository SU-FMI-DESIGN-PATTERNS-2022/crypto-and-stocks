package env

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type ServerConfig struct {
	Port int `envconfig:"SERVER_PORT" default:"8081"`
}

type WebSocetConfig struct {
	CryptoURL    string `envconfig:"WS_CRYPTO_URL" default:"wss://stream.data.alpaca.markets/v1beta1/crypto"`
	StockURL     string `envconfig:"WS_STOCK_URL" default:"wss://stream.data.alpaca.markets/v2/iex"`
	CryptoQuotes []string
	StockQuotes  []string
	Key          string `envconfig:"WS_KEY" required:"true"`
	Secret       string `envconfig:"WS_SECRET" required:"true"`
}

func LoadServerConfig() (ServerConfig, error) {
	var serverConfig ServerConfig
	if err := envconfig.Process("", &serverConfig); err != nil {
		return ServerConfig{}, fmt.Errorf("failed to load server config from environment: %w", err)
	}

	return serverConfig, nil
}

func LoadWebSocetConfig() (WebSocetConfig, error) {
	var wsConfig WebSocetConfig
	if err := envconfig.Process("", &wsConfig); err != nil {
		return WebSocetConfig{}, fmt.Errorf("failed to load web socket config from environment: %w", err)
	}

	wsConfig.CryptoQuotes = []string{"BTCUSD", "ETHUSD", "ADAUSD", "DOTUSD", "USDTUSD", "SOLUSD", "MATICUSD", "LINKUSD", "ATOMUSD", "BMBUSD", "LTCUSD"}
	wsConfig.StockQuotes = []string{"AAPL", "AMZN"}

	return wsConfig, nil
}
