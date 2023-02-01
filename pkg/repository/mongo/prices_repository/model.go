package prices_repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Prices struct {
	ID       primitive.ObjectID `bson:"_id, omitempty"`
	Symbol   string             `bson:"symbol,omitempty"`
	BidPrice float64            `bson:"bid_price,omitempty"`
	BidSize  float64            `bson:"bid_size,omitempty"`
	AskPrice float64            `bson:"ask_price,omitempty"`
	AskSize  float64            `bson:"ask_size,omitempty"`
	Date     time.Time          `bson:"date,omitempty"`
}

type CryptoPrice struct {
	Prices
	Exchange string `bson:"exchange,omitempty"`
}
type StockPrices struct {
	Prices
	AskExchange string   `bson:"ask_exchange, omitempty"`
	BidExchange string   `bson:"bid_exchange, omitempty"`
	TradeSize   float64  `bson:"trade_size, omitempty"`
	Conditions  []string `bson:"conditions, omitempty"`
	Tape        string   `bson:"tape, omitempty"`
}
