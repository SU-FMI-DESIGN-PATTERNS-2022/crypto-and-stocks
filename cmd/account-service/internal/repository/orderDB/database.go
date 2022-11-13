package orderDB

import (
	"database/sql"
	"strconv"

	_ "github.com/lib/pq"
)

type Database struct {
	instance *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{
		instance: db,
	}
}

func (db *Database) StoreOrder(order Model) error {
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

func (db *Database) GetAllOrders() ([]Model, error) {
	rows, err := db.instance.Query(selectAllSQL)

	defer rows.Close()

	var orders []Model
	for rows.Next() {
		var order Model
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

func (db *Database) Close() {
	defer db.instance.Close()
}
