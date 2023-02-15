package order

import (
	"errors"
	"math"

	"net/http"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/order_repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/user_repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
	"github.com/gorilla/websocket"
)

type OrderRepository interface {
	StoreOrder(order order_repository.Order) error
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
	GetUserByUserId(id int64) (user_repository.User, error)
	GetUserAmount(id int64) (float64, error)
	UpdateUserAmount(id int64, amount float64) error
	DeleteUserById(id int64) error
}

type PricesRepository[Prices database.CryptoPrices | database.StockPrices] interface {
	GetMostRecentPriceBySymbol(symbol string) (Prices, error)
}

type CryptoPricesRepository interface {
	PricesRepository[database.CryptoPrices]
}

type StockPricesRepository interface {
	PricesRepository[database.StockPrices]
}

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error)
}

type OrderPresenter struct {
	orderRepo  OrderRepository
	userRepo   UserRepository
	cryptoRepo CryptoPricesRepository
	stockRepo  StockPricesRepository
	upgrader   Upgrader
}

func NewOrderPresenter(orderRepo OrderRepository, userRepo UserRepository, cryptoRepo CryptoPricesRepository, stockRepo StockPricesRepository, upgrader Upgrader) OrderPresenter {
	return OrderPresenter{
		orderRepo:  orderRepo,
		userRepo:   userRepo,
		cryptoRepo: cryptoRepo,
		stockRepo:  stockRepo,
		upgrader:   upgrader,
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

func (orderPresenter *OrderPresenter) CreateUser(userId int64, name string) error {
	_, err := orderPresenter.userRepo.GetUserByUserId(userId)

	if err == nil {
		return errors.New("user with this id already exists")
	}

	return orderPresenter.userRepo.CreateUser(userId, name)
}

func (orderPresenter *OrderPresenter) CreateBot(creatorId int64, amount float64) error {
	user, userErr := orderPresenter.userRepo.GetUserById(creatorId)
	if userErr != nil {
		return userErr
	}

	if user.Amount < amount {
		return errors.New("could not create bot! Insufficient amount")
	}

	if user.IsBot {
		return errors.New("bots can't have their own bots")
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
		return errors.New("can't merge 2 users - one must be bot")
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

func (orderPresenter *OrderPresenter) EstimateUserAmount(userId int64) (float64, error) {
	orders, err := orderPresenter.orderRepo.GetAllOrdersByUserId(userId)
	if err != nil {
		return 0, err
	}

	quantityMap := make(map[string]float64)
	for _, o := range orders {
		if o.Type == "buy" {
			quantityMap[o.Symbol] += o.Amount
		} else {
			quantityMap[o.Symbol] -= o.Amount
		}
	}

	amount, err := orderPresenter.userRepo.GetUserAmount(userId)

	if err != nil {
		return 0, err
	}

	for s, a := range quantityMap {
		cp, err := orderPresenter.cryptoRepo.GetMostRecentPriceBySymbol(s)
		if err != nil {
			sp, err := orderPresenter.stockRepo.GetMostRecentPriceBySymbol(s)
			if err != nil {
				return 0, err
			}
			amount += a * sp.BidPrice
		} else {
			amount += a * cp.BidPrice
		}
	}

	return amount, nil
}
