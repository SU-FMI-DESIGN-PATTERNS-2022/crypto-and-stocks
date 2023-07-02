package prices

import "github.com/gin-gonic/gin"

type PricesHandler struct {
	presenter Presenter
}

func NewPricesHandler(pricesPresenter Presenter) *PricesHandler {
	return &PricesHandler{
		presenter: pricesPresenter,
	}
}

func (pricesHandler *PricesHandler) GetCryptoPrices(context *gin.Context) {
	pricesHandler.presenter.CryptoHandler(context.Writer, context.Request)
}

func (pricesHandler *PricesHandler) GetStocksPrices(context *gin.Context) {
	pricesHandler.presenter.StockHandler(context.Writer, context.Request)
}
