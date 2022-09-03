package game

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"miner/db"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type GameStatus uint8

const (
	StatusFinish GameStatus = iota
	StatusContinue
)

type Game struct {
	Id     string `redis:"-"`
	Result int    `redis:"result"`
	Count  int    `redis:"count"`
	Min    int    `redis:"min"`
	Max    int    `redis:"max"`
	Name   string `redis:"name"`
}

func (g *Game) persistence() error {
	rdb := db.GetRedis()
	gameKey := fmt.Sprintf("Game:%s", g.Id)
	ctx := context.Background()

	if _, err := rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		p.HSet(ctx, gameKey, "name", g.Name)
		p.HSet(ctx, gameKey, "result", g.Result)
		p.HSet(ctx, gameKey, "count", g.Count)
		p.HSet(ctx, gameKey, "max", g.Max)
		p.HSet(ctx, gameKey, "min", g.Min)
		return nil
	}); err != nil {
		return err
	}

	rdb.Expire(ctx, gameKey, 10*time.Minute)

	return nil
}

func (g *Game) Dig(id string, pos int) (GameStatus, bool, error) {
	g.Count++
	if pos == g.Result {
		return StatusFinish, true, nil
	}

	if pos < g.Result {
		g.Min = pos
	} else {
		g.Max = pos
	}

	if g.Max-g.Min <= 2 {
		return StatusFinish, false, nil
	}

	err := g.persistence()
	if err != nil {
		return StatusContinue, false, err
	}

	return StatusContinue, false, nil
}

func New(name string, min int, max int) (*Game, error) {
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	rand.Seed(time.Now().UnixMilli())
	g := Game{
		Name:   name,
		Id:     id,
		Min:    min,
		Max:    max,
		Result: rand.Intn(max-1) + 1,
	}

	err := g.persistence()
	if err != nil {
		return nil, err
	}

	log.Printf("start new game %s", id)
	return &g, nil
}

func Get(id string) (*Game, error) {
	rdb := db.GetRedis()
	ctx := context.Background()
	gameKey := fmt.Sprintf("Game:%s", id)

	exists, err := rdb.Exists(ctx, gameKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	if exists != 1 {
		return nil, nil
	}

	var game Game
	if err := rdb.HGetAll(ctx, gameKey).Scan(&game); err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	game.Id = id
	return &game, nil
}

func Del(id string) {
	rdb := db.GetRedis()
	ctx := context.Background()
	gameKey := fmt.Sprintf("Game:%s", id)
	rdb.Del(ctx, gameKey)
}
