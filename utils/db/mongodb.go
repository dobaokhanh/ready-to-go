package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConfig struct {
	URI          string
	DBName       string
	PoolSize     int
	MaxIdleConns int
}

type MongoDBImpl struct {
	db     *mongo.Database
	config MongoDBConfig
}

func NewMongoDBImpl(config MongoDBConfig) (Database, error) {
	mongoDB := &MongoDBImpl{
		config: config,
	}

	if err := mongoDB.Connect(); err != nil {
		return nil, err
	}

	return mongoDB, nil
}

func (mongoDB *MongoDBImpl) Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDB.buildURI()))

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)

	if err != nil {
		return err
	}

	mongoDB.db = client.Database(mongoDB.config.DBName)

	return nil
}

func (mongoDB *MongoDBImpl) Disconnect() error {
	return mongoDB.db.Client().Disconnect(context.Background())
}

func (mongoDB *MongoDBImpl) GetConnection() interface{} {
	return mongoDB.db
}

func (mongoDB *MongoDBImpl) buildURI() string {
	return fmt.Sprintf("%s?maxIdleTimeMS=%d&maxPoolSize=%d",
		mongoDB.config.URI, mongoDB.config.MaxIdleConns, mongoDB.config.PoolSize)
}
