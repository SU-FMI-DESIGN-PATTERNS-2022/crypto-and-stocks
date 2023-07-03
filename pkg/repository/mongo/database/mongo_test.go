package database

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
	"time"
)

// Create a new collection with the given name. The collection will be dropped when the test is over.

func TestFindOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		//cryptoCollection = mt.Coll
		//stockCollection = mt.Coll

		expectedStocks := StockPrices{
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
		}
		expectedCryptos := CryptoPrices{
			Prices: Prices{
				Symbol:   "BTC",
				BidPrice: 100.0,
				BidSize:  100.0,
				AskPrice: 100.0,
				AskSize:  100.0,
				Date:     time.Now(),
			},
			Exchange: "NASDAQ",
		}
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "crypto", mtest.FirstBatch, bson.D{
			{"_id", "1"},
			{"symbol", expectedCryptos.Symbol},
			{"bid_price", expectedCryptos.BidPrice},
			{"bid_size", expectedCryptos.BidSize},
			{"ask_price", expectedCryptos.AskPrice},
			{"ask_size", expectedCryptos.AskSize},
			{"date", expectedCryptos.Date},
			{"exchange", expectedCryptos.Exchange},
		}))
		cryptoResponse, err := cryptoCollection.GetMostRecentPriceBySymbol("BTC")
		assert.Nil(t, err)
		assert.Equal(t, expectedCryptos, cryptoResponse)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "stocks", mtest.FirstBatch, bson.D{
			{"_id", "1"},
			{"symbol", expectedStocks.Symbol},
			{"bid_price", expectedStocks.BidPrice},
			{"bid_size", expectedStocks.BidSize},
			{"ask_price", expectedStocks.AskPrice},
			{"ask_size", expectedStocks.AskSize},
			{"date", expectedStocks.Date},
			{"ask_exchange", expectedStocks.AskExchange},
			{"bid_exchange", expectedStocks.BidExchange},
			{"trade_size", expectedStocks.TradeSize},
			{"conditions", expectedStocks.Conditions},
			{"tape", expectedStocks.Tape},
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