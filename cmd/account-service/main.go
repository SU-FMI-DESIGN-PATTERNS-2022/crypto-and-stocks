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

	order.NewPresenter(orderRepository, userRepository)
}
