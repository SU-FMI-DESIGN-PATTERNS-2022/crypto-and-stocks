package order

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/order_repository"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Connection, error)
}

type Connection interface {
	Close() error
	ReadMessage() (messageType int, payload []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

type OrderRepository interface {
	StoreOrder(order order_repository.Order) error
	GetAllOrders() ([]order_repository.Order, error)
	GetAllOrdersByUserId(userId int64) ([]order_repository.Order, error)
	GetAllOrdersBySymbol(symbol string) ([]order_repository.Order, error)
	GetAllOrdersByUserIdAndSymbol(userId int64, symbol string) ([]order_repository.Order, error)
}

type UserRepository interface {
	GetUserAmount(id int64) (float64, error)
	UpdateUserAmount(id int64, amount float64) error
}

type Presenter struct {
	orderRepo OrderRepository
	userRepo  UserRepository
	upgrader  Upgrader
}

func NewPresenter(orderRepo OrderRepository, userRepo UserRepository, upgrader Upgrader) *Presenter {
	return &Presenter{
		orderRepo: orderRepo,
		userRepo:  userRepo,
		upgrader:  upgrader,
	}
}

func (presenter *Presenter) GetAllOrders(context *gin.Context) {
	orders, err := presenter.orderRepo.GetAllOrders()
	if err != nil {
		context.JSON(http.StatusInternalServerError, "Could not fetch orders"+err.Error())
		return
	}

	context.JSON(http.StatusOK, orders)
}

func (presenter *Presenter) StoreOrder(context *gin.Context) {
	conn, err := presenter.upgrader.Upgrade(context.Writer, context.Request, nil)

	if err != nil {
		return
	}

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Something is wrong."))
			break
		}

		var order order_repository.Order
		if err = json.Unmarshal(message, &order); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("The message is not in right json object structure."))
			break
		}

		if err = presenter.storeOrder(order); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("We have a problem with storing your order."))
			break
		}
	}
}

func (presenter *Presenter) GetAllOrdersByUserId(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Parameter id should be number")
		return
	}

	orders, err := presenter.orderRepo.GetAllOrdersByUserId(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Could not fetch orders"+err.Error())
		return
	}

	context.JSON(http.StatusOK, orders)
}

func (presenter *Presenter) GetAllOrdersBySymbol(context *gin.Context) {
	symbol := context.Param("symbol")

	orders, err := presenter.orderRepo.GetAllOrdersBySymbol(symbol)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Could not fetch orders"+err.Error())
		return
	}

	context.JSON(http.StatusOK, orders)
}

func (presenter *Presenter) GetAllOrdersByUserIdAndSymbol(context *gin.Context) {
	symbol := context.Param("symbol")
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Parameter id should be number")
		return
	}

	orders, err := presenter.orderRepo.GetAllOrdersByUserIdAndSymbol(id, symbol)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Could not fetch orders:"+err.Error())
		return
	}

	context.JSON(http.StatusOK, orders)
}

func (presenter *Presenter) storeOrder(order order_repository.Order) error {
	switch order.Type {
	case "buy":
		if err := presenter.storeBuyOrder(order); err != nil {
			return err
		}
	case "sell":
		if err := presenter.storeSellOrder(order); err != nil {
			return err
		}
	default:
		return errors.New("invalid order type")
	}

	return presenter.orderRepo.StoreOrder(order)
}

func (presenter *Presenter) storeBuyOrder(order order_repository.Order) error {
	amount, err := presenter.userRepo.GetUserAmount(order.UserID)
	if err != nil {
		return err
	}

	if amount < order.Amount*order.Price {
		return errors.New("not enough amount")
	}

	updateErr := presenter.userRepo.UpdateUserAmount(order.UserID, amount-order.Amount*order.Price)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (presenter *Presenter) storeSellOrder(order order_repository.Order) error {
	orders, err := presenter.orderRepo.GetAllOrdersByUserIdAndSymbol(order.UserID, order.Symbol)
	if err != nil {
		return err
	}

	var amount float64
	for _, o := range orders {
		if o.Type == "buy" {
			amount += o.Amount
		} else {
			amount -= o.Amount
		}
	}

	if amount < order.Amount {
		return errors.New("not enough amount")
	}

	updateErr := presenter.userRepo.UpdateUserAmount(order.UserID, amount+order.Amount*order.Price)
	if updateErr != nil {
		return updateErr
	}

	return nil
}
