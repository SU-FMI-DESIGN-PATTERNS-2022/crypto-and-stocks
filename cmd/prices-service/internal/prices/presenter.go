package prices

import (
	"context"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/env"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/repositories/prices_repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var collection string = "prices"

type PricesRepository interface {
	StoreEntry(prices prices_repository.Prices, client *mongo.Client) error
	GetAllPrices(client *mongo.Client) ([]prices_repository.Prices, error)
	GetAllPricesBySymbol(symbol string, client *mongo.Client) ([]prices_repository.Prices, error)
	GetAllPricesByExchange(exchange string, client *mongo.Client) ([]prices_repository.Prices, error)
	GetAllPricesInPeriod(from time.Time, to time.Time, client *mongo.Client) ([]prices_repository.Prices, error)
	GetAllPricesInPeriodSymbol(from time.Time, to time.Time, symbol string, client *mongo.Client) ([]prices_repository.Prices, error)
}

type Presenter struct {
	pricesRepository PricesRepository
}

func NewPresenter(repository PricesRepository) Presenter {
	return Presenter{
		pricesRepository: repository,
	}
}

func StoreEntry(price prices_repository.Prices, client *mongo.Client) error {
	collection := client.Database(env.LoadMongoConfig().Database).Collection(collection)

	_, err := collection.InsertOne(context.TODO(), price)
	return err
}

func GetAllPrices(client *mongo.Client) ([]prices_repository.Prices, error) {
	collection := client.Database(env.LoadMongoConfig().Database).Collection(collection)

	result, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		panic(err)
	}
	var prices []prices_repository.Prices

	if err = result.All(context.TODO(), &prices); err != nil {
		panic(err)
	}
	return prices, err
}

func GetAllPricesBySymbol(symbol string, client mongo.Client) ([]prices_repository.Prices, error) {
	collection := client.Database(env.LoadMongoConfig().Database).Collection(collection)

	filter := bson.D{{"symbol", symbol}}

	result, err := collection.Find(context.TODO(), filter)

	if err != nil {
		panic(err)
	}
	var prices []prices_repository.Prices

	if err = result.All(context.TODO(), &prices); err != nil {
		panic(err)
	}
	return prices, err
}

func GetAllPricesByExchange(exchange string, client *mongo.Client) ([]prices_repository.Prices, error) {
	collection := client.Database(env.LoadMongoConfig().Database).Collection(collection)

	filter := bson.D{{"exchange", exchange}}

	result, err := collection.Find(context.TODO(), filter)

	if err != nil {
		panic(err)
	}

	var prices []prices_repository.Prices

	if err = result.All(context.TODO(), &prices); err != nil {
		panic(err)
	}

	return prices, err
}

func GetAllPricesInPeriod(from time.Time, to time.Time, client *mongo.Client) ([]prices_repository.Prices, error) {
	collection := client.Database(env.LoadMongoConfig().Database).Collection(collection)

	filter := bson.M{"date": bson.M{
		"$gte": primitive.NewDateTimeFromTime(from),
		"$lte": primitive.NewDateTimeFromTime(to)}}

	result, err := collection.Find(context.TODO(), filter)

	var prices []prices_repository.Prices

	if err = result.All(context.TODO(), &prices); err != nil {
		panic(err)
	}

	return prices, err
}

func GetAllPricesInPeriodSymbol(from time.Time, to time.Time, symbol string, client *mongo.Client) ([]prices_repository.Prices, error) {
	collection := client.Database(env.LoadMongoConfig().Database).Collection(collection)

	filter := bson.M{"date": bson.M{
		"$gte": primitive.NewDateTimeFromTime(from),
		"$lte": primitive.NewDateTimeFromTime(to)},
		"symbol": symbol,
	}

	result, err := collection.Find(context.TODO(), filter)

	var prices []prices_repository.Prices

	if err = result.All(context.TODO(), &prices); err != nil {
		panic(err)
	}

	return prices, err
}
