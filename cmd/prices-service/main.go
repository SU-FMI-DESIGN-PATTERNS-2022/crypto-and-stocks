package main

import (
	"context"

	"log"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
	mongoEnv "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/asaskevich/EventBus"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/prices"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/stream"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	mongoConfig := mongoEnv.LoadMongoConfig()
	wsConfig := env.LoadWebSocetConfig()
	serverConfig := env.LoadServerConfig()
	cryptoStreamConfig := stream.NewStreamConfig(wsConfig)
	stockStreamConfig := stream.NewStockConfig(wsConfig)

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

	cryptoStream, err := stream.NewPriceStream(cryptoStreamConfig)
	if err != nil {
		panic(err)
	}

	stockStream, err := stream.NewPriceStream(stockStreamConfig)
	if err != nil {
		panic(err)
	}

	repoController := prices.NewRepositoryController(cryptoRepo, stocksRepo)
	bus := EventBus.New()
	streamController := stream.NewController(cryptoStream, stockStream, bus)

	repoController.ListenForStoring(bus)
	errCh := streamController.StartStreamsToWrite()

	pricesPresenter := prices.NewPresenter(upgrader, bus)

	router := gin.Default()
	prices.HandleRoutes(&router.RouterGroup, *pricesPresenter)

	select {
	case err := <-errCh:
		log.Fatal(err)
		streamController.StopStreams()
		return
	default:
		router.Run("localhost:" + serverConfig.Port)
	}
}
