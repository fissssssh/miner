package web

import (
	"miner/web/router"

	"github.com/gin-gonic/gin"
)

func Init() {
	r := gin.Default()
	router.Init(r)
	r.Run()
}
