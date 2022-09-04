package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

var bucketsWrapper = struct {
	sync.RWMutex
	Buckets map[string]*ratelimit.Bucket
}{Buckets: make(map[string]*ratelimit.Bucket)}

func RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.Request.RemoteAddr
		bucketsWrapper.RLock()
		bucket, exists := bucketsWrapper.Buckets[ip]
		bucketsWrapper.RUnlock()
		if exists {
			if bucket.TakeAvailable(1) < 1 {
				c.String(http.StatusForbidden, "rate limit...")
				c.Abort()
				return
			}
		} else {
			bucketsWrapper.Lock()
			bucketsWrapper.Buckets[ip] = ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
			bucketsWrapper.Unlock()
		}
		c.Next()
	}
}
