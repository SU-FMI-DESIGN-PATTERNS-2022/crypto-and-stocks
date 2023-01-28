package prices_repository

import (
	"context"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDb struct {
	Client     *mongo.Client
	Context    context.Context
	CancelFunc context.CancelFunc
}
type Database struct {
	instance *MongoDb
}

func NewDatabase(db *MongoDb) *Database {
	return &Database{
		instance: db,
	}
}

func insertOne(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error) {

	collection := client.Database(dataBase).Collection(col)

	result, err := collection.InsertOne(ctx, doc)
	return result, err
}

func insertMany(client *mongo.Client, ctx context.Context, dataBase, col string, docs []interface{}) (*mongo.InsertManyResult, error) {

	collection := client.Database(dataBase).Collection(col)

	result, err := collection.InsertMany(ctx, docs)
	return result, err
}
