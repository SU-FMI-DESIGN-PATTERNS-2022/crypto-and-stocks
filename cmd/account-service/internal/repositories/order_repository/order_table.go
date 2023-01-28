package order_repository

import (
	"database/sql"
	"strconv"

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

func (db *OrderTable) StoreOrder(order Order) error {
	_, err := db.instance.Exec(insertSQL,
		order.UserID,
		order.Type,
		order.Symbol,
		strconv.FormatFloat(order.Amount, 'E', -1, 64),
		strconv.FormatFloat(order.Price, 'E', -1, 64),
		order.Date.Format("2006-1-2"),
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
