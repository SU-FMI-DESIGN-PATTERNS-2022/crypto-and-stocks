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

type UserPresenter struct {
	orderRepo  OrderRepository
	userRepo   UserRepository
	cryptoRepo CryptoPricesRepository
	stockRepo  StockPricesRepository
}

func NewUserPresenter(orderRepo OrderRepository, userRepo UserRepository, cryptoRepo CryptoPricesRepository, stockRepo StockPricesRepository) *UserPresenter {
	return &UserPresenter{
		orderRepo:  orderRepo,
		userRepo:   userRepo,
		cryptoRepo: cryptoRepo,
		stockRepo:  stockRepo,
	}
}

func (userPresenter *UserPresenter) CreateUser(userId int64, name string) error {
	_, err := userPresenter.userRepo.GetUserByUserId(userId)

	if err == nil {
		return errors.New("user with this id already exists")
	}

	return userPresenter.userRepo.CreateUser(userId, name)
}

func (userPresenter *UserPresenter) CreateBot(creatorId int64, amount float64) error {
	user, err := userPresenter.userRepo.GetUserById(creatorId)
	if err != nil {
		return err
	}

	if user.Amount < amount {
		return errors.New("could not create bot! Insufficient amount")
	}

	if user.IsBot {
		return errors.New("bots can't have their own bots")
	}

	if err = userPresenter.userRepo.UpdateUserAmount(creatorId, math.Round((user.Amount-amount)*100)/100); err != nil {
		return err
	}

	return userPresenter.userRepo.CreateBot(creatorId, amount)
}

func (userPresenter *UserPresenter) MergeUserAndBot(botId int64) error {
	bot, botErr := userPresenter.userRepo.GetUserById(botId)
	if botErr != nil {
		return botErr
	}

	if !bot.IsBot {
		return errors.New("can't merge 2 users - one must be bot")
	}

	user, userErr := userPresenter.userRepo.GetUserById(bot.CreatorID.Int64)
	if userErr != nil {
		return userErr
	}

	if ordersErr := userPresenter.orderRepo.UpdateOrdersCreatorByUserId(bot.ID, user.ID); ordersErr != nil {
		return ordersErr
	}

	if amountErr := userPresenter.userRepo.UpdateUserAmount(user.ID, math.Round((user.Amount+bot.Amount)*100)/100); amountErr != nil {
		return amountErr
	}

	return userPresenter.userRepo.DeleteUserById(bot.ID)
}

func (userPresenter *UserPresenter) EstimateUserAmount(userId int64) (float64, error) {
	orders, err := userPresenter.orderRepo.GetAllOrdersByUserId(userId)
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

	userAmount, err := userPresenter.userRepo.GetUserAmount(userId)

	if err != nil {
		return 0, err
	}

	for symbol, amount := range quantityMap {
		price, err := userPresenter.cryptoRepo.GetMostRecentPriceBySymbol(symbol)
		if err != nil {
			price, err := userPresenter.stockRepo.GetMostRecentPriceBySymbol(symbol)
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
