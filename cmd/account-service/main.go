package main

import (
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/order"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repository/orderDB"
)

func main() {
	dbConfig := env.LoadDBConfig()
	db, err := repository.Connect(dbConfig)
	if err != nil {
		panic(err)
	}

	orderDB := orderDB.NewDatabase(db)

	order.NewPresenter(orderDB)

	orderDB.Close()
}
