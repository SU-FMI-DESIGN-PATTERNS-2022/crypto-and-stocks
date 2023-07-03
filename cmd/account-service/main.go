package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/order"
	repository "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/order_repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/user_repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/user"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
	mongoEnv "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type upgrader struct {
	wsUpgrader *websocket.Upgrader
}

func (u *upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (order.Connection, error) {
	return u.wsUpgrader.Upgrade(w, r, responseHeader)
}

func main() {
	dbConfig, err := env.LoadPostgreDBConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.Connect(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	serverConfig, err := env.LoadServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	mongoConfig, err := mongoEnv.LoadMongoDBConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err := database.Connect(mongoConfig, database.Remote)
	if err != nil {
		log.Fatal(err)
	}

	orderRepository := order_repository.NewOrderTable(db)
	userRepository := user_repository.NewUserTable(db)
	cryptoRepository := database.NewCollection[database.CryptoPrices](client, mongoConfig.Database, "CryptoPrices")
	stockRepository := database.NewCollection[database.StockPrices](client, mongoConfig.Database, "StockPrices")

	wsUpgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	orderPresenter := order.NewOrderPresenter(orderRepository, userRepository, &upgrader{wsUpgrader})
	userPresenter := user.NewUserPresenter(orderRepository, userRepository, cryptoRepository, stockRepository)

	router := gin.Default()
	ordersGroup := router.Group("orders")
	usersGroup := router.Group("users")

	order.HandleRoutes(ordersGroup, *orderPresenter)
	user.HandleRoutes(usersGroup, *userPresenter)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverConfig.Port),
		Handler: router,
	}

	log.Println("Starting server...")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")

	contextWithTimout, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()

	if err = client.Disconnect(contextWithTimout); err != nil {
		log.Fatal(err)
	}

	db.Close()

	if err := server.Shutdown(context.TODO()); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
}
