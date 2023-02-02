package database

import (
	"time"
)

type Prices struct {
	Symbol   string    `bson:"symbol,omitempty"`
	BidPrice float64   `bson:"bid_price,omitempty"`
	BidSize  float64   `bson:"bid_size,omitempty"`
	AskPrice float64   `bson:"ask_price,omitempty"`
	AskSize  float64   `bson:"ask_size,omitempty"`
	Date     time.Time `bson:"date,omitempty"`
}
