package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"

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
	ctx := context.TODO()
	client, cancel, connectErr := mongo.Connect(mongoConfig)
	//var pricesRepo = prices_repository.NewDatabase(client)
	//var pricesPresenter = prices.NewPresenter(pricesRepo)

	if connectErr != nil {
		panic(connectErr)
	}
	defer func() {
		if connectErr = client.Disconnect(ctx); connectErr != nil {
			panic(connectErr)
		}
	}()

	defer mongo.Close(client, ctx, cancel)
	// Checking whether the connection was successful
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged MongoDB.")

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
