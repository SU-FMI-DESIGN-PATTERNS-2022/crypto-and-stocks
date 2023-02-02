package crypto_prices_repository

import (
	"context"
	"time"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CryptoPricesCollection struct {
	database.Collection
}

func NewCryptoPricesCollection(client *mongo.Client, db string, col string) *CryptoPricesCollection {
	return &CryptoPricesCollection{
		Collection: database.Collection{
			Instance:       client,
			Database:       db,
			CollectionName: col,
		},
	}
}

func (c *CryptoPricesCollection) StoreEntry(price CryptoPrices) error {
	_, err := c.InsertOne(c.CollectionName, price)
	return err
}

func (c *CryptoPricesCollection) GetAllPrices() ([]CryptoPrices, error) {
	collection := c.Instance.Database(c.Database).Collection(c.CollectionName)

	result, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		panic(err)
	}
	var prices []CryptoPrices

	if err = result.All(context.TODO(), &prices); err != nil {
		panic(err)
	}
	return prices, err
}

func (c *CryptoPricesCollection) GetAllPricesBySymbol(symbol string) ([]CryptoPrices, error) {
	collection := c.Instance.Database(c.Database).Collection(c.CollectionName)

	filter := bson.D{primitive.E{Key: "symbol", Value: symbol}}

	result, err := collection.Find(context.TODO(), filter)

	if err != nil {
		panic(err)
	}
	var prices []CryptoPrices

	if err = result.All(context.TODO(), &prices); err != nil {
		panic(err)
	}
	return prices, err
}

func (c *CryptoPricesCollection) GetAllPricesByExchange(exchange string) ([]CryptoPrices, error) {
	collection := c.Instance.Database(c.Database).Collection(c.CollectionName)

	filter := bson.D{primitive.E{Key: "exchange", Value: exchange}}

	result, err := collection.Find(context.TODO(), filter)

	if err != nil {
		panic(err)
	}

	var prices []CryptoPrices

	if err = result.All(context.TODO(), &prices); err != nil {
		panic(err)
	}

	return prices, err
}

func (c *CryptoPricesCollection) GetAllPricesInPeriod(from time.Time, to time.Time) ([]CryptoPrices, error) {
	collection := c.Instance.Database(c.Database).Collection(c.CollectionName)

	filter := bson.M{"date": bson.M{
		"$gte": primitive.NewDateTimeFromTime(from),
		"$lte": primitive.NewDateTimeFromTime(to)}}

	result, err := collection.Find(context.TODO(), filter)

	var prices []CryptoPrices

	if err = result.All(context.TODO(), &prices); err != nil {
		panic(err)
	}

	return prices, err
}

func (c *CryptoPricesCollection) GetAllPricesInPeriodSymbol(from time.Time, to time.Time, symbol string) ([]CryptoPrices, error) {
	collection := c.Instance.Database(c.Database).Collection(c.CollectionName)

	filter := bson.M{"date": bson.M{
		"$gte": primitive.NewDateTimeFromTime(from),
		"$lte": primitive.NewDateTimeFromTime(to)},
		"symbol": symbol,
	}

	result, err := collection.Find(context.TODO(), filter)

	var prices []CryptoPrices

	if err = result.All(context.TODO(), &prices); err != nil {
		panic(err)
	}

	return prices, err
}

func (c *CryptoPricesCollection) GetMostRecentPriceBySymbol(symbol string) (CryptoPrices, error) {
	collection := c.Instance.Database(c.Database).Collection(c.CollectionName)

	filter := bson.D{primitive.E{Key: "symbol", Value: symbol}}
	opts := options.FindOne().SetSort(bson.M{"$natural": -1})

	var lastRecord CryptoPrices
	err := collection.FindOne(context.TODO(), filter, opts).Decode(&lastRecord)
	if err != nil {
		panic(err)
	}

	return lastRecord, err
}
