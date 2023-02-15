package main

import (
	"fmt"
	"net/http"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/order"
	repository "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/order_repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/user_repository"
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

	orderRepository := order_repository.NewOrderTable(db)
	userRepository := user_repository.NewUserTable(db)

	mux := http.NewServeMux()

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	orderPresenter := order.NewOrderPresenter(orderRepository, userRepository, &upgrader)
	orderHandler := order.NewOrderHandler(orderPresenter)

	order.HandleRoutes(mux, orderHandler)

	fmt.Println("Starting server on", serverConfig.Port)
	serverErr := http.ListenAndServe(fmt.Sprintf("localhost:%d", serverConfig.Port), mux)
	if serverErr != nil {
		fmt.Println("Failed to start server:", err)
	}
}
