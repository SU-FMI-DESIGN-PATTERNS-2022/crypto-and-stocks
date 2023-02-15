package main

import (
	"context"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"
)

// func cryptoHandler(b []byte) {
// 	var cryptoResponse []stream.CryptoResponse
// 	if err := json.Unmarshal(b, &cryptoResponse); err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println(cryptoResponse)
// }

// func stockHandler(b []byte) {
// 	var stockResponse []stream.StockResponse
// 	if err := json.Unmarshal(b, &stockResponse); err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println(stockResponse)
// }

func main() {
	mongoConfig := env.LoadMongoConfig()
	client, err := database.Connect(mongoConfig, database.Remote)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	_ = database.NewCollection[database.CryptoPrices](client, mongoConfig.Database, "CryptoPrices")
	_ = database.NewCollection[database.StockPrices](client, mongoConfig.Database, "StockPrices")

	// cryptoPricesCollection.StoreEntry(database.CryptoPrices{
	// 	Prices: database.Prices{
	// 		Symbol:   "BTC",
	// 		BidPrice: 15728.36,
	// 		BidSize:  0.0472,
	// 		AskPrice: 15701.92,
	// 		AskSize:  0.0453,
	// 		Date:     time.Now(),
	// 	},
	// 	Exchange: "Nexo",
	// })
	// fmt.Println(cryptoPricesCollection.GetAllPrices())

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

	// time.Sleep(time.Second * 3)

	// cryptoStream.Stop()
	// stockStream.Stop()
}
