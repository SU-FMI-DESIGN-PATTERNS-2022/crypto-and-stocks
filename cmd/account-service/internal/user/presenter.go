package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Presenter struct {
	controller Controller
}

func NewPresenter(controller Controller) *Presenter {
	return &Presenter{
		controller,
	}
}

func (presenter *Presenter) CreateUser(context *gin.Context) {
	name := context.Param("name")
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Parameter id should be number")
		return
	}

	if err = presenter.controller.CreateUser(id, name); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}

	context.JSON(http.StatusCreated, "User created successfully")
}

func (presenter *Presenter) CreateBot(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Parameter id should be number")
		return
	}

	amount, err := strconv.ParseFloat(context.Param("id"), 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Parameter amount should be number")
		return
	}

	if err = presenter.controller.CreateBot(id, amount); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}

	context.JSON(http.StatusCreated, "Bot created successfully")
}

func (presenter *Presenter) MergeUserAndBot(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Parameter id should be number")
		return
	}

	if err = presenter.controller.MergeUserAndBot(id); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}

	context.JSON(http.StatusOK, "Successfully merged user and bot")
}

func (presenter *Presenter) EstimateUserAmount(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Parameter id should be number")
		return
	}

	amount, err := presenter.controller.EstimateUserAmount(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Could not find user with that id")
		return
	}

	context.JSON(http.StatusOK, amount)
}
