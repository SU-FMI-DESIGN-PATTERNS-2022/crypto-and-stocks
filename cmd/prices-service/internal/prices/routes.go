package prices

import "github.com/gin-gonic/gin"

func HandleRoutes(router *gin.RouterGroup, pricesPresenter Presenter) {
	pricesHandler := NewPricesHandler(pricesPresenter)

	router.GET("/crypto", pricesHandler.GetCryptoPrices)
	router.GET("/stocks", pricesHandler.GetStocksPrices)
}
