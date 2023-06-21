package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleRoutes(router *gin.RouterGroup, orderPresenter OrderPresenter) {
	orderHandler := NewOrderHandler(orderPresenter, &upgrader)

	router.GET("/all", orderHandler.GetAllOrders)
	router.GET("/store", orderHandler.StoreOrder)
	router.GET("/user/:id", orderHandler.GetAllOrdersByUserId)
	router.GET("/symbol/:symbol", orderHandler.GetAllOrdersBySymbol)
	router.GET("/:id/:symbol", orderHandler.GetAllOrdersByUserIdAndSymbol)
}
