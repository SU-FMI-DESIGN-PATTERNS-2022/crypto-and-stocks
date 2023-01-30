package main

import (
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

	orderRepository := order_repository.NewOrderTable(db)
	userRepository := user_repository.NewUserTable(db)

	// userRepository.CreateBot(1, 59.99)
	// userRepository.CreateBot(1, 20.99)
	// userRepository.MergeAllUserOrders(1)

	order.NewPresenter(orderRepository, userRepository)
}
