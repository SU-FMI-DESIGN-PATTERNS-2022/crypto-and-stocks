package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/env"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/database"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/repositories/crypto_prices_repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/stream"
)

func cryptoHandler(b []byte) {
	var cryptoResponse []stream.CryptoResponse
	if err := json.Unmarshal(b, &cryptoResponse); err != nil {
		fmt.Println(err)
	}
	//TODO: Make a method that saves the response into the corresponding collection
	fmt.Println(cryptoResponse)
}

func stockHandler(b []byte) {
	var stockResponse []stream.StockResponse
	if err := json.Unmarshal(b, &stockResponse); err != nil {
		fmt.Println(err)
	}
	//TODO: same here
	fmt.Println(stockResponse)
}

func main() {
	mongoConfig := env.LoadMongoConfig()
	client, err := database.Connect(mongoConfig, database.Remote)

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	cryptoPricesCollection := crypto_prices_repository.NewCryptoPricesCollection(client, mongoConfig.Database, "CryptoPrices")
	// stockPricesCollection := stock_prices_repository.NewStockPricesCollection(client, mongoConfig.Database, "StockPrices")
	// pricesPresenter := prices.NewPricesPresenter(cryptoPricesCollection, stockPricesCollection)
	fmt.Println(cryptoPricesCollection.GetAllPrices())
	// cryptoPricesCollection.StoreEntry(crypto_prices_repository.CryptoPrices{
	// 	Prices: database.Prices{
	// 		Symbol:   "BTC",
	// 		BidPrice: 13452.23,
	// 		BidSize:  0.0024,
	// 		AskPrice: 13452.23,
	// 		AskSize:  0.0024,
	// 		Date:     time.Now(),
	// 	},
	// 	Exchange: "Binance",
	// })

	//===============================================================================================
	// wsConfig := env.LoadWebSocetConfig()
	// cryptoStreamConfig := stream.StreamConfig{
	// 	URL:    wsConfig.CryptoURL,
	// 	Quotes: wsConfig.CryptoQuotes,
	// 	Key:    wsConfig.Key,
	// 	Secret: wsConfig.Secret,
	// }

	// stockStreamConfig := stream.StreamConfig{
	// 	URL:    wsConfig.StockURL,
	// 	Quotes: wsConfig.StockQuotes,
	// 	Key:    wsConfig.Key,
	// 	Secret: wsConfig.Secret,
	// }

	// cryptoStream, err := stream.NewPriceStream(cryptoStreamConfig)
	// if err != nil {
	// 	panic(err)
	// }
	// stockStream, err := stream.NewPriceStream(stockStreamConfig)
	// if err != nil {
	// 	panic(err)
	// }

	// go func() {
	// 	if err := cryptoStream.Start(cryptoHandler); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// go func() {
	// 	if err := stockStream.Start(stockHandler); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// time.Sleep(time.Minute)
	// cryptoStream.Stop()
	// stockStream.Stop()
}
