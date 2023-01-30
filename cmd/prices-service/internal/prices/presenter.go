package prices

import (
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/prices_repository"
)

type PricesRepository interface {
	StoreEntry(prices prices_repository.Prices) error
}

type Presenter struct {
	pricesRepository PricesRepository
}

// TODO: PriceRepository -> CryptoPriceRepo & StockPriceRepo
func NewPresenter(repository PricesRepository) Presenter {
	return Presenter{
		pricesRepository: repository,
	}
}
