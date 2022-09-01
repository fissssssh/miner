package middleware

import (
	"miner/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, e := range c.Errors {
			err := e.Err
			if customErr, ok := err.(*web.RequestError); ok {
				c.JSON(customErr.Status, gin.H{"message": customErr.Msg})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
			}
		}
	}
}
