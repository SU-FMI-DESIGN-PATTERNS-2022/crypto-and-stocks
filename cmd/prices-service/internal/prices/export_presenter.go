package prices

import (
	"time"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
)

type ExportPricesRepository[Prices database.CryptoPrices | database.StockPrices] interface {
	GetAllPrices() ([]Prices, error)
	GetAllPricesBySymbol(symbol string) ([]Prices, error)
	GetAllPricesByExchange(exchange string) ([]Prices, error)
	GetAllPricesInPeriod(from time.Time, to time.Time) ([]Prices, error)
	GetAllPricesInPeriodSymbol(from time.Time, to time.Time, symbol string) ([]Prices, error)
	GetMostRecentPriceBySymbol(symbol string) (Prices, error)
}

type ExportStockPriceRepository interface {
	ExportPricesRepository[database.StockPrices]
}

type ExportCryptoPriceRepository interface {
	ExportPricesRepository[database.CryptoPrices]
}

type ExportPricesPresenter struct {
	cryptoPricesRepo ExportCryptoPriceRepository
	stockPricesRepo  ExportStockPriceRepository
}

func NewExportPricesPresenter(cryptoPricesRepository ExportCryptoPriceRepository, stockPricesRepository ExportStockPriceRepository) ExportPricesPresenter {
	return ExportPricesPresenter{
		cryptoPricesRepo: cryptoPricesRepository,
		stockPricesRepo:  stockPricesRepository,
	}
}
