package main

import (
	"context"
	"net/http"
	"time"
	db "utils/db"
	id_generator "utils/id_generator"
	rate_limiter "utils/rate_limiter"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	router = gin.Default()
)

// const (
// 	// DefaultPrefix url prefix of pprof
// 	DefaultPrefix = "/debug/pprof"
// )

// func getPrefix(prefixOptions ...string) string {
// 	prefix := DefaultPrefix
// 	if len(prefixOptions) > 0 {
// 		prefix = prefixOptions[0]
// 	}
// 	return prefix
// }

// // Register the standard HandlerFuncs from the net/http/pprof package with
// // the provided gin.Engine. prefixOptions is a optional. If not prefixOptions,
// // the default path prefix is used, otherwise first prefixOptions will be path prefix.
// func Register(r *gin.Engine, prefixOptions ...string) {
// 	RouteRegister(&(r.RouterGroup), prefixOptions...)
// }

// // RouteRegister the standard HandlerFuncs from the net/http/pprof package with
// // the provided gin.GrouterGroup. prefixOptions is a optional. If not prefixOptions,
// // the default path prefix is used, otherwise first prefixOptions will be path prefix.
// func RouteRegister(rg *gin.RouterGroup, prefixOptions ...string) {
// 	prefix := getPrefix(prefixOptions...)

// 	prefixRouter := rg.Group(prefix)
// 	{
// 		prefixRouter.GET("/", gin.WrapF(pprof.Index))
// 		prefixRouter.GET("/cmdline", gin.WrapF(pprof.Cmdline))
// 		prefixRouter.GET("/profile", gin.WrapF(pprof.Profile))
// 		prefixRouter.POST("/symbol", gin.WrapF(pprof.Symbol))
// 		prefixRouter.GET("/symbol", gin.WrapF(pprof.Symbol))
// 		prefixRouter.GET("/trace", gin.WrapF(pprof.Trace))
// 		prefixRouter.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
// 		prefixRouter.GET("/block", gin.WrapH(pprof.Handler("block")))
// 		prefixRouter.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
// 		prefixRouter.GET("/heap", gin.WrapH(pprof.Handler("heap")))
// 		prefixRouter.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
// 		prefixRouter.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
// 	}
// }

func main() {
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

	rateLimiterConfig := rate_limiter.RateLimiterConfig{
		Strategy:    rate_limiter.NewSortedSetCounterStrategy(dbRDB, time.Now),
		Expiration:  time.Second,
		MaxRequests: 100000,
	}

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

	gin.SetMode("debug")

	router.Use(rate_limiter.UseRateLimiter(rateLimiterConfig))

	router.GET("/id", func(c *gin.Context) {
		id, _ := id_generator.NextID()
		testCollection.InsertOne(context.Background(), bson.D{
			{Key: "_id", Value: id},
		})
		c.JSON(http.StatusOK, gin.H{
			"message": id,
		})
	})
	router.Run()
}
