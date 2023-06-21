package main

import (
	"context"
	"fmt"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/order"
	repository "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/order_repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/user_repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/user"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
	mongo_env "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"
	"github.com/gin-gonic/gin"
)

func main() {

	dbConfig := env.LoadDBConfig()
	db, err := repository.Connect(dbConfig)

	if err != nil {
		fmt.Println("Failed to open database:", err)
		return
	}

	defer db.Close()

	serverConfig := env.LoadServerConfig()

	mongoConfig := mongo_env.LoadMongoConfig()
	client, err := database.Connect(mongoConfig, database.Remote)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	orderRepository := order_repository.NewOrderTable(db)
	userRepository := user_repository.NewUserTable(db)
	cryptoRepository := database.NewCollection[database.CryptoPrices](client, mongoConfig.Database, "CryptoPrices")
	stockRepository := database.NewCollection[database.StockPrices](client, mongoConfig.Database, "StockPrices")

	orderPresenter := order.NewOrderPresenter(orderRepository, userRepository)
	userPresenter := user.NewUserPresenter(orderRepository, userRepository, cryptoRepository, stockRepository)

	router := gin.Default()

	ordersGroup := router.Group("orders")
	usersGroup := router.Group("users")

	order.HandleRoutes(ordersGroup, *orderPresenter)
	user.HandleRoutes(usersGroup, *userPresenter)

	router.Run("localhost:" + serverConfig.Port)
}
