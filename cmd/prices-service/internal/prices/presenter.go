package prices

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/stream"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/prices_repository"
	"github.com/asaskevich/EventBus"
	"github.com/gorilla/websocket"
)

type PricesRepository interface {
	StoreEntry(prices prices_repository.Prices) error
}

type Presenter struct {
	cryptoRepo   *prices_repository.Collection
	stocksRepo   *prices_repository.Collection
	bus          EventBus.Bus
	cryptoStream *stream.Stream
	stockStream  *stream.Stream
}

func (self *Presenter) Handler(w http.ResponseWriter, r *http.Request, topic string) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	go self.writer(conn, topic)
}

func (self *Presenter) writer(conn *websocket.Conn, topic string) {
	self.bus.Subscribe(topic, func(json interface{}) {
		conn.WriteJSON(json)
	})
}

func (self *Presenter) BusHandler(b []byte, topic string) {
	if topic == "crypto" {
		var resp []stream.CryptoResponse
		if err := json.Unmarshal(b, &resp); err != nil {
			fmt.Println(err)
		}
		self.bus.Publish(topic, resp[len(resp)-1])
	}
	if topic == "stocks" {
		var resp []stream.StockResponse
		if err := json.Unmarshal(b, &resp); err != nil {
			fmt.Println(err)
		}
		self.bus.Publish(topic, resp[len(resp)-1])
	}
}

func (self *Presenter) StartStream(wsConfig env.WebSocetConfig) {
	cryptoStreamConfig := stream.NewStreamConfig(wsConfig)
	stockStreamConfig := stream.NewStockConfig(wsConfig)

	cryptoStream, err := stream.NewPriceStream(cryptoStreamConfig)
	if err != nil {
		panic(err)
	}
	self.cryptoStream = cryptoStream

	stockStream, err := stream.NewPriceStream(stockStreamConfig)
	if err != nil {
		panic(err)
	}
	self.stockStream = stockStream

	go func() {
		if err := self.cryptoStream.Start(func(b []byte) { self.BusHandler(b, "crypto") }); err != nil {
			panic(err)
		}
	}()
	go func() {
		if err := self.stockStream.Start(func(b []byte) { self.BusHandler(b, "stocks") }); err != nil {
			panic(err)
		}
	}()

}
func (self *Presenter) StopStream() {
	self.cryptoStream.Stop()
	self.stockStream.Stop()
}

// TODO: PriceRepository -> CryptoPriceRepo & StockPriceRepo
func NewPresenter(cryptoRepo prices_repository.Collection, stocksRepo prices_repository.Collection, bus EventBus.Bus) Presenter {
	result := Presenter{
		cryptoRepo: &cryptoRepo,
		stocksRepo: &stocksRepo,
		bus:        bus,
	}

	bus.Subscribe("crypto", func(resp stream.CryptoResponse) {
		temp := prices_repository.NewCryptoPrice(resp.Symbol, resp.Exchange, resp.BidPrice, resp.AskPrice, resp.Date)

		err := result.cryptoRepo.StoreCryptoPrice(temp)
		if err != nil {
			log.Fatal(err)
		}
	})

	bus.Subscribe("stocks", func(resp stream.StockResponse) {
		price := prices_repository.NewPrice(resp.Symbol, resp.Type, resp.BidPrice, resp.AskPrice, resp.Date)
		stockPrice := prices_repository.NewStockPrice(price, resp.AskExchange, resp.BidExchange, resp.TradeSize, resp.Conditions, resp.Type)

		err := result.stocksRepo.StoreStockPrice(stockPrice)
		if err != nil {
			log.Fatal(err)
		}
	})

	return result
}
