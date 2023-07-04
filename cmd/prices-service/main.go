package main

import (
	"context"
	"fmt"
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

type upgrader struct {
	wsUpgarder *websocket.Upgrader
}

func (u *upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (prices.Connection, error) {
	return u.wsUpgarder.Upgrade(w, r, responseHeader)
}

func main() {
	mongoConfig, err := mongoEnv.LoadMongoDBConfig()
	if err != nil {
		panic(err)
	}

	mongoClient, err := database.Connect(mongoConfig, database.Remote)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

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

	go func() {
		if err := streamController.StartStreamsToWrite(); err != nil {
			log.Fatal(err)
		}
	}()
	defer streamController.StopStreams()

	wsUpgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	pricesPresenter := prices.NewPresenter(&upgrader{wsUpgrader}, bus)

	http.HandleFunc("/crypto", pricesPresenter.CryptoHandler)
	http.HandleFunc("/stocks", pricesPresenter.StockHandler)

	serverConfig, err := env.LoadServerConfig()

	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", serverConfig.Port), nil))
}
