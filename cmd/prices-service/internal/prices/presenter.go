package prices

import (
	"time"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/repositories/crypto_prices_repository"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/repositories/stock_prices_repository"
)

type StockPriceRepository interface {
	StoreEntry(price stock_prices_repository.StockPrices) error
	GetAllPrices() ([]stock_prices_repository.StockPrices, error)
	GetAllPricesBySymbol(symbol string) ([]stock_prices_repository.StockPrices, error)
	GetAllPricesByExchange(exchange string) ([]stock_prices_repository.StockPrices, error)
	GetAllPricesInPeriod(from time.Time, to time.Time) ([]stock_prices_repository.StockPrices, error)
	GetAllPricesInPeriodSymbol(from time.Time, to time.Time, symbol string) ([]stock_prices_repository.StockPrices, error)
	GetMostRecentPriceBySymbol(symbol string) (stock_prices_repository.StockPrices, error)
}

type CryptoPriceRepository interface {
	StoreEntry(price crypto_prices_repository.CryptoPrices) error
	GetAllPrices() ([]crypto_prices_repository.CryptoPrices, error)
	GetAllPricesBySymbol(symbol string) ([]crypto_prices_repository.CryptoPrices, error)
	GetAllPricesByExchange(exchange string) ([]crypto_prices_repository.CryptoPrices, error)
	GetAllPricesInPeriod(from time.Time, to time.Time) ([]crypto_prices_repository.CryptoPrices, error)
	GetAllPricesInPeriodSymbol(from time.Time, to time.Time, symbol string) ([]crypto_prices_repository.CryptoPrices, error)
	GetMostRecentPriceBySymbol(symbol string) (crypto_prices_repository.CryptoPrices, error)
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
