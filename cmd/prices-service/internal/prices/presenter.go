package prices

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (presenter *Presenter) StockHandler(context *gin.Context) {
	conn, err := presenter.upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	presenter.subscribeForResponding(conn, "stocks")
}

func (presenter *Presenter) CryptoHandler(context *gin.Context) {
	conn, err := presenter.upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	presenter.subscribeForResponding(conn, "crypto")
}

func (presenter *Presenter) subscribeForResponding(conn Connection, topic string) {
	err := presenter.bus.Subscribe(topic, func(resp interface{}) {
		err := conn.WriteJSON(resp)
		if err != nil {
			log.Println(err)
			return
		}
	})
	if err != nil {
		log.Println(err)
		return
	}
}
