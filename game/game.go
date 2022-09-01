package game

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var rdb *redis.Client

type GameStatus uint8

const (
	StatusFinish GameStatus = iota
	StatusContinue
)

type Game struct {
	Id     string `redis:"-"`
	Result int    `redis:"result"`
	Min    int    `redis:"min"`
	Max    int    `redis:"max"`
}

func (g *Game) persistence() error {
	gameKey := fmt.Sprintf("Game:%s", g.Id)
	ctx := context.Background()

	if _, err := rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		p.HSet(ctx, gameKey, "result", g.Result)
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

var once sync.Once

func Init() {
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "",
			DB:       0,
		})
		_, err := rdb.Ping(context.Background()).Result()
		if err != nil {
			panic(err)
		}
	})
}

func New(min int, max int) (*Game, error) {
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	g := Game{Id: id, Min: min, Max: max, Result: rand.Intn(max-1) + 1}

	err := g.persistence()
	if err != nil {
		return nil, err
	}

	log.Printf("start new game %s", id)
	return &g, nil
}

func Get(id string) (*Game, error) {
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
	ctx := context.Background()
	gameKey := fmt.Sprintf("Game:%s", id)
	rdb.Del(ctx, gameKey)
}
