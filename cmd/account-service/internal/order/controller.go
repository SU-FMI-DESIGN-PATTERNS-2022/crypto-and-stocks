package order

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderPresenter OrderPresenter
}

func NewOrderController(orderPresenter OrderPresenter) OrderController {
	return OrderController{
		orderPresenter: orderPresenter,
	}
}

func (orderController *OrderController) GetAllOrders(context *gin.Context) {
	orders, err := orderController.orderPresenter.GetAllOrders()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, orders)
}

func (orderController *OrderController) GetAllOrdersByUserId(context *gin.Context) {
	idParam := context.Param("id")
	id, paramErr := strconv.ParseInt(idParam, 10, 64)
	if paramErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Paramater should be a number",
		})
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	orders, err := orderController.orderPresenter.GetAllOrdersByUserId(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, orders)
}

func (orderController *OrderController) CreateBot(context *gin.Context) {
	creatorIdQuery := context.Query("id")
	creatorId, err := strconv.ParseInt(creatorIdQuery, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": `Query parameter "id" should be a number`,
		})
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	amountQuery := context.Query("amount")
	amount, err := strconv.ParseFloat(amountQuery, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": `Query parameter "amount" should be a float`,
		})
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	reqErr := orderController.orderPresenter.CreateBot(creatorId, amount)
	if reqErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": reqErr.Error(),
		})
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "Successfully created bot",
	})
}
