package prices_repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewCryptoPrice(symbol string, exchange string, bidPrice float64, askPrice float64, date time.Time) CryptoPrice {
	return CryptoPrice{
		Prices: Prices{
			ID:       primitive.NewObjectID(),
			Symbol:   symbol,
			Exchange: exchange,
			BidPrice: bidPrice,
			AskPrice: askPrice,
			Date:     date,
		},
	}

}
func NewPrice(symbol string, exchange string, bidPrice float64, askPrice float64, date time.Time) Prices {
	return Prices{
		ID:       primitive.NewObjectID(),
		Symbol:   symbol,
		Exchange: exchange,
		BidPrice: bidPrice,
		AskPrice: askPrice,
		Date:     date,
	}

}

func NewStockPrice(price Prices, askExchange string, bidExchange string, tradeSize float64, conditions []string, tape string) StockPrice {
	return StockPrice{
		Prices:      price,
		AskExchange: askExchange,
		BidExchange: bidExchange,
		TradeSize:   tradeSize,
		Conditions:  conditions,
		Tape:        tape,
	}
}

type Prices struct {
	ID       primitive.ObjectID `bson:"_id, omitempty"`
	Symbol   string             `bson:"symbol,omitempty"`
	Exchange string             `bson:"exchange,omitempty"`
	BidPrice float64            `bson:"bid_price,omitempty"`
	AskPrice float64            `bson:"ask_price,omitempty"`
	Date     time.Time          `bson:"date,omitempty"`
}

type CryptoPrice struct {
	Prices
}
type StockPrice struct {
	Prices
	AskExchange string   `bson:"ask_exchange,omitempty"`
	BidExchange string   `bson:"bid_exchange,omitempty"`
	TradeSize   float64  `bson:"trade_size,omitempty"`
	Conditions  []string `bson:"conditions,omitempty"`
	Tape        string   `bson:"tape,omitempty"`
}
