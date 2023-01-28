package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/stream"
)

func cryptoHandler(b []byte) {
	var cryptoResponse []stream.CryptoResponse
	if err := json.Unmarshal(b, &cryptoResponse); err != nil {
		fmt.Println(err)
	}

	fmt.Println(cryptoResponse)
}

func stockHandler(b []byte) {
	var stockResponse []stream.StockResponse
	if err := json.Unmarshal(b, &stockResponse); err != nil {
		fmt.Println(err)
	}

	fmt.Println(stockResponse)
}

func main() {
	wsConfig := env.LoadWebSocetConfig()
	cryptoStreamConfig := stream.StreamConfig{
		URL:    wsConfig.CryptoURL,
		Quotes: wsConfig.CryptoQuotes,
		Key:    wsConfig.Key,
		Secret: wsConfig.Secret,
	}

	stockStreamConfig := stream.StreamConfig{
		URL:    wsConfig.StockURL,
		Quotes: wsConfig.StockQuotes,
		Key:    wsConfig.Key,
		Secret: wsConfig.Secret,
	}

	cryptoStream, err := stream.NewPriceStream(cryptoStreamConfig)
	if err != nil {
		panic(err)
	}
	stockStream, err := stream.NewPriceStream(stockStreamConfig)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := cryptoStream.Start(cryptoHandler); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := stockStream.Start(stockHandler); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Minute)
	cryptoStream.Stop()
	stockStream.Stop()
}
