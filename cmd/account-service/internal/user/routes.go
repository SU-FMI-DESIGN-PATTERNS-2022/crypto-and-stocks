package user

import "github.com/gin-gonic/gin"

func HandleRoutes(router *gin.RouterGroup, userPresenter UserPresenter) {
	userHandler := NewUserHandler(userPresenter)

	router.POST("/create/user/:id/:name", userHandler.CreateUser)
	router.POST("/create/bot/:id/:amount", userHandler.CreateBot)
	router.PUT("/merge/:id", userHandler.MergeUserAndBot)
	router.GET("/amount/:id", userHandler.EstimateUserAmount)
}
