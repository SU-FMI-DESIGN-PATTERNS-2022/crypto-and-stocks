package prices_repository

import (
	"context"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Collection struct {
	instance *mongo.Client
	database string
}

var collection string = "prices"

func NewCollection(client *mongo.Client, db string) *Collection {
	return &Collection{
		instance: client,
		database: db,
	}
}

func (c *Collection) insertOne(col string, doc interface{}) (*mongo.InsertOneResult, error) {

	collection := c.instance.Database(c.database).Collection(col)

	result, err := collection.InsertOne(context.TODO(), doc)
	return result, err
}

// TODO: Yet to be used
func (c *Collection) insertMany(col string, docs []interface{}) (*mongo.InsertManyResult, error) {

	collection := c.instance.Database(c.database).Collection(col)

	result, err := collection.InsertMany(context.TODO(), docs)
	return result, err
}

func (c *Collection) StoreEntry(price Prices) error {
	_, err := c.insertOne(collection, price)
	return err
}

func (c *Collection) GetAllPrices() ([]Prices, error) {
	collection := c.instance.Database(env.LoadMongoConfig().Database).Collection(collection)

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

func (c *Collection) GetAllPricesBySymbol(symbol string) ([]Prices, error) {
	collection := c.instance.Database(c.database).Collection(collection)

	filter := bson.D{{"symbol", symbol}}

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

func (c *Collection) GetAllPricesByExchange(exchange string) ([]Prices, error) {
	collection := c.instance.Database(env.LoadMongoConfig().Database).Collection(collection)

	filter := bson.D{{"exchange", exchange}}

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

func (c *Collection) GetAllPricesInPeriod(from time.Time, to time.Time) ([]Prices, error) {
	collection := c.instance.Database(env.LoadMongoConfig().Database).Collection(collection)

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

func (c *Collection) GetAllPricesInPeriodSymbol(from time.Time, to time.Time, symbol string) ([]Prices, error) {
	collection := c.instance.Database(env.LoadMongoConfig().Database).Collection(collection)

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
