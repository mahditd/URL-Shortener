package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {

		ip := c.ClientIP()

		if !limiter.Allow(ip) {
			c.Header("Retry-After", "60")

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
