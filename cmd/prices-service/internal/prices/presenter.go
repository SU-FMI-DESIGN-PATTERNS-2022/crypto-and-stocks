package prices

import (
	"time"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/repositories/crypto_prices_repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/repositories/stock_prices_repository"
)

type PricesRepository[Prices crypto_prices_repository.CryptoPrices | stock_prices_repository.StockPrices] interface {
	StoreEntry(price Prices) error
	GetAllPrices() ([]Prices, error)
	GetAllPricesBySymbol(symbol string) ([]Prices, error)
	GetAllPricesByExchange(exchange string) ([]Prices, error)
	GetAllPricesInPeriod(from time.Time, to time.Time) ([]Prices, error)
	GetAllPricesInPeriodSymbol(from time.Time, to time.Time, symbol string) ([]Prices, error)
	GetMostRecentPriceBySymbol(symbol string) (Prices, error)
}

type StockPriceRepository interface {
	PricesRepository[crypto_prices_repository.CryptoPrices]
}

type CryptoPriceRepository interface {
	PricesRepository[stock_prices_repository.StockPrices]
}

type PricesPresenter struct {
	cryptoPricesRepo CryptoPriceRepository
	stockPricesRepo  StockPriceRepository
}

func NewPricesPresenter(cryptoPricesRepository CryptoPriceRepository, stockPricesRepository StockPriceRepository) PricesPresenter {
	return PricesPresenter{
		cryptoPricesRepo: cryptoPricesRepository,
		stockPricesRepo:  stockPricesRepository,
	}
}
