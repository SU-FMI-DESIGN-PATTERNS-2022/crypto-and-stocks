package orderDB

import "time"

type Model struct {
	ID     int
	UserID int
	Type   string
	Symbol string
	Amount float64
	Price  float64
	Date   time.Time
}
