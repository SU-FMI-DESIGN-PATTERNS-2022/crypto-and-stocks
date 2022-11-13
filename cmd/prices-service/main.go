package main

import (
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/prices"
)

func main() {
	wsConfig := env.LoadWebSocetConfig()
	cryptoClientConfig := prices.ClientSocetConfig{
		URL:    wsConfig.CryptoURL,
		Quotes: wsConfig.CryptoQuotes,
		Key:    wsConfig.Key,
		Secret: wsConfig.Secret,
	}

	cryptoClient, err := prices.NewClientSocket(cryptoClientConfig)
	if err != nil {
		panic(err)
	}

	cryptoClient.Read()
	cryptoClient.Read()
	cryptoClient.Read()
	cryptoClient.Read()
}
