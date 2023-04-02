package prices

import (
	"fmt"
	"log"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/stream"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
)

type PricesRepository[P database.CryptoPrices | database.StockPrices] interface {
	StoreEntry(prices P) error
}

type RepositoryController struct {
	cryptoRepo PricesRepository[database.CryptoPrices]
	stocksRepo PricesRepository[database.StockPrices]
}

func NewRepositoryController(cryptoRepo PricesRepository[database.CryptoPrices], stocksRepo PricesRepository[database.StockPrices]) *RepositoryController {
	return &RepositoryController{
		cryptoRepo: cryptoRepo,
		stocksRepo: stocksRepo,
	}
}

func (r *RepositoryController) ListenForStoring(bus EventBus) error {
	if err := bus.Subscribe("crypto", r.handleCryptoMessage); err != nil {
		return fmt.Errorf("error with subscribing to crypto messages: %w", err)
	}

	if err := bus.Subscribe("stocks", r.handleStocksMessage); err != nil {
		return fmt.Errorf("error with subscribing to stocks messages: %w", err)
	}

	return nil
}

func (r *RepositoryController) handleCryptoMessage(resp stream.CryptoResponse) {
	cryptoPrices := database.CryptoPrices{
		Prices: database.Prices{
			Symbol:   resp.Symbol,
			BidPrice: resp.BidPrice,
			BidSize:  resp.BidSize,
			AskPrice: resp.AskPrice,
			AskSize:  resp.AskSize,
			Date:     resp.Date,
		},
		Exchange: resp.Exchange,
	}

	if err := r.cryptoRepo.StoreEntry(cryptoPrices); err != nil {
		log.Println(err)
	}
}

func (r *RepositoryController) handleStocksMessage(resp stream.StockResponse) {
	stockPrice := database.StockPrices{
		Prices: database.Prices{
			Symbol:   resp.Symbol,
			BidPrice: resp.BidPrice,
			BidSize:  resp.BidSize,
			AskPrice: resp.AskPrice,
			AskSize:  resp.AskSize,
			Date:     resp.Date,
		},
		AskExchange: resp.AskExchange,
		BidExchange: resp.BidExchange,
		TradeSize:   resp.TradeSize,
		Conditions:  resp.Conditions,
		Tape:        resp.Tape,
	}

	if err := r.stocksRepo.StoreEntry(stockPrice); err != nil {
		log.Println(err)
	}
}
