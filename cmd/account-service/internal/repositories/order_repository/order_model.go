package order_repository

import "time"

type Order struct {
	ID     int64
	UserID int64
	Type   string
	Symbol string
	Amount float64
	Price  float64
	Date   time.Time
}
