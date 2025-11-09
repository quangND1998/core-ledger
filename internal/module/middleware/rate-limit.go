package middleware

import (
	"fmt"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware() gin.HandlerFunc {
	limit := tollbooth.NewLimiter(60, nil)
	limit.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})
	limit.SetBurst(5)
	limit.SetTokenBucketExpirationTTL(time.Minute)
	limit.SetMessage(`{"error": "Too many requests, please try again later."}`)
	limit.SetMessageContentType("application/json")
	limit.SetStatusCode(429)

	fmt.Println("✅ RateLimitMiddleware initialized") // Debug

	return tollbooth_gin.LimitHandler(limit)
}
