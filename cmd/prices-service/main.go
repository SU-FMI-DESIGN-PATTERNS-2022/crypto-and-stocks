package main

import (
	"context"

	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
	mongoEnv "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"

	"github.com/asaskevich/EventBus"
	"github.com/gorilla/websocket"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/prices"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/stream"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

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

func main() {
	bus := EventBus.New()
	mongoConfig := mongoEnv.LoadMongoConfig()
	client, err := database.Connect(mongoConfig, database.Remote)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	var cryptoRepo = database.NewCollection[database.CryptoPrices](client, mongoConfig.Database, "CryptoPrices")
	var stocksRepo = database.NewCollection[database.StockPrices](client, mongoConfig.Database, "StockPrices")

	var pricesPresenter = prices.NewPresenter(cryptoRepo, stocksRepo, bus)

	wsConfig := env.LoadWebSocetConfig()
	pricesPresenter.StartStream(wsConfig)

	http.HandleFunc("/crypto", pricesPresenter.CryptoHandler)
	http.HandleFunc("/stocks", pricesPresenter.StockHandler)
	go log.Fatal(http.ListenAndServe(*addr, nil))

	time.Sleep(time.Minute)
	pricesPresenter.StopStream()
}
