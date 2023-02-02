package main

import (
	"fmt"
	"net/http"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/order"
	repository "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/order_repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/user_repository"
)

func main() {
	dbConfig := env.LoadDBConfig()
	db, err := repository.Connect(dbConfig)

	defer db.Close()

	if err != nil {
		panic(err)
	}

	serverConfig := env.LoadServerConfig()

	orderRepository := order_repository.NewOrderTable(db)
	userRepository := user_repository.NewUserTable(db)

	orderPresenter := order.NewOrderPresenter(orderRepository, userRepository)
	orderHandler := order.NewOrderHandler(orderPresenter)

	// router := gin.Default()

	// order.SetupRoutes(router, orderController)

	// router.Run()

	mux := http.NewServeMux()

	order.HandleRoutes(mux, orderHandler)

	serverErr := http.ListenAndServe(fmt.Sprintf("localhost:%d", serverConfig.Port), mux)

	if serverErr != nil {
		panic(serverErr)
	}
}
