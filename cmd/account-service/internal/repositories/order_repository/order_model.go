package order_repository

import "time"

type Order struct {
	ID     int64     `db:"id" json:"id"`
	UserID int64     `db:"user_id" json:"user_id"`
	Type   string    `db:"type" json:"type"`
	Symbol string    `db:"symbol" json:"symbol"`
	Amount float64   `db:"amount" json:"amount"`
	Price  float64   `db:"price" json:"price"`
	Date   time.Time `db:"date" json:"date"`
}
