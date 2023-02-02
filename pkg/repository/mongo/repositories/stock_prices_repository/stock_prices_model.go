package stock_prices_repository

import "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"

type StockPrices struct {
	database.Prices
	AskExchange string   `bson:"ask_exchange, omitempty"`
	BidExchange string   `bson:"bid_exchange, omitempty"`
	TradeSize   float64  `bson:"trade_size, omitempty"`
	Conditions  []string `bson:"conditions, omitempty"`
	Tape        string   `bson:"tape, omitempty"`
}
