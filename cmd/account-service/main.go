package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/order"
	repository "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/order_repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/user_repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/user"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
	mongo_env "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	orderPresenter := order.NewOrderPresenter(orderRepository, userRepository, &upgrader)
	userPresenter := user.NewUserPresenter(orderRepository, userRepository, cryptoRepository, stockRepository)

	router := gin.Default()

	ordersGroup := router.Group("orders")
	ordersGroup.GET("/all", func(context *gin.Context) {
		orders, err := orderPresenter.GetAllOrders()
		if err != nil {
			context.JSON(http.StatusInternalServerError, "Could not fetch orders")
			return
		}

		context.JSON(http.StatusOK, orders)
	})
	ordersGroup.GET("/store", func(context *gin.Context) {
		conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)

		if err != nil {
			return
		}

		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Something is wrong."))
				break
			}

			var order order_repository.Order
			if err = json.Unmarshal(message, &order); err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("The message is not in right json object structure."))
				break
			}

			if err = orderPresenter.StoreOrder(order); err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("We have a problem with storing your order."))
				break
			}
		}
	})

	usersGroup := router.Group("users")
	usersGroup.POST("/create/user/:id/:name", func(context *gin.Context) {
		name := context.Param("name")
		id, err := strconv.ParseInt(context.Param("id"), 10, 64)
		if err != nil {
			context.JSON(http.StatusBadRequest, "Parameter id should be number")
			return
		}

		if err = userPresenter.CreateUser(id, name); err != nil {
			context.JSON(http.StatusBadRequest, err.Error())
			return
		}

		context.JSON(http.StatusCreated, "User created successfully")
	})

	router.Run("localhost:" + serverConfig.Port)
}
