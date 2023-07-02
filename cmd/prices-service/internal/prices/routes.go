package prices

import "github.com/gin-gonic/gin"

func HandleRoutes(router *gin.RouterGroup, pricesPresenter *Presenter) {
	router.GET("/crypto", pricesPresenter.CryptoHandler)
	router.GET("/stocks", pricesPresenter.StockHandler)
}
