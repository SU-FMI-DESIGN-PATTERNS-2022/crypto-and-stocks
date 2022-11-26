package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/prices"
)

func streamHandler(b []byte) {
	var cryptoResponse []prices.CryptoResponse
	if err := json.Unmarshal(b, &cryptoResponse); err != nil {
		fmt.Println(err)
	}

	fmt.Println(cryptoResponse)
}

func main() {
	wsConfig := env.LoadWebSocetConfig()
	cryptoStreamConfig := prices.StreamConfig{
		URL:    wsConfig.CryptoURL,
		Quotes: wsConfig.CryptoQuotes,
		Key:    wsConfig.Key,
		Secret: wsConfig.Secret,
	}

	cryptoStream, err := prices.NewStream(cryptoStreamConfig)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := cryptoStream.Start(streamHandler); err != nil {
			panic(err)
		}
	}()

	time.Sleep(4 * time.Second)
	cryptoStream.Stop()
}
