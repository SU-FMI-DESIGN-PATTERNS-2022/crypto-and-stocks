package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	presenter UserPresenter
}

func NewUserHandler(userPresenter UserPresenter) *UserHandler {
	return &UserHandler{
		presenter: userPresenter,
	}
}

func (userHandler *UserHandler) CreateUser(context *gin.Context) {
	name := context.Param("name")
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Parameter id should be number")
		return
	}

	if err = userHandler.presenter.CreateUser(id, name); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}

	context.JSON(http.StatusCreated, "User created successfully")
}

func (userHandler *UserHandler) CreateBot(context *gin.Context) {
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

	if err = userHandler.presenter.CreateBot(id, amount); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}

	context.JSON(http.StatusCreated, "Bot created successfully")
}

func (userHandler *UserHandler) MergeUserAndBot(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Parameter id should be number")
		return
	}

	if err = userHandler.presenter.MergeUserAndBot(id); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}

	context.JSON(http.StatusOK, "Successfully merged user and bot")
}

func (userHandler *UserHandler) EstimateUserAmount(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Parameter id should be number")
		return
	}

	amount, err := userHandler.presenter.EstimateUserAmount(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, "Could not find user with that id")
		return
	}

	context.JSON(http.StatusOK, amount)
}
