package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestFindOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		cryptoCollection := NewCollection[CryptoPrices](mt.Client, "CryptoStocks", "CryptoPrices")
		stocksCollection := NewCollection[StockPrices](mt.Client, "CryptoStocks", "StockPrices")

		expectedStocks := StockPrices{
			Prices: Prices{
				Symbol:   "AAPL",
				BidPrice: 100.0,
				BidSize:  100.0,
				AskPrice: 100.0,
				AskSize:  100.0,
			},
			AskExchange: "NASDAQ",
			BidExchange: "NASDAQ",
			TradeSize:   100.0,
			Conditions:  []string{"A", "B"},
			Tape:        "A",
		}
		expectedCryptos := CryptoPrices{
			Prices: Prices{
				Symbol:   "BTC",
				BidPrice: 100.0,
				BidSize:  100.0,
				AskPrice: 100.0,
				AskSize:  100.0,
			},
			Exchange: "NASDAQ",
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "db.crypto", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: "1"},
			{Key: "prices", Value: bson.D{
				{Key: "symbol", Value: expectedCryptos.Prices.Symbol},
				{Key: "bid_price", Value: expectedCryptos.Prices.BidPrice},
				{Key: "bid_size", Value: expectedCryptos.Prices.BidSize},
				{Key: "ask_price", Value: expectedCryptos.Prices.AskPrice},
				{Key: "ask_size", Value: expectedCryptos.Prices.AskSize},
			}},
			{Key: "exchange", Value: expectedCryptos.Exchange},
		}))
		cryptoResponse, err := cryptoCollection.GetMostRecentPriceBySymbol("BTC")

		assert.Nil(t, err)
		assert.Equal(t, expectedCryptos, cryptoResponse)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "db.stocks", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: "1"},
			{Key: "prices", Value: bson.D{
				{Key: "symbol", Value: expectedStocks.Symbol},
				{Key: "bid_price", Value: expectedStocks.BidPrice},
				{Key: "bid_size", Value: expectedStocks.BidSize},
				{Key: "ask_price", Value: expectedStocks.AskPrice},
				{Key: "ask_size", Value: expectedStocks.AskSize},
			}},
			{Key: "ask_exchange", Value: expectedStocks.AskExchange},
			{Key: "bid_exchange", Value: expectedStocks.BidExchange},
			{Key: "trade_size", Value: expectedStocks.TradeSize},
			{Key: "conditions", Value: expectedStocks.Conditions},
			{Key: "tape", Value: expectedStocks.Tape},
		}))
		stockResponse, err := stocksCollection.GetMostRecentPriceBySymbol("AAPL")
		assert.Nil(t, err)
		assert.Equal(t, expectedStocks, stockResponse)

	})
}

func TestInsertOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		stocksCollection := NewCollection[StockPrices](mt.Client, "CryptoStocks", "StockPrices")

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := stocksCollection.StoreEntry(
			StockPrices{
				Prices: Prices{
					Symbol:   "AAPL",
					BidPrice: 100.0,
					BidSize:  100.0,
					AskPrice: 100.0,
					AskSize:  100.0,
					Date:     time.Now(),
				},
				AskExchange: "NASDAQ",
				BidExchange: "NASDAQ",
				TradeSize:   100.0,
				Conditions:  []string{"A", "B"},
				Tape:        "A",
			})
		assert.Nil(t, err)
	})
}

func TestInsertMany(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		stocksCollection := NewCollection[StockPrices](mt.Client, "CryptoStocks", "StockPrices")

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		_, err := stocksCollection.insertMany("StockPrices", []interface{}{
			StockPrices{
				Prices: Prices{
					Symbol:   "AAPL",
					BidPrice: 100.0,
					BidSize:  100.0,
					AskPrice: 100.0,
					AskSize:  100.0,
					Date:     time.Now(),
				},
				AskExchange: "NASDAQ",
				BidExchange: "NASDAQ",
				TradeSize:   100.0,
				Conditions:  []string{"A", "B"},
				Tape:        "A",
			},
			StockPrices{
				Prices: Prices{
					Symbol:   "TSLA",
					BidPrice: 100.0,
					BidSize:  100.0,
					AskPrice: 100.0,
					AskSize:  100.0,
					Date:     time.Now(),
				},
				AskExchange: "NASDAQ",
				BidExchange: "NASDAQ",
				TradeSize:   100.0,
				Conditions:  []string{"A", "B"},
				Tape:        "A",
			},
		})
		assert.Nil(t, err)
	})
}
