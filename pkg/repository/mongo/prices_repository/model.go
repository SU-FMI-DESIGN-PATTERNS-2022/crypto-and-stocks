package prices_repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Prices struct {
	ID       primitive.ObjectID `bson:"_id, omitempty"`
	Symbol   string             `bson:"symbol,omitempty"`
	Exchange string             `bson:"exchange,omitempty"`
	BidPrice float64            `bson:"bid_price,omitempty"`
	AskPrice float64            `bson:"ask_price,omitempty"`
	Date     time.Time          `bson:"date,omitempty"`
}
