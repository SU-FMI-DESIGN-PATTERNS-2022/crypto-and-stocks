package repositories

import (
	"context"
	"fmt"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strconv"
)

// Connect Connects The Go app with a MongoDb database
func Connect(mongoConfig env.MongoConfig) (*mongo.Client, error) {

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoConfig.Driver+"://"+mongoConfig.Host+":"+strconv.FormatInt(int64(mongoConfig.Port), 10)))
	return client, err
}

// Ping This method used to ping the mongoDB, return error if any.
func Ping(client *mongo.Client, ctx context.Context) error {

	// mongo.Client has Ping to ping mongoDB, deadline of
	// the Ping method will be determined by cxt
	// Ping method return error if any occurred, then
	// the error can be handled.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("Connected to MongoDb successfully")
	return nil
}

// Close This method closes mongoDB connection and cancel context.
func Close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func main() {

	mongoConfig := env.LoadMongoConfig()
	client, connectErr := Connect(mongoConfig)

	if connectErr != nil {
		panic(connectErr)
	}
	defer func() {
		if connectErr = client.Disconnect(context.TODO()); connectErr != nil {
			panic(connectErr)
		}
	}()
	// Checking whether the connection was successful
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged MongoDB.")
}
