package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/prices_repository"
	"github.com/asaskevich/EventBus"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/stream"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func handler(w http.ResponseWriter, r *http.Request, bus EventBus.Bus) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	go writerCrypto(conn, bus)
}

func writerCrypto(conn *websocket.Conn, bus EventBus.Bus) {
	bus.Subscribe("crypto", func(json interface{}) {
		conn.WriteJSON(json)
	})
}

func cryptoHandler(b []byte) {
	var cryptoResponse []stream.CryptoResponse
	if err := json.Unmarshal(b, &cryptoResponse); err != nil {
		fmt.Println(err)
	}

	fmt.Println(cryptoResponse)
}

func cryptoHandlerBus(b []byte, bus EventBus.Bus) {
	var cryptoResponse []stream.CryptoResponse
	if err := json.Unmarshal(b, &cryptoResponse); err != nil {
		fmt.Println(err)
	}
	bus.Publish("crypto", cryptoResponse[len(cryptoResponse)-1])
}

func stockHandler(b []byte) {
	var stockResponse []stream.StockResponse
	if err := json.Unmarshal(b, &stockResponse); err != nil {
		fmt.Println(err)
	}

	fmt.Println(stockResponse)
}

func main() {
	bus := EventBus.New()

	mongoConfig := env.LoadMongoConfig()
	ctx := context.TODO()
	client, cancel, connectErr := mongo.Connect(mongoConfig)

	var pricesRepo = prices_repository.NewCollection(client, "crypto", "prices")

	//var pricesPresenter = prices.NewPresenter(pricesRepo)

	if connectErr != nil {
		panic(connectErr)
	}
	defer func() {
		if connectErr = client.Disconnect(ctx); connectErr != nil {
			panic(connectErr)
		}
	}()

	defer mongo.Close(client, ctx, cancel)
	// Checking whether the connection was successful
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged MongoDB.")

	wsConfig := env.LoadWebSocetConfig()
	cryptoStreamConfig := stream.StreamConfig{
		URL:    wsConfig.CryptoURL,
		Quotes: wsConfig.CryptoQuotes,
		Key:    wsConfig.Key,
		Secret: wsConfig.Secret,
	}

	stockStreamConfig := stream.StreamConfig{
		URL:    wsConfig.StockURL,
		Quotes: wsConfig.StockQuotes,
		Key:    wsConfig.Key,
		Secret: wsConfig.Secret,
	}

	cryptoStream, err := stream.NewPriceStream(cryptoStreamConfig)
	if err != nil {
		panic(err)
	}
	stockStream, err := stream.NewPriceStream(stockStreamConfig)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := cryptoStream.Start(func(b []byte) { cryptoHandlerBus(b, bus) }); err != nil {
			panic(err)
		}
	}()

	bus.Subscribe("crypto", func(cryptoResponse stream.CryptoResponse) {

		id := primitive.NewObjectID()
		temp := prices_repository.Prices{
			ID:       id,
			Symbol:   cryptoResponse.Symbol,
			Exchange: cryptoResponse.Exchange,
			BidPrice: cryptoResponse.BidPrice,
			AskPrice: cryptoResponse.AskPrice,
			Date:     cryptoResponse.Date,
		}

		err := pricesRepo.StoreEntry(temp)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) { handler(rw, r, bus) })
	go log.Fatal(http.ListenAndServe(*addr, nil))

	go func() {
		if err := stockStream.Start(stockHandler); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Minute)
	cryptoStream.Stop()
	stockStream.Stop()
	print("ALL PRICES ________________________________________________________________________________________________________")
}
