package prices

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/stream"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
	"github.com/asaskevich/EventBus"
	"github.com/gorilla/websocket"
)

type PricesRepository[Prices database.CryptoPrices | database.StockPrices] interface {
	StoreEntry(prices Prices) error
}

type Presenter struct {
	cryptoRepo   PricesRepository[database.CryptoPrices]
	stocksRepo   PricesRepository[database.StockPrices]
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
func (self *Presenter) StockHandler(w http.ResponseWriter, r *http.Request){
	self.Handler(w,r,"stocks")
}
func (self *Presenter) CryptoHandler(w http.ResponseWriter, r *http.Request){
	self.Handler(w,r,"crypto")
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


func NewPresenter(cryptoRepo   PricesRepository[database.CryptoPrices], stocksRepo   PricesRepository[database.StockPrices], bus EventBus.Bus) Presenter {
	result := Presenter{
		cryptoRepo: cryptoRepo,
		stocksRepo: stocksRepo,
		bus:        bus,
	}

	bus.Subscribe("crypto", func(resp stream.CryptoResponse) {
		temp := database.CryptoPrices{
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

		err := result.cryptoRepo.StoreEntry(temp)
		if err != nil {
			log.Fatal(err)
		}
	})

	bus.Subscribe("stocks", func(resp stream.StockResponse) {

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
			TradeSize: resp.TradeSize,
			Conditions:resp.Conditions,
			Tape: resp.Tape,

		}

		
		err := result.stocksRepo.StoreEntry(stockPrice)
		if err != nil {
			log.Fatal(err)
		}
	})

	return result
}
