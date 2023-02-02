package prices

import "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/repositories/crypto_prices_repository"

type StockPriceRepository interface {
	//TODO: Add needed methods
}
type CryptoPriceRepository interface {
	StoreEntry(price crypto_prices_repository.CryptoPrices) error
}

type PricesPresenter struct {
	cryptoPricesRepo CryptoPriceRepository
	stockPricesRepo  StockPriceRepository
}

// TODO: PriceRepository -> CryptoPriceRepo & StockPriceRepo
func NewPricesPresenter(cryptoPricesRepository CryptoPriceRepository, stockPricesRepository StockPriceRepository) PricesPresenter {
	return PricesPresenter{
		cryptoPricesRepo: cryptoPricesRepository,
		stockPricesRepo:  stockPricesRepository,
	}
}
