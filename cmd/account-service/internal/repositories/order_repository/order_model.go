package order_repository

import "time"

type Order struct {
	ID     int64     `db:"id"`
	UserID int64     `db:"user_id"`
	Type   string    `db:"type"`
	Symbol string    `db:"symbol"`
	Amount float64   `db:"amount"`
	Price  float64   `db:"price"`
	Date   time.Time `db:"date"`
}
