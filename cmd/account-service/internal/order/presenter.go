package order

import "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repository/orderDB"

type Repository interface { // bridge - structural designa pattern
	StoreOrder(order orderDB.Model) error
	GetAllOrders() ([]orderDB.Model, error)
}

type Presenter struct {
	repo Repository
}

func NewPresenter(repo Repository) Presenter {
	return Presenter{
		repo: repo,
	}
}
