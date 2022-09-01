package router

import (
	"fmt"
	"miner/game"
	"miner/web"
	"miner/web/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	r.Use(middleware.ErrorHandler())

	r.StaticFile("/", "./index.html")

	r.POST("games/new", func(c *gin.Context) {
		min, max := 1, 100000
		g, err := game.New(min, max)
		if err != nil {
			c.Status(500)
			return
		}
		c.JSON(200, gin.H{
			"id":  g.Id,
			"min": g.Min,
			"max": g.Max,
		})
	})

	r.POST("games/:id/dig/:pos", func(c *gin.Context) {
		id := c.Param("id")
		posStr := c.Param("pos")
		pos, err := strconv.Atoi(posStr)
		if err != nil {
			c.Error(web.NewRequestError(400, "参数错误"))
			return
		}
		g, err := game.Get(id)
		if err != nil {
			c.Error(err)
			return
		}
		if g == nil {
			c.Error(web.NewRequestError(404, "游戏不存在"))
			return
		}
		if pos <= g.Min || pos >= g.Max {
			c.Error(web.NewRequestError(400, fmt.Sprintf("挖掘位置必须在%v~%v之间", g.Min, g.Max)))
			return
		}
		status, win, err := g.Dig(id, pos)
		if err != nil {
			c.Error(err)
			return
		}
		if status == game.StatusFinish {
			game.Del(id)
			c.JSON(200, gin.H{
				"status": status,
				"win":    win,
			})
		} else {
			c.JSON(200, gin.H{
				"status": status,
				"min":    g.Min,
				"max":    g.Max,
			})
		}
	})
}
