package game

import (
	"context"
	"encoding/json"
	"log"
	"miner/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Winner struct {
	Name  string `json:"name" redis:"name"`
	Count int    `json:"count" redis:"count"`
	Time  int64  `json:"time" redis:"time"`
}

const rankKey = "rank"
const rankCount = 10

func AddWinner(name string, count int) {
	rdb := db.GetRedis()
	ctx := context.Background()

	winners, _ := GetRank(true)
	// 若排行榜未满，则删除排行榜缓存
	if len(*winners) < rankCount {
		_ = rdb.Del(ctx, rankKey)

	} else {
		// 判断是否进入排行榜 如果进入则删除排行缓存
		for _, winner := range *winners {
			if count < winner.Count {
				_ = rdb.Del(ctx, rankKey)
				break
			}
		}
	}

	mongodb := db.GetMongoDb()
	collection := mongodb.Collection("winners")

	winner := Winner{
		Name:  name,
		Count: count,
		Time:  time.Now().UnixMilli(),
	}

	collection.InsertOne(ctx, winner)
}

func GetRank(cacheOnly bool) (*[]*Winner, error) {
	ctx := context.Background()

	// 从缓存读取
	var winners = []*Winner{}
	rdb := db.GetRedis()
	if exists, _ := rdb.Exists(ctx, rankKey).Result(); exists == 1 {
		if rankJson, _ := rdb.Get(ctx, rankKey).Result(); rankJson != "" {
			err := json.Unmarshal([]byte(rankJson), &winners)
			if err == nil {
				log.Println("get rank cache hit!")
				return &winners, nil
			}
		}
	}

	if cacheOnly {
		return &winners, nil
	}

	mongodb := db.GetMongoDb()
	collection := mongodb.Collection("winners")

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "count", Value: 1},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{{Key: "name", Value: "$name"}}},
				{Key: "record", Value: bson.D{{Key: "$first", Value: "$$ROOT"}}},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "name", Value: "$record.name"},
				{Key: "count", Value: "$record.count"},
				{Key: "time", Value: "$record.time"},
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "count", Value: 1},
				{Key: "time", Value: 1},
			}},
		},
		bson.D{
			{Key: "$limit", Value: rankCount},
		},
	}

	cur, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &winners); err != nil {
		return nil, err
	}

	// 查询到结果写入缓存
	rankJson, err := json.Marshal(&winners)
	if err != nil {
		return nil, err
	}
	rdb.Set(ctx, rankKey, string(rankJson), 1*time.Hour)

	return &winners, nil
}
