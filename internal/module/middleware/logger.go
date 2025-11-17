package middleware

import (
	"core-ledger/pkg/logging"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
// Highly optimized for performance using zap.Logger with batched fields.
func LogRequest(c *gin.Context) {
	start := time.Now()
	method := c.Request.Method
	path := c.Request.URL.Path
	requestID := c.GetHeader(HeaderRequestID)

	// Generate UUID only if needed
	if requestID == "" {
		requestID = uuid.NewString()
	}

	ctx := c.Request.Context()
	sugaredLogger := logging.From(ctx).Named("http")

	// Attach logger to context (minimal overhead)
	c.Request = c.Request.WithContext(logging.WithLogger(ctx, sugaredLogger.With("request_id", requestID)))

	c.Next()

	status := c.Writer.Status()
	latency := time.Since(start)

	// Get zap.Logger (not SugaredLogger) for better performance
	baseLogger := sugaredLogger.Desugar()

	// Check log level BEFORE any expensive operations
	core := baseLogger.Core()
	var shouldLog bool
	var logLevel zapcore.Level
	var msg string

	switch {
	case status >= 500:
		shouldLog = core.Enabled(zap.ErrorLevel)
		logLevel = zap.ErrorLevel
		msg = serverError
	case status >= 400:
		shouldLog = core.Enabled(zap.WarnLevel)
		logLevel = zap.WarnLevel
		msg = clientError
	case status >= 200:
		shouldLog = core.Enabled(zap.InfoLevel)
		logLevel = zap.InfoLevel
		msg = ok
	default:
		shouldLog = core.Enabled(zap.InfoLevel)
		logLevel = zap.InfoLevel
		msg = "_unknown"
	}

	// Early return if logging is disabled for this level
	if !shouldLog {
		return
	}

	// Build all fields in one slice (minimal allocations)
	fields := make([]zap.Field, 0, 12) // Pre-allocate with estimated capacity
	fields = append(fields,
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status", status),
		zap.Int64("latency_ms", latency.Milliseconds()),
	)

	// Add optional fields only if they exist
	if requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}
	if ip := c.ClientIP(); ip != "" {
		fields = append(fields, zap.String("remote_ip", ip))
	}
	if query := c.Request.URL.RawQuery; query != "" {
		fields = append(fields, zap.String("query", query))
	}

	// Only add expensive fields for errors
	if status >= 400 {
		if ua := c.Request.UserAgent(); ua != "" {
			fields = append(fields, zap.String("user_agent", ua))
		}
		if ref := c.Request.Referer(); ref != "" {
			fields = append(fields, zap.String("referer", ref))
		}
	}

	// Only log sizes for non-zero values
	if c.Request.ContentLength > 0 {
		fields = append(fields, zap.Int64("request_size", c.Request.ContentLength))
	}
	if size := c.Writer.Size(); size > 0 {
		fields = append(fields, zap.Int("response_size", size))
	}

	// Add service-specific fields if present
	if code, ok := c.Get(ServiceCode); ok {
		fields = append(fields, zap.Any("service_code", code))
	}
	if payload, ok := c.Get(ServicePayload); ok {
		fields = append(fields, zap.Any("service_payload", payload))
	}

	// Build logger with all fields in ONE call (single allocation)
	enrichedLogger := baseLogger.With(fields...)

	// Log with batched fields (single allocation)
	switch logLevel {
	case zap.ErrorLevel:
		enrichedLogger.Error(msg)
	case zap.WarnLevel:
		enrichedLogger.Warn(msg)
	default:
		enrichedLogger.Info(msg)
	}
}
