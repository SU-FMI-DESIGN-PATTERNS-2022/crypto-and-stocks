package order_repository

import (
	"database/sql"
	"strconv"
	"time"

	"errors"

	_ "github.com/lib/pq"
)

type OrderTable struct {
	instance *sql.DB
}

func NewOrderTable(db *sql.DB) *OrderTable {
	return &OrderTable{
		instance: db,
	}
}

func (db *OrderTable) getOrderRequest(query string, args ...any) ([]Order, error) {
	rows, err := db.instance.Query(query, args...)

	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Type,
			&order.Symbol,
			&order.Amount,
			&order.Price,
			&order.Date,
		)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return orders, nil
}

func (db *OrderTable) StoreOrder(userId int64, orderType string, symbol string, amount float64, price float64) error {
	if orderType == "buy" {
		row := db.instance.QueryRow(selectUserAmountWhereUserIdSQL, userId)

		var userAmount float64

		userErr := row.Scan(&userAmount)

		if userErr != nil {
			return userErr
		}

		if userAmount < amount*price {
			return errors.New("Insufficient amount")
		}
	}

	//TODO: Validate if user has enough crypto to sell

	_, err := db.instance.Exec(insertSQL,
		userId,
		orderType,
		symbol,
		strconv.FormatFloat(amount, 'E', -1, 64),
		strconv.FormatFloat(price, 'E', -1, 64),
		time.Now(),
	)

	return err
}

func (db *OrderTable) GetAllOrders() ([]Order, error) {
	return db.getOrderRequest(selectAllSQL)
}

func (db *OrderTable) GetAllOrdersBySymbol(symbol string) ([]Order, error) {
	return db.getOrderRequest(selectAllWhereSymbolSQL, symbol)
}

func (db *OrderTable) GetAllOrdersByUserId(userId int64) ([]Order, error) {
	return db.getOrderRequest(selectAllWhereUserIdSQL, userId)
}
