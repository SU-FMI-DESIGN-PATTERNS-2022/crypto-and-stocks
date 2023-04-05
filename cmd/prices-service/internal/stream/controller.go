package stream

import (
	"encoding/json"
	"log"
)

//go:generate mockgen -source=controller.go -destination=mocks/controller.go

type Bus interface {
	Publish(topic string, args ...interface{})
}

type controller struct {
	cryptoStream *Stream
	stockStream  *Stream
	bus          Bus
}

func NewController(cryptoStream *Stream, stockStream *Stream, bus Bus) *controller {
	return &controller{
		cryptoStream: cryptoStream,
		stockStream:  stockStream,
		bus:          bus,
	}
}

func (c *controller) StartStreamsToWrite() {
	go func() {
		if err := c.cryptoStream.Start(c.publishInCrypto); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := c.stockStream.Start(c.publishInStocks); err != nil {
			panic(err)
		}
	}()

}

func (c *controller) StopStreams() {
	c.cryptoStream.Stop()
	c.stockStream.Stop()
}

func (c *controller) publishInCrypto(b []byte) {
	var resp []CryptoResponse

	if err := json.Unmarshal(b, &resp); err != nil {
		log.Println(err)
	}

	for _, price := range resp {
		c.bus.Publish("crypto", price)
	}
}

func (c *controller) publishInStocks(b []byte) {
	var resp []StockResponse

	if err := json.Unmarshal(b, &resp); err != nil {
		log.Println(err)
	}

	for _, price := range resp {
		c.bus.Publish("stocks", price)
	}
}
