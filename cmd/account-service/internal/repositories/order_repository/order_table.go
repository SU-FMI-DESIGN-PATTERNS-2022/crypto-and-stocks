package order_repository

import (
	"strconv"
	"time"

	"errors"
	"math"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type OrderTable struct {
	instance *sqlx.DB
}

func NewOrderTable(db *sqlx.DB) *OrderTable {
	return &OrderTable{
		instance: db,
	}
}

func (table *OrderTable) getOrderRequest(query string, args ...any) ([]Order, error) {
	var orders []Order
	err := table.instance.Select(&orders, query, args...)

	return orders, err
}

func (table *OrderTable) StoreOrder(userId int64, orderType string, symbol string, amount float64, price float64) error {
	if orderType == "buy" {
		var userAmount float64
		userErr := table.instance.Get(&userAmount, selectUserAmountWhereUserIdSQL, userId)

		if userErr != nil {
			return userErr
		}

		if userAmount < amount*price {
			return errors.New("Insufficient amount")
		}

		_, amountErr := table.instance.Exec(updateUserAmountSQL, math.Round((userAmount-amount*price)*100)/100, userId)

		if amountErr != nil {
			return amountErr
		}
	}

	//TODO: Validate if user has enough crypto to sell

	_, err := table.instance.Exec(insertSQL,
		userId,
		orderType,
		symbol,
		strconv.FormatFloat(amount, 'E', -1, 64),
		strconv.FormatFloat(price, 'E', -1, 64),
		time.Now(),
	)

	return err
}

func (table *OrderTable) GetAllOrders() ([]Order, error) {
	return table.getOrderRequest(selectAllSQL)
}

func (table *OrderTable) GetAllOrdersBySymbol(symbol string) ([]Order, error) {
	return table.getOrderRequest(selectAllWhereSymbolSQL, symbol)
}

func (table *OrderTable) GetAllOrdersByUserId(userId int64) ([]Order, error) {
	return table.getOrderRequest(selectAllWhereUserIdSQL, userId)
}
