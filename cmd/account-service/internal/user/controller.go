package user

import (
	"errors"
	"math"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/order_repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/user_repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
)

type OrderRepository interface {
	GetAllOrdersByUserId(userId int64) ([]order_repository.Order, error)
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

type Controller struct {
	orderRepo  OrderRepository
	userRepo   UserRepository
	cryptoRepo CryptoPricesRepository
	stockRepo  StockPricesRepository
}

func NewController(orderRepo OrderRepository, userRepo UserRepository, cryptoRepo CryptoPricesRepository, stockRepo StockPricesRepository) *Controller {
	return &Controller{
		orderRepo:  orderRepo,
		userRepo:   userRepo,
		cryptoRepo: cryptoRepo,
		stockRepo:  stockRepo,
	}
}

func (controller *Controller) CreateUser(userId int64, name string) error {
	_, err := controller.userRepo.GetUserByUserId(userId)

	if err == nil {
		return errors.New("user with this id already exists")
	}

	return controller.userRepo.CreateUser(userId, name)
}

func (controller *Controller) CreateBot(creatorId int64, amount float64) error {
	user, err := controller.userRepo.GetUserById(creatorId)
	if err != nil {
		return err
	}

	if user.Amount < amount {
		return errors.New("could not create bot! Insufficient amount")
	}

	if user.IsBot {
		return errors.New("bots can't have their own bots")
	}

	if err = controller.userRepo.UpdateUserAmount(creatorId, math.Round((user.Amount-amount)*100)/100); err != nil {
		return err
	}

	return controller.userRepo.CreateBot(creatorId, amount)
}

func (controller *Controller) MergeUserAndBot(botId int64) error {
	bot, botErr := controller.userRepo.GetUserById(botId)
	if botErr != nil {
		return botErr
	}

	if !bot.IsBot {
		return errors.New("can't merge 2 users - one must be bot")
	}

	user, userErr := controller.userRepo.GetUserById(bot.CreatorID.Int64)
	if userErr != nil {
		return userErr
	}

	if ordersErr := controller.orderRepo.UpdateOrdersCreatorByUserId(bot.ID, user.ID); ordersErr != nil {
		return ordersErr
	}

	if amountErr := controller.userRepo.UpdateUserAmount(user.ID, math.Round((user.Amount+bot.Amount)*100)/100); amountErr != nil {
		return amountErr
	}

	return controller.userRepo.DeleteUserById(bot.ID)
}

func (controller *Controller) EstimateUserAmount(userId int64) (float64, error) {
	orders, err := controller.orderRepo.GetAllOrdersByUserId(userId)
	if err != nil {
		return 0, err
	}

	quantityMap := make(map[string]float64)
	for _, order := range orders {
		if order.Type == "buy" {
			quantityMap[order.Symbol] += order.Amount
		} else {
			quantityMap[order.Symbol] -= order.Amount
		}
	}

	userAmount, err := controller.userRepo.GetUserAmount(userId)

	if err != nil {
		return 0, err
	}

	for symbol, amount := range quantityMap {
		price, err := controller.cryptoRepo.GetMostRecentPriceBySymbol(symbol)
		if err != nil {
			price, err := controller.stockRepo.GetMostRecentPriceBySymbol(symbol)
			if err != nil {
				return 0, err
			}
			userAmount += amount * price.BidPrice
		} else {
			userAmount += amount * price.BidPrice
		}
	}

	return userAmount, nil
}
