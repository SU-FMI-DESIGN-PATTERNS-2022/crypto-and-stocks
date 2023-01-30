package mongo

import (
	"context"
	"fmt"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strconv"
	"time"
)

// Connect Connects The Go app with a MongoDb database
func Connect(mongoConfig env.MongoConfig) (*mongo.Client, context.CancelFunc, error) {

	_, cancel := context.WithTimeout(context.Background(),
		30*time.Second)

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(
		mongoConfig.Driver+"://"+mongoConfig.User+":"+mongoConfig.Password+"@"+mongoConfig.Host+":"+strconv.FormatInt(int64(mongoConfig.Port), 10)))
	return client, cancel, err
}

// Ping This method used to ping the mongoDB, return error if any.
func Ping(client *mongo.Client, ctx context.Context) error {

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("Connected to MongoDb successfully")
	return nil
}

// Close This method closes mongoDB connection and cancel context.
func Close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {
	defer cancel()

	defer func() {

		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}
