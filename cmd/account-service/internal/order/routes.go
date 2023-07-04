package order

import (
	"github.com/gin-gonic/gin"
)

func HandleRoutes(router *gin.RouterGroup, orderPresenter Presenter) {
	router.GET("/all", orderPresenter.GetAllOrders)
	router.GET("/store", orderPresenter.StoreOrder)
	router.GET("/user/:id", orderPresenter.GetAllOrdersByUserId)
	router.GET("/symbol/:symbol", orderPresenter.GetAllOrdersBySymbol)
	router.GET("/:id/:symbol", orderPresenter.GetAllOrdersByUserIdAndSymbol)
}
