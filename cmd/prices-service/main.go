package main

import (
	"context"

	"flag"
	"log"
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

type upgrader struct {
	wsUpgarder *websocket.Upgrader
}

func (u *upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (prices.Connection, error) {
	return u.wsUpgarder.Upgrade(w, r, responseHeader)
}

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	mongoConfig := mongoEnv.LoadMongoConfig()

	client, err := database.Connect(mongoConfig, database.Remote)
	if err != nil {
		panic(err)
	}

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
	streamController.StartStreamsToWrite()

	wsUpgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	pricesPresenter := prices.NewPresenter(&upgrader{wsUpgrader}, bus)

	http.HandleFunc("/crypto", pricesPresenter.CryptoHandler)
	http.HandleFunc("/stocks", pricesPresenter.StockHandler)

	go log.Fatal(http.ListenAndServe(*addr, nil))

	time.Sleep(time.Minute)
	streamController.StopStreams()
	if err = client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
