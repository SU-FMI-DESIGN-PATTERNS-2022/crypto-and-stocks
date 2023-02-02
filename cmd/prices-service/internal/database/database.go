package database

import (
	"context"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection struct {
	Instance       *mongo.Client
	Database       string
	CollectionName string
}

func (c *Collection) InsertOne(col string, doc interface{}) (*mongo.InsertOneResult, error) {
	collection := c.Instance.Database(c.Database).Collection(col)

	return collection.InsertOne(context.TODO(), doc)
}

func (c *Collection) InsertMany(col string, docs []interface{}) (*mongo.InsertManyResult, error) {
	collection := c.Instance.Database(c.Database).Collection(col)

	return collection.InsertMany(context.TODO(), docs)
}
