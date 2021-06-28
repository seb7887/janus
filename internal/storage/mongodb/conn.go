package mongodb

import (
	"context"
	"fmt"
	"sync"

	"github.com/seb7887/janus/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DB = "janus"
)

var (
	mongoOnce      sync.Once
	connection     = fmt.Sprintf("mongodb://%s:%d", config.GetConfig().MongoHost, config.GetConfig().MongoPort)
	connError      error
	clientInstance *mongo.Client
)

func GetMongoClient() (*mongo.Client, error) {
	// perform connection creation operation only once
	mongoOnce.Do(func() {
		options := options.Client().ApplyURI(connection)
		client, err := mongo.Connect(context.TODO(), options)
		if err != nil {
			connError = err
		}
		clientInstance = client
	})

	return clientInstance, connError
}
