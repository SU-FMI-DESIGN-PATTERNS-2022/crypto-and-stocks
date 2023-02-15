package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Connection int

const (
	Local Connection = iota
	Remote
)

func Connect(config env.MongoConfig, connection Connection) (*mongo.Client, error) {
	var mongoconn string

	switch connection {
	case Local:
		mongoconn = fmt.Sprintf("%s://localhost:%s", config.LocalDriver, config.Port)
	case Remote:
		mongoconn = fmt.Sprintf("%s://%s:%s@%s/?%s", config.RemoteDriver, config.User, config.Password, config.Host, config.Options)
	default:
		return nil, errors.New("Unrecognized connection type")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoconn))

	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}
