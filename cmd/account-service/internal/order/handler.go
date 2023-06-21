package order

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/order_repository"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error)
}

type OrderHandler struct {
	presenter OrderPresenter
	upgrader  Upgrader
}

func NewOrderHandler(orderPresenter OrderPresenter, upgrader Upgrader) OrderHandler {
	return OrderHandler{
		presenter: orderPresenter,
		upgrader:  upgrader,
	}
}

func (orderHandler *OrderHandler) GetAllOrders(context *gin.Context) {
	orders, err := orderHandler.presenter.GetAllOrders()
	if err != nil {
		context.JSON(http.StatusInternalServerError, "Could not fetch orders"+err.Error())
		return
	}

	context.JSON(http.StatusOK, orders)
}

func (orderHandler *OrderHandler) StoreOrder(context *gin.Context) {
	conn, err := orderHandler.upgrader.Upgrade(context.Writer, context.Request, nil)

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

		if err = orderHandler.presenter.StoreOrder(order); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("We have a problem with storing your order."))
			break
		}
	}
}

func (orderHandler *OrderHandler) GetAllOrdersByUserId(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Parameter id should be number")
		return
	}

	orders, err := orderHandler.presenter.GetAllOrdersByUserId(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Could not fetch orders"+err.Error())
		return
	}

	context.JSON(http.StatusOK, orders)
}

func (orderHandler *OrderHandler) GetAllOrdersBySymbol(context *gin.Context) {
	symbol := context.Param("symbol")

	orders, err := orderHandler.presenter.GetAllOrdersBySymbol(symbol)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Could not fetch orders"+err.Error())
		return
	}

	context.JSON(http.StatusOK, orders)
}

func (orderHandler *OrderHandler) GetAllOrdersByUserIdAndSymbol(context *gin.Context) {
	symbol := context.Param("symbol")
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Parameter id should be number")
		return
	}

	orders, err := orderHandler.presenter.GetAllOrdersByUserIdAndSymbol(id, symbol)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Could not fetch orders:"+err.Error())
		return
	}

	context.JSON(http.StatusOK, orders)
}
