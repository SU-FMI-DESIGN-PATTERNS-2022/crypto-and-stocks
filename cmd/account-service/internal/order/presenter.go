package order

import (
	"errors"
	"math"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/order_repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/user_repository"
)

type OrderRepository interface {
	StoreOrder(userId int64, orderType string, symbol string, amount float64, price float64) error
	GetAllOrders() ([]order_repository.Order, error)
	GetAllOrdersByUserId(userId int64) ([]order_repository.Order, error)
	GetAllOrdersBySymbol(symbol string) ([]order_repository.Order, error)
}

type UserRepository interface {
	CreateUser(userId int64, name string) error
	CreateBot(creatorId int64, amount float64) error
	GetUserById(id int64) (user_repository.User, error)
	UpdateUserAmount(id int64, amount float64) error
	MergeUserAndBot(id int64) error
	MergeAllUserOrders(id int64) error
	GetAmountByUserId(id int64) (float64, error)
}

type OrderPresenter struct {
	orderRepo OrderRepository
	userRepo  UserRepository
}

func NewOrderPresenter(orderRepo OrderRepository, userRepo UserRepository) OrderPresenter {
	return OrderPresenter{
		orderRepo: orderRepo,
		userRepo:  userRepo,
	}
}

func (orderPresenter *OrderPresenter) CreateBot(creatorId int64, amount float64) error {
	user, userErr := orderPresenter.userRepo.GetUserById(creatorId)
	if userErr != nil {
		return userErr
	}

	if user.Amount < amount {
		return errors.New("Could not create bot! Insufficient amount!")
	}

	if user.IsBot {
		return errors.New("Bots can't have their own bots!")
	}

	err := orderPresenter.userRepo.UpdateUserAmount(creatorId, math.Round((user.Amount-amount)*100)/100)

	if err != nil {
		return err
	}

	return orderPresenter.userRepo.CreateBot(creatorId, amount)
}

func (orderPresenter *OrderPresenter) GetAllOrders() ([]order_repository.Order, error) {
	orders, err := orderPresenter.orderRepo.GetAllOrders()

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (orderPresenter *OrderPresenter) GetAllOrdersByUserId(userId int64) ([]order_repository.Order, error) {
	orders, err := orderPresenter.orderRepo.GetAllOrdersByUserId(userId)

	if err != nil {
		return nil, err
	}

	return orders, nil
}
