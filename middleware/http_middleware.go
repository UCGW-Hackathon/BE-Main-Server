package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimitHeaders(limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		reset := time.Now().Add(time.Minute).Unix()
		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(limit-1))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(reset, 10))
		c.Next()
	}
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Next()
	}
}
