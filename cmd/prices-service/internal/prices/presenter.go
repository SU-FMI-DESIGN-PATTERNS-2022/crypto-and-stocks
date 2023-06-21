package prices

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//go:generate mockgen -source=presenter.go -destination=mocks/presenter.go

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error)
}

type EventBus interface {
	Subscribe(topic string, fn interface{}) error
}

type PricesPresenter struct {
	upgrader Upgrader
	bus      EventBus
}

func NewPresenter(upgrader Upgrader, bus EventBus) *PricesPresenter {
	return &PricesPresenter{
		upgrader: upgrader,
		bus:      bus,
	}
}

func (p *PricesPresenter) StockHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := p.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	p.subscribeForResponding(conn, "stocks")
}

func (p *PricesPresenter) CryptoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := p.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	p.subscribeForResponding(conn, "crypto")
}

func (p *PricesPresenter) subscribeForResponding(conn *websocket.Conn, topic string) {
	p.bus.Subscribe(topic, func(resp interface{}) {
		conn.WriteJSON(resp)
	})
}
