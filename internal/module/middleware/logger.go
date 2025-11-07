package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
	"core-ledger/pkg/logging"
)

const (
	ok          = "OK"
	clientError = "CLIENT_ERROR"
	serverError = "SERVER_ERROR"
)

const (
	HeaderRequestID = "X-Request-SystemPaymentID"

	ServiceCode    = "_code"
	ServicePayload = "_payload"
)

// LogRequest is a middleware that logs HTTP requests with detailed fields suitable for Grafana/Prometheus statistics.
// This logging approach is quite common in modern Go web services, especially those using Gin or similar frameworks.
// It provides structured logging with useful fields for observability and monitoring, which is popular for production systems.
// Many teams use similar patterns to enrich logs for tools like Grafana, Prometheus, or ELK/EFK stacks.
func LogRequest(c *gin.Context) {
	// This pattern of logging request/response metadata is widely adopted in Go web APIs.
	method := c.Request.Method
	url := c.Request.URL
	requestID := c.GetHeader(HeaderRequestID)
	start := time.Now()

	if requestID == "" {
		requestID = uuid.NewString()
	}

	ctx := c.Request.Context()
	logger := logging.
		From(ctx).
		Named("http").
		With("request_id", requestID)

	// Attach logger with request_id to context for downstream handlers
	c.Request = c.Request.WithContext(logging.WithLogger(ctx, logger))
	c.Next()

	status := c.Writer.Status()
	latency := time.Since(start)

	// Enrich log with common HTTP request/response fields
	logger = logger.
		With("timestamp", start.UTC().Format(time.RFC3339Nano)).
		With("latency_ms", latency.Milliseconds()).
		With("method", method).
		With("path", url.Path).
		With("full_path", url.String()).
		With("query", url.RawQuery).
		With("status", status).
		With("remote_ip", c.ClientIP()).
		With("user_agent", c.Request.UserAgent()).
		With("referer", c.Request.Referer()).
		With("request_id", requestID)

	if code, ok := c.Get(ServiceCode); ok {
		logger = logger.With("service_code", code)
	}
	if payload, ok := c.Get(ServicePayload); ok {
		logger = logger.With("service_payload", payload)
	}

	// Log request and response sizes, which is also a popular practice
	logger = logger.
		With("request_size", c.Request.ContentLength).
		With("response_size", c.Writer.Size())

	// This switch for log level based on status code is also a common pattern
	switch {
	case status >= 400 && status < 500:
		logger.Warn(clientError)
	case status >= 500 && status <= 599:
		logger.Error(serverError)
	case status >= 200 && status < 400:
		logger.Info(ok)
	default:
		logger.Info("_unknown")
	}
}
