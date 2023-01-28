package prices

import (
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/repositories/prices_repository"
	"time"
)

type PricesRepository interface {
	StoreEntry(prices prices_repository.Prices) error
	GetAllPrices() ([]prices_repository.Prices, error)
	GetAllPricesBySymbol(symbol string) ([]prices_repository.Prices, error)
	GetAllPricesByExchange(exchange string) ([]prices_repository.Prices, error)
	GetAllPricesInPeriod(from time.Time, to time.Time) ([]prices_repository.Prices, error)
	GetAllPricesInPeriodSymbol(from time.Time, to time.Time, symbol string) ([]prices_repository.Prices, error)
}

type Presenter struct {
	pricesRepository PricesRepository
}

func NewPresenter(repository PricesRepository) Presenter {
	return Presenter{
		pricesRepository: repository,
	}
}
