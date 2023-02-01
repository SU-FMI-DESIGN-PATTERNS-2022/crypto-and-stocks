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
	GetAllOrdersByUserIdAndSymbol(userId int64, symbol string) ([]order_repository.Order, error)
	UpdateOrdersCreatorByUserId(prevUserId int64, newUserId int64) error
}

type UserRepository interface {
	CreateUser(userId int64, name string) error
	CreateBot(creatorId int64, amount float64) error
	GetUserById(id int64) (user_repository.User, error)
	UpdateUserAmount(id int64, amount float64) error
	DeleteUserById(id int64) error
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

func (orderPresenter *OrderPresenter) CreateUser(userId int64, name string) error {
	_, err := orderPresenter.userRepo.GetUserById(userId)

	if err == nil {
		return errors.New("User with this id already exists")
	}

	return orderPresenter.userRepo.CreateUser(userId, name)
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

func (orderPresenter *OrderPresenter) GetAllOrdersBySymbol(symbol string) ([]order_repository.Order, error) {
	orders, err := orderPresenter.orderRepo.GetAllOrdersBySymbol(symbol)

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (orderPresenter *OrderPresenter) GetAllOrdersByUserIdAndSymbol(userId int64, symbol string) ([]order_repository.Order, error) {
	orders, err := orderPresenter.orderRepo.GetAllOrdersByUserIdAndSymbol(userId, symbol)

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (orderPresenter *OrderPresenter) MergeUserAndBot(botId int64) error {
	bot, botErr := orderPresenter.userRepo.GetUserById(botId)
	if botErr != nil {
		return botErr
	}

	if !bot.IsBot {
		return errors.New("Can't merge 2 users - one must be bot")
	}

	user, userErr := orderPresenter.userRepo.GetUserById(bot.CreatorID.Int64)
	if userErr != nil {
		return userErr
	}

	ordersErr := orderPresenter.orderRepo.UpdateOrdersCreatorByUserId(bot.ID, user.ID)
	if ordersErr != nil {
		return ordersErr
	}

	amountErr := orderPresenter.userRepo.UpdateUserAmount(user.ID, math.Round((user.Amount+bot.Amount)*100)/100)
	if amountErr != nil {
		return amountErr
	}

	return orderPresenter.userRepo.DeleteUserById(bot.ID)
}
