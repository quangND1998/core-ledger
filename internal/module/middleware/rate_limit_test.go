package middleware

import (
	// "fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	// "time"

	// "github.com/didip/tollbooth/v7"
	// "github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

// func RateLimitMiddleware() gin.HandlerFunc {
// 	limit := tollbooth.NewLimiter(60, nil)
// 	limit.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})
// 	limit.SetBurst(5)
// 	limit.SetTokenBucketExpirationTTL(time.Minute)
// 	limit.SetMessage(`{"error": "Too many requests, please try again later."}`)
// 	limit.SetMessageContentType("application/json")
// 	limit.SetStatusCode(429)

// 	fmt.Println("✅ RateLimitMiddleware initialized") // Debug

// 	return tollbooth_gin.LimitHandler(limit)
// }

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RateLimitMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Gửi 65 request
	for i := 0; i < 20; i++ {
		r.ServeHTTP(w, req)
		if i < 60 && w.Code != 200 {
			t.Errorf("expected 200, got %d  time %d", w.Code, i)
		}
		if i >= 60 && w.Code != 429 {
			t.Errorf("expected 429 after limit exceeded, got %d time %d", w.Code, i)
		}
	}
}
