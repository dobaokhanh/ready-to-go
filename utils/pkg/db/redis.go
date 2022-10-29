package db

import (
	"context"

	"github.com/go-redis/redis/v9"
)

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

type RedisImpl struct {
	db     *redis.Client
	config RedisConfig
}

func NewRedisImpl(config RedisConfig) (Database, error) {
	rdb := &RedisImpl{
		config: config,
	}

	if err := rdb.Connect(); err != nil {
		return nil, err
	}

	return rdb, nil
}

func (redisDB *RedisImpl) Connect() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisDB.config.URL,
		Password: redisDB.config.Password,
		DB:       redisDB.config.DB,
	})

	cmd := rdb.Ping(context.Background())
	redisDB.db = rdb
	return cmd.Err()
}

func (redisDB *RedisImpl) Disconnect() error {
	return redisDB.db.Close()
}

func (redisDB *RedisImpl) GetConnection() interface{} {
	return redisDB.db
}
