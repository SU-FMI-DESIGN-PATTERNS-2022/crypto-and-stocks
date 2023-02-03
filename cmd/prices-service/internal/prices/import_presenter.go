package prices

import (
	"encoding/json"
	"fmt"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/stream"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
)

type ImportPricesRepository[Prices database.CryptoPrices | database.StockPrices] interface {
	StoreEntry(price Prices) error
}

type ImportStockPriceRepository interface {
	ImportPricesRepository[database.StockPrices]
}

type ImportCryptoPriceRepository interface {
	ImportPricesRepository[database.CryptoPrices]
}

type ImportPricesPresenter struct {
	cryptoPricesRepo ImportCryptoPriceRepository
	stockPricesRepo  ImportStockPriceRepository
	cryptoStream     *stream.Stream
	stockStream      *stream.Stream
}

func NewImportPricesPresenter(cryptoPricesRepository ImportCryptoPriceRepository, stockPricesRepository ImportStockPriceRepository, cryptoStream *stream.Stream, stockStream *stream.Stream) ImportPricesPresenter {
	return ImportPricesPresenter{
		cryptoPricesRepo: cryptoPricesRepository,
		stockPricesRepo:  stockPricesRepository,
		cryptoStream:     cryptoStream,
		stockStream:      stockStream,
	}
}

func (presenter *ImportPricesPresenter) cryptoHandler(b []byte) {
	var cryptoResponse []stream.CryptoResponse
	if err := json.Unmarshal(b, &cryptoResponse); err != nil {
		fmt.Println(err)
	}
	//TODO: Make a method that saves the response into the corresponding collection
	fmt.Println(cryptoResponse)
}

func (presenter ImportPricesPresenter) stockHandler(b []byte) {
	var stockResponse []stream.StockResponse
	if err := json.Unmarshal(b, &stockResponse); err != nil {
		fmt.Println(err)
	}
	//TODO: same here
	fmt.Println(stockResponse)
}
