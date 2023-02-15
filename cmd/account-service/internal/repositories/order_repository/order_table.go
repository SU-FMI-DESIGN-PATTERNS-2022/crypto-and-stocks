package order_repository

import (
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
	orders := make([]Order, 0)
	err := table.instance.Select(&orders, query, args...)

	return orders, err
}

func (table *OrderTable) StoreOrder(order Order) error {
	_, err := table.instance.Exec(insertSQL,
		order.UserID,
		order.Type,
		order.Symbol,
		order.Amount,
		order.Price,
		order.Date,
	)

	return err
}

func (table *OrderTable) UpdateOrdersCreatorByUserId(prevUserId int64, newUserId int64) error {
	_, err := table.instance.Exec(updateOrdersAfterMergeSQL, newUserId, prevUserId)
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

func (table *OrderTable) GetAllOrdersByUserIdAndSymbol(userId int64, symbol string) ([]Order, error) {
	return table.getOrderRequest(selectAllWhereUserIdAndSymbolSQL, userId, symbol)
}
