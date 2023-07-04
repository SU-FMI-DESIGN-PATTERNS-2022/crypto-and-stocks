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

func (collection *Collection[Prices]) insertOne(col string, doc interface{}) (*mongo.InsertOneResult, error) {

	dbCollection := collection.instance.Database(collection.database).Collection(col)

	result, err := dbCollection.InsertOne(context.TODO(), doc)
	return result, err
}

// TODO: Yet to be used
func (collection *Collection[Prices]) insertMany(col string, docs []interface{}) (*mongo.InsertManyResult, error) {

	dbCollection := collection.instance.Database(collection.database).Collection(col)

	result, err := dbCollection.InsertMany(context.TODO(), docs)
	return result, err
}

func (collection *Collection[Prices]) StoreEntry(price Prices) error {
	_, err := collection.insertOne(collection.collectionName, price)
	return err
}

func (collection *Collection[Prices]) GetAllPrices() ([]Prices, error) {
	dbCollection := collection.instance.Database(collection.database).Collection(collection.collectionName)

	result, err := dbCollection.Find(context.TODO(), bson.D{})

	if err != nil {
		return nil, err
	}
	var prices []Prices

	if err = result.All(context.TODO(), &prices); err != nil {
		return nil, err
	}
	return prices, err
}

func (collection *Collection[Prices]) GetAllPricesBySymbol(symbol string) ([]Prices, error) {
	dbCollection := collection.instance.Database(collection.database).Collection(collection.collectionName)

	filter := bson.D{primitive.E{Key: "symbol", Value: symbol}}

	result, err := dbCollection.Find(context.TODO(), filter)

	if err != nil {
		return nil, err
	}
	var prices []Prices

	if err = result.All(context.TODO(), &prices); err != nil {
		return nil, err
	}
	return prices, err
}

func (collection *Collection[Prices]) GetAllPricesByExchange(exchange string) ([]Prices, error) {
	dbCollection := collection.instance.Database(collection.database).Collection(collection.collectionName)

	filter := bson.D{primitive.E{Key: "exchange", Value: exchange}}

	result, err := dbCollection.Find(context.TODO(), filter)

	if err != nil {
		return nil, err
	}

	var prices []Prices

	if err = result.All(context.TODO(), &prices); err != nil {
		return nil, err
	}

	return prices, err
}

func (collection *Collection[Prices]) GetAllPricesInPeriod(from time.Time, to time.Time) ([]Prices, error) {
	dbCollection := collection.instance.Database(collection.database).Collection(collection.collectionName)

	filter := bson.M{"date": bson.M{
		"$gte": primitive.NewDateTimeFromTime(from),
		"$lte": primitive.NewDateTimeFromTime(to)}}

	result, err := dbCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var prices []Prices

	if err = result.All(context.TODO(), &prices); err != nil {
		return nil, err
	}

	return prices, err
}

func (collection *Collection[Prices]) GetAllPricesInPeriodSymbol(from time.Time, to time.Time, symbol string) ([]Prices, error) {
	dbCollection := collection.instance.Database(collection.database).Collection(collection.collectionName)

	filter := bson.M{"date": bson.M{
		"$gte": primitive.NewDateTimeFromTime(from),
		"$lte": primitive.NewDateTimeFromTime(to)},
		"symbol": symbol,
	}

	result, err := dbCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var prices []Prices

	if err = result.All(context.TODO(), &prices); err != nil {
		return nil, err
	}

	return prices, err
}

func (collection *Collection[Prices]) GetMostRecentPriceBySymbol(symbol string) (Prices, error) {
	dbCollection := collection.instance.Database(collection.database).Collection(collection.collectionName)

	filter := bson.D{primitive.E{Key: "prices.symbol", Value: symbol}}
	opts := options.FindOne().SetSort(bson.M{"$natural": -1})

	var lastRecord Prices
	err := dbCollection.FindOne(context.TODO(), filter, opts).Decode(&lastRecord)

	return lastRecord, err
}
