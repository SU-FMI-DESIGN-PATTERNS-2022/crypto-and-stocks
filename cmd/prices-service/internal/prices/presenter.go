package prices

import (
	"log"
	"net/http"
)

//go:generate mockgen -source=presenter.go -destination=mocks/presenter.go

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Connection, error)
}

type EventBus interface {
	Subscribe(topic string, fn interface{}) error
}

type Connection interface {
	WriteJSON(v interface{}) error
}

type Presenter struct {
	upgrader Upgrader
	bus      EventBus
}

func NewPresenter(upgrader Upgrader, bus EventBus) *Presenter {
	return &Presenter{
		upgrader: upgrader,
		bus:      bus,
	}
}

func (p *Presenter) StockHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := p.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	p.subscribeForResponding(conn, "stocks")
}

func (p *Presenter) CryptoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := p.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	p.subscribeForResponding(conn, "crypto")
}

func (p *Presenter) subscribeForResponding(conn Connection, topic string) {
	p.bus.Subscribe(topic, func(resp interface{}) {
		conn.WriteJSON(resp)
	})
}
