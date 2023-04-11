package stream

import (
	"encoding/json"
	"fmt"
	"log"
)

//go:generate mockgen -source=controller.go -destination=mocks/controller.go

type EventBus interface {
	Publish(topic string, args ...interface{})
}

type PriceStream interface {
	Start(msgHandler func([]byte)) error
	Stop()
}

type Controller struct {
	cryptoStream PriceStream
	stockStream  PriceStream
	bus          EventBus
}

func NewController(cryptoStream PriceStream, stockStream PriceStream, bus EventBus) *Controller {
	return &Controller{
		cryptoStream: cryptoStream,
		stockStream:  stockStream,
		bus:          bus,
	}
}

func (c *Controller) StartStreamsToWrite() error {
	errChan := make(chan error)

	go func() {
		if err := c.cryptoStream.Start(c.publishInCrypto); err != nil {
			errChan <- fmt.Errorf("start crypto stream fails: %w", err)
			return
		}
	}()

	go func() {
		if err := c.stockStream.Start(c.publishInStocks); err != nil {
			errChan <- fmt.Errorf("start stocks stream fails: %w", err)
			return
		}
	}()

	return <-errChan
}

func (c *Controller) StopStreams() {
	c.cryptoStream.Stop()
	c.stockStream.Stop()
}

func (c *Controller) publishInCrypto(b []byte) {
	var resp []CryptoResponse

	if err := json.Unmarshal(b, &resp); err != nil {
		log.Println(err)
	}

	for _, price := range resp {
		c.bus.Publish("crypto", price)
	}
}

func (c *Controller) publishInStocks(b []byte) {
	var resp []StockResponse

	if err := json.Unmarshal(b, &resp); err != nil {
		log.Println(err)
	}

	for _, price := range resp {
		c.bus.Publish("stocks", price)
	}
}
