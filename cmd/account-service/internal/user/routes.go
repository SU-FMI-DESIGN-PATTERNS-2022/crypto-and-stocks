package user

import "github.com/gin-gonic/gin"

func HandleRoutes(router *gin.RouterGroup, presenter Presenter) {
	router.POST("/create/user/:id/:name", presenter.CreateUser)
	router.POST("/create/bot/:id/:amount", presenter.CreateBot)
	router.PUT("/merge/:id", presenter.MergeUserAndBot)
	router.GET("/amount/:id", presenter.EstimateUserAmount)
}
