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

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/prices"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/stream"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

// func handler(w http.ResponseWriter, r *http.Request, bus EventBus.Bus, topic string) {
// 	var upgrader = websocket.Upgrader{
// 		ReadBufferSize:  1024,
// 		WriteBufferSize: 1024,
// 	}
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	go writer(conn, bus, topic)
// }

// func writer(conn *websocket.Conn, bus EventBus.Bus, topic string) {
// 	bus.Subscribe(topic, func(json interface{}) {
// 		conn.WriteJSON(json)
// 	})
// }

func cryptoHandler(b []byte) {
	var cryptoResponse []stream.CryptoResponse
	if err := json.Unmarshal(b, &cryptoResponse); err != nil {
		fmt.Println(err)
	}

	fmt.Println(cryptoResponse)
}

func BusHandler(b []byte, bus EventBus.Bus, topic string) {
	if topic == "crypto" {
		var resp []stream.CryptoResponse
		if err := json.Unmarshal(b, &resp); err != nil {
			fmt.Println(err)
		}
		bus.Publish(topic, resp[len(resp)-1])
	}
	if topic == "stocks" {
		var resp []stream.StockResponse
		if err := json.Unmarshal(b, &resp); err != nil {
			fmt.Println(err)
		}
		bus.Publish(topic, resp[len(resp)-1])
	}
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

	var cryptoRepo = prices_repository.NewCollection(client, "crypto", "crypto")
	var stocksRepo = prices_repository.NewCollection(client, "crypto", "stocks")

	var pricesPresenter = prices.NewPresenter(*cryptoRepo, *stocksRepo, bus)

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
	pricesPresenter.StartStream(wsConfig)

	http.HandleFunc("/crypto", func(rw http.ResponseWriter, r *http.Request) { pricesPresenter.Handler(rw, r, "crypto") })
	http.HandleFunc("/stocks", func(rw http.ResponseWriter, r *http.Request) { pricesPresenter.Handler(rw, r, "stocks") })
	go log.Fatal(http.ListenAndServe(*addr, nil))

	time.Sleep(time.Minute)
	pricesPresenter.StopStream()
}
