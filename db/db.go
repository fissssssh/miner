package db

import (
	"context"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var rdb *redis.Client
var once sync.Once

func Init() {
	once.Do(func() {
		ctx := context.Background()
		redisConnStr := viper.GetString("connectionStrings.redis")
		rdb = redis.NewClient(&redis.Options{
			Addr:     redisConnStr,
			Password: "",
			DB:       0,
		})

		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			panic(err)
		}

	})
}

func GetRedis() *redis.Client {
	return rdb
}

func GetMongoDb() *mongo.Database {
	ctx := context.Background()
	mongodbConnStr := viper.GetString("connectionStrings.mongodb")
	mongodb, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbConnStr))
	if err != nil {
		log.Fatal(err)
	}

	err = mongodb.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	return mongodb.Database("miner")
}
