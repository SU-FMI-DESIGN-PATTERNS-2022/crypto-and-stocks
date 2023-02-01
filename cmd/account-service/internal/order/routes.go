package order

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, orderController OrderController) {
	router.GET("/orders/all", orderController.GetAllOrders)
	router.GET("/orders/user/:id", orderController.GetAllOrdersByUserId)
	router.GET("/orders/user/:id/:symbol", orderController.GetAllOrdersByUserIdAndSymbol)
	router.GET("/orders/symbol/:symbol", orderController.GetAllOrdersBySymbol)
	router.POST("/create/user", orderController.CreateUser)
	router.POST("/create/bot", orderController.CreateBot)
	router.PUT("/merge", orderController.MergeUserAndBot)
}
