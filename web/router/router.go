package router

import (
	"fmt"
	"miner/game"
	"miner/web/errors"
	"miner/web/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	r.Use(middleware.ErrorHandler())

	r.StaticFile("/", "./index.html")

	apiGourp := r.Group("api")
	{
		apiGourp.GET("rank", func(c *gin.Context) {
			winners, err := game.GetRank(false)
			if err != nil {
				c.Error(err)
				return
			}
			c.JSON(200, winners)
		})

		apiGourp.POST("newgame", func(c *gin.Context) {
			name := c.PostForm("name")
			if name == "" {
				c.Error(errors.NewRequestError(400, "玩家名称不可为空"))
				return
			}
			min, max := 1, 100000
			g, err := game.New(name, min, max)
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

		apiGourp.POST("games/:id/dig/:pos", func(c *gin.Context) {
			id := c.Param("id")
			posStr := c.Param("pos")
			pos, err := strconv.Atoi(posStr)
			if err != nil {
				c.Error(errors.NewRequestError(400, "参数错误"))
				return
			}
			g, err := game.Get(id)
			if err != nil {
				c.Error(err)
				return
			}
			if g == nil {
				c.Error(errors.NewRequestError(404, "游戏不存在"))
				return
			}
			if pos <= g.Min || pos >= g.Max {
				c.Error(errors.NewRequestError(400, fmt.Sprintf("挖掘位置必须在%v~%v之间", g.Min, g.Max)))
				return
			}
			status, win, err := g.Dig(id, pos)
			if err != nil {
				c.Error(err)
				return
			}
			if status == game.StatusFinish {
				game.Del(id)
				if win {
					game.AddWinner(g.Name, g.Count)
				}
				c.JSON(200, gin.H{
					"status": status,
					"win":    win,
					"result": g.Result,
					"count":  g.Count,
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
}
