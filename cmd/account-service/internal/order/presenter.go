package order

import (
	"errors"

	"net/http"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/order_repository"
	"github.com/gorilla/websocket"
)

type OrderRepository interface {
	StoreOrder(order order_repository.Order) error
	GetAllOrders() ([]order_repository.Order, error)
	GetAllOrdersByUserId(userId int64) ([]order_repository.Order, error)
	GetAllOrdersBySymbol(symbol string) ([]order_repository.Order, error)
	GetAllOrdersByUserIdAndSymbol(userId int64, symbol string) ([]order_repository.Order, error)
}

type UserRepository interface {
	GetUserAmount(id int64) (float64, error)
	UpdateUserAmount(id int64, amount float64) error
}

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error)
}

type OrderPresenter struct {
	orderRepo OrderRepository
	userRepo  UserRepository
	upgrader  Upgrader
}

func NewOrderPresenter(orderRepo OrderRepository, userRepo UserRepository, upgrader Upgrader) OrderPresenter {
	return OrderPresenter{
		orderRepo: orderRepo,
		userRepo:  userRepo,
		upgrader:  upgrader,
	}
}

func (orderPresenter *OrderPresenter) StoreOrder(order order_repository.Order) error {
	switch order.Type {
	case "buy":
		amount, err := orderPresenter.userRepo.GetUserAmount(order.UserID)
		if err != nil {
			return err
		}

		if amount < order.Amount*order.Price {
			return errors.New("not enough amount")
		}

		updateErr := orderPresenter.userRepo.UpdateUserAmount(order.UserID, amount-order.Amount*order.Price)
		if updateErr != nil {
			return updateErr
		}
	case "sell":
		orders, err := orderPresenter.orderRepo.GetAllOrdersByUserIdAndSymbol(order.UserID, order.Symbol)
		if err != nil {
			return err
		}

		var amount float64
		for _, o := range orders {
			if o.Type == "buy" {
				amount += o.Amount
			} else {
				amount -= o.Amount
			}
		}

		if amount < order.Amount {
			return errors.New("not enough amount")
		}

		updateErr := orderPresenter.userRepo.UpdateUserAmount(order.UserID, amount+order.Amount*order.Price)
		if updateErr != nil {
			return updateErr
		}
	default:
		return errors.New("invalid order type")
	}

	return orderPresenter.orderRepo.StoreOrder(order)
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
