package database

import (
	"context"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection[Prices CryptoPrices | StockPrices] struct {
	instance       *mongo.Client
	database       string
	collectionName string
}

func NewCollection[Prices CryptoPrices | StockPrices](client *mongo.Client, db string, col string) *Collection[Prices] {
	return &Collection[Prices]{
		instance:       client,
		database:       db,
		collectionName: col,
	}
}

func (c *Collection[Prices]) insertOne(col string, doc interface{}) (*mongo.InsertOneResult, error) {

	collection := c.instance.Database(c.database).Collection(col)

	result, err := collection.InsertOne(context.TODO(), doc)
	return result, err
}

// TODO: Yet to be used
func (c *Collection[Prices]) insertMany(col string, docs []interface{}) (*mongo.InsertManyResult, error) {

	collection := c.instance.Database(c.database).Collection(col)

	result, err := collection.InsertMany(context.TODO(), docs)
	return result, err
}

func (c *Collection[Prices]) StoreEntry(price Prices) error {
	_, err := c.insertOne(c.collectionName, price)
	return err
}

func (c *Collection[Prices]) GetAllPrices() ([]Prices, error) {
	collection := c.instance.Database(c.database).Collection(c.collectionName)

	result, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		panic(err)
	}
	var prices []Prices

	if err = result.All(context.TODO(), &prices); err != nil {
		panic(err)
	}
	return prices, err
}

func (c *Collection[Prices]) GetAllPricesBySymbol(symbol string) ([]Prices, error) {
	collection := c.instance.Database(c.database).Collection(c.collectionName)

	filter := bson.D{primitive.E{Key: "symbol", Value: symbol}}

	result, err := collection.Find(context.TODO(), filter)

	if err != nil {
		panic(err)
	}
	var prices []Prices

	if err = result.All(context.TODO(), &prices); err != nil {
		panic(err)
	}
	return prices, err
}

func (c *Collection[Prices]) GetAllPricesByExchange(exchange string) ([]Prices, error) {
	collection := c.instance.Database(c.database).Collection(c.collectionName)

	filter := bson.D{primitive.E{Key: "exchange", Value: exchange}}

	result, err := collection.Find(context.TODO(), filter)

	if err != nil {
		panic(err)
	}

	var prices []Prices

	if err = result.All(context.TODO(), &prices); err != nil {
		panic(err)
	}

	return prices, err
}

func (c *Collection[Prices]) GetAllPricesInPeriod(from time.Time, to time.Time) ([]Prices, error) {
	collection := c.instance.Database(c.database).Collection(c.collectionName)

	filter := bson.M{"date": bson.M{
		"$gte": primitive.NewDateTimeFromTime(from),
		"$lte": primitive.NewDateTimeFromTime(to)}}

	result, err := collection.Find(context.TODO(), filter)

	var prices []Prices

	if err = result.All(context.TODO(), &prices); err != nil {
		panic(err)
	}

	return prices, err
}

func (c *Collection[Prices]) GetAllPricesInPeriodSymbol(from time.Time, to time.Time, symbol string) ([]Prices, error) {
	collection := c.instance.Database(c.database).Collection(c.collectionName)

	filter := bson.M{"date": bson.M{
		"$gte": primitive.NewDateTimeFromTime(from),
		"$lte": primitive.NewDateTimeFromTime(to)},
		"symbol": symbol,
	}

	result, err := collection.Find(context.TODO(), filter)

	var prices []Prices

	if err = result.All(context.TODO(), &prices); err != nil {
		panic(err)
	}

	return prices, err
}

func (c *Collection[Prices]) GetMostRecentPriceBySymbol(symbol string) (Prices, error) {
	collection := c.instance.Database(c.database).Collection(c.collectionName)

	filter := bson.D{primitive.E{Key: "prices.symbol", Value: symbol}}
	opts := options.FindOne().SetSort(bson.M{"$natural": -1})

	var lastRecord Prices
	err := collection.FindOne(context.TODO(), filter, opts).Decode(&lastRecord)

	return lastRecord, err
}
