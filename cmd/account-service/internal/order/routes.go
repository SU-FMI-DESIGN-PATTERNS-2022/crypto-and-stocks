package order

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, orderController OrderController) {
	router.GET("/orders/all", orderController.GetAllOrders)
	router.GET("/orders/all/:id", orderController.GetAllOrdersByUserId)
	router.POST("/create/bot", orderController.CreateBot)
}
