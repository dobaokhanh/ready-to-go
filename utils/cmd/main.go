package main

import (
	"context"
	"fmt"
	"utils/db"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	configMG := db.MongoDBConfig{
		URI:          "mongodb://testmongo:testmongo@127.0.0.1:27017/",
		DBName:       "demo",
		PoolSize:     200,
		MaxIdleConns: 10000,
	}

	mongoDB, err := db.NewMongoDBImpl(configMG)
	if err != nil {
		panic(err)
	}
	defer mongoDB.Disconnect()
	dbMG := mongoDB.GetConnection().(*mongo.Database)
	testCollection := dbMG.Collection("test")
	result, err := testCollection.InsertOne(context.Background(), bson.D{
		{Key: "test2", Value: "1"},
		{Key: "test2", Value: "2"},
	})

	fmt.Println(result)

	configRDB := db.RedisConfig{
		URL:      "localhost:6379",
		Password: "redis_master",
		DB:       0,
	}

	rdb, err := db.NewRedisImpl(configRDB)

	if err != nil {
		panic(err)
	}
	defer rdb.Disconnect()

	dbRDB := rdb.GetConnection().(*redis.Client)
	err = dbRDB.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := dbRDB.Get(context.Background(), "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
}
