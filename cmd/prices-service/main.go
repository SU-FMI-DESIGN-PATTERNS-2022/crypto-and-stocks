package main

import (
	"context"

	"flag"
	"log"
	"net/http"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
	mongoEnv "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"
	"github.com/gorilla/websocket"

	"github.com/asaskevich/EventBus"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/prices"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/stream"
)

var addr = flag.String("addr", "localhost:8000", "http service address")

func main() {
	mongoConfig := mongoEnv.LoadMongoConfig()

	client, err := database.Connect(mongoConfig, database.Remote)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	cryptoRepo := database.NewCollection[database.CryptoPrices](client, mongoConfig.Database, "CryptoPrices")
	stocksRepo := database.NewCollection[database.StockPrices](client, mongoConfig.Database, "StockPrices")

	repoController := prices.NewRepositoryController(cryptoRepo, stocksRepo)

	bus := EventBus.New()
	repoController.ListenForStoring(bus)

	wsConfig := env.LoadWebSocetConfig()
	cryptoStreamConfig := stream.NewStreamConfig(wsConfig)
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

	go func() {
		if err := streamController.StartStreamsToWrite(); err != nil {
			log.Fatal(err)
		}
	}()
	defer streamController.StopStreams()

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	pricesPresenter := prices.NewPresenter(upgrader, bus)

	http.HandleFunc("/crypto", pricesPresenter.CryptoHandler)
	http.HandleFunc("/stocks", pricesPresenter.StockHandler)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
