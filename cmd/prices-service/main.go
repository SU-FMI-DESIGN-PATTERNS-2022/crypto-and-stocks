package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"
	"log"

	"net/http"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
	mongoEnv "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/asaskevich/EventBus"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/prices"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/stream"
)

type upgrader struct {
	wsUpgrader *websocket.Upgrader
}

func (u *upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (prices.Connection, error) {
	return u.wsUpgrader.Upgrade(w, r, responseHeader)
}

func main() {
	mongoConfig, err := mongoEnv.LoadMongoDBConfig()
	if err != nil {
		log.Fatal(err)
	}

	mongoClient, err := database.Connect(mongoConfig, database.Remote)
	if err != nil {
		log.Fatal(err)
	}

	wsConfig, err := env.LoadWebSocetConfig()
	if err != nil {
		log.Fatal(err)
	}

	cryptoStreamConfig := stream.NewCryptoConfig(wsConfig)
	stockStreamConfig := stream.NewStockConfig(wsConfig)

	cryptoStream, err := stream.NewPriceStream(cryptoStreamConfig)
	if err != nil {
		log.Fatal(err)
	}

	stockStream, err := stream.NewPriceStream(stockStreamConfig)
	if err != nil {
		log.Fatal(err)
	}

	cryptoRepo := database.NewCollection[database.CryptoPrices](mongoClient, mongoConfig.Database, "CryptoPrices")
	stocksRepo := database.NewCollection[database.StockPrices](mongoClient, mongoConfig.Database, "StockPrices")

	repoController := prices.NewRepositoryController(cryptoRepo, stocksRepo)
	bus := EventBus.New()
	streamController := stream.NewController(cryptoStream, stockStream, bus)

	if err := repoController.ListenForStoring(bus); err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := streamController.StartStreamsToWrite(); err != nil {
			log.Fatal(err)
		}
	}()

	wsUpgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	pricesPresenter := prices.NewPresenter(&upgrader{wsUpgrader}, bus)

	router := gin.Default()
	prices.HandleRoutes(&router.RouterGroup, pricesPresenter)

	serverConfig, err := env.LoadServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverConfig.Port),
		Handler: router,
	}

	log.Println("Starting server...")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")

	streamController.StopStreams()

	contextWithTimout, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()

	if err = mongoClient.Disconnect(contextWithTimout); err != nil {
		log.Fatal(err)
	}

	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
}
