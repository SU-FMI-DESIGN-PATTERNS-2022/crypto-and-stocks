package order

import (
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/order_repository"
)

type OrderRepository interface {
	StoreOrder(order order_repository.Order) error
	GetAllOrders() ([]order_repository.Order, error)
	GetAllOrdersByUserId(userId int64) ([]order_repository.Order, error)
	GetAllOrdersBySymbol(symbol string) ([]order_repository.Order, error)
}

type UserRepository interface {
	CreateUser(userId int64, name string) error
	CreateBot(creatorID int64, amount float64) error
	AddOrder(userId int64, orderId int64) error
	MergeUserOrders(id int64) error
}

type Presenter struct {
	orderRepo OrderRepository
	userRepo  UserRepository
}

func NewPresenter(orderRepo OrderRepository, userRepo UserRepository) Presenter {
	return Presenter{
		orderRepo: orderRepo,
		userRepo:  userRepo,
	}
}
