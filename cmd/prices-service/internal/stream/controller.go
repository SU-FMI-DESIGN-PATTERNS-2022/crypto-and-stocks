package stream

import (
	"encoding/json"
)

//go:generate mockgen -source=controller.go -destination=mocks/controller.go

type EventBus interface {
	Publish(topic string, args ...interface{})
}

type PriceStream interface {
	Start(msgHandler func([]byte) error) error
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

func (c *Controller) StartStreamsToWrite() <-chan error {
	errCh := make(chan error, 2)
	go func() {
		if err := c.cryptoStream.Start(c.publishInCrypto); err != nil {
			errCh <- err
		}
	}()

	go func() {
		if err := c.stockStream.Start(c.publishInStocks); err != nil {
			errCh <- err
		}
	}()
	return errCh
}

func (c *Controller) StopStreams() {
	c.cryptoStream.Stop()
	c.stockStream.Stop()
}

func (c *Controller) publishInCrypto(b []byte) error {
	var resp []CryptoResponse

	if err := json.Unmarshal(b, &resp); err != nil {
		// log.Println(err)
		return err
	}

	for _, price := range resp {
		c.bus.Publish("crypto", price)
	}

	return nil
}

func (c *Controller) publishInStocks(b []byte) error {
	var resp []StockResponse

	if err := json.Unmarshal(b, &resp); err != nil {
		// log.Println(err)
		return err
	}

	for _, price := range resp {
		c.bus.Publish("stocks", price)
	}

	return nil
}
