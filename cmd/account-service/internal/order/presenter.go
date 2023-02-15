package order

import (
	"encoding/json"
	"fmt"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/order_repository"
	"github.com/gorilla/websocket"
	"net/http"
)

type OrderRepository interface {
	StoreOrder(order order_repository.Order) error
	GetAllOrders() ([]order_repository.Order, error)
	GetAllOrdersByUserId(userId int64) ([]order_repository.Order, error)
	GetAllOrdersBySymbol(symbol string) ([]order_repository.Order, error)
}

type UserRepository interface {
	CreateUser(userId int64, name string) error
	CreateBot(creatorID int64, amount float64) error
	AddOrder(userId int64, orderId int64) error
	MergeUserOrders(id int64) error
}

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error)
}

type Presenter struct {
	orderRepo OrderRepository
	userRepo  UserRepository
	upgrader  Upgrader
}

func NewPresenter(orderRepo OrderRepository, userRepo UserRepository, upgrader Upgrader) Presenter {
	return Presenter{
		orderRepo: orderRepo,
		userRepo:  userRepo,
		upgrader:  upgrader,
	}
}

func (p *Presenter) StoreOrder(w http.ResponseWriter, r *http.Request) {
	conn, err := p.upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Failed to upgrade connection:", err)
		return
	}

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage() // json format
		if err != nil {
			fmt.Println("Failed to read message:", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Hello, something is wrong."))
			break
		}
		// deserialize to struct Order
		var order order_repository.Order
		if err := json.Unmarshal(message, &order); err != nil {
			fmt.Println("Failed to unmarshal message:", err)
			conn.WriteMessage(websocket.TextMessage, []byte("The message is not in right json object structure."))
			break
		}
		if err := p.orderRepo.StoreOrder(order); err != nil {
			fmt.Println("Failed to store order:", err)
			conn.WriteMessage(websocket.TextMessage, []byte("We have a problem with storing your order."))
			break
		}
	}
}
