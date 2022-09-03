package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

var buckets map[string]*ratelimit.Bucket = make(map[string]*ratelimit.Bucket)

func RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.Request.RemoteAddr
		if bucket, exists := buckets[ip]; exists {
			if bucket.TakeAvailable(1) < 1 {
				c.String(http.StatusForbidden, "rate limit...")
				c.Abort()
				return
			}
		} else {
			buckets[ip] = ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
		}
		c.Next()
	}
}
