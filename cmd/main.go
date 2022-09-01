package main

import (
	"miner/game"
	"miner/web/router"

	"github.com/gin-gonic/gin"
)

func main() {
	game.Init()
	r := gin.Default()
	router.Init(r)
	r.Run()
}
