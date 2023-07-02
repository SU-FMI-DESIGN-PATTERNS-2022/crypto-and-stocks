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

func (controller *Controller) StartStreamsToWrite() error {
	errCh := make(chan error, 1)

	go func() {
		if err := controller.cryptoStream.Start(controller.publishInCrypto); err != nil {
			errCh <- err
		}
	}()

	go func() {
		if err := controller.stockStream.Start(controller.publishInStocks); err != nil {
			errCh <- err
		}
	}()

	return <-errCh
}

func (controller *Controller) StopStreams() {
	controller.cryptoStream.Stop()
	controller.stockStream.Stop()
}

func (controller *Controller) publishInCrypto(b []byte) error {
	var resp []CryptoResponse

	if err := json.Unmarshal(b, &resp); err != nil {
		return err
	}

	for _, price := range resp {
		controller.bus.Publish("crypto", price)
	}

	return nil
}

func (controller *Controller) publishInStocks(b []byte) error {
	var resp []StockResponse

	if err := json.Unmarshal(b, &resp); err != nil {
		return err
	}

	for _, price := range resp {
		controller.bus.Publish("stocks", price)
	}

	return nil
}
