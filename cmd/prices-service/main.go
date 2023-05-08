package main

import (
	"context"
	"fmt"

	"net/http"
	"time"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
	mongoEnv "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"
	"github.com/gorilla/websocket"

	"github.com/asaskevich/EventBus"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/prices"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/stream"
)

func main() {
	mongoConfig, err := mongoEnv.LoadMongoDBConfig()
	if err != nil {
		panic(err)
	}

	mongoClient, err := database.Connect(mongoConfig, database.Remote)
	if err != nil {
		panic(err)
	}

	cryptoRepo := database.NewCollection[database.CryptoPrices](mongoClient, mongoConfig.Database, "CryptoPrices")
	stocksRepo := database.NewCollection[database.StockPrices](mongoClient, mongoConfig.Database, "StockPrices")

	repoController := prices.NewRepositoryController(cryptoRepo, stocksRepo)

	bus := EventBus.New()
	repoController.ListenForStoring(bus)

	wsConfig, err := env.LoadWebSocetConfig()
	if err != nil {
		panic(err)
	}

	cryptoStreamConfig := stream.NewCryptoConfig(wsConfig)
	stockStreamConfig := stream.NewStockConfig(wsConfig)

	cryptoStream, err := stream.NewPriceStream(cryptoStreamConfig)
	if err != nil {
		panic(err)
	}

	stockStream, err := stream.NewPriceStream(stockStreamConfig)
	if err != nil {
		panic(err)
	}

	streamController := stream.NewController(cryptoStream, stockStream, bus)
	streamController.StartStreamsToWrite()

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	pricesPresenter := prices.NewPresenter(upgrader, bus)

	http.HandleFunc("/crypto", pricesPresenter.CryptoHandler)
	http.HandleFunc("/stocks", pricesPresenter.StockHandler)

	serverConfig, err := env.LoadServerConfig()
	if err != nil {
		panic(err)
	}

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", serverConfig.Port), nil); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	time.Sleep(time.Minute)

	streamController.StopStreams()
	if err := mongoClient.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}
