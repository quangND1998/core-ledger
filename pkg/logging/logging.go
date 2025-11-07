// Package logging sets up and configures logging.
package logging

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	//gin
	"github.com/gin-gonic/gin"
)

var (
	defaultLogLevel = "INFO"       // Default log level if not set in environment
	defaultLogMode  = "production" // Default log mode if not set in environment
)

// contextKey is a private string type to prevent collisions in the context map.
type contextKey string

const loggerKey = contextKey("logger")

var (
	defaultLogger     *zap.SugaredLogger
	defaultLoggerOnce sync.Once
)

// NewLogger creates a new logger with the given configuration.
func NewLogger(level string, development bool) *zap.SugaredLogger {
	var config *zap.Config
	if development {
		config = &zap.Config{
			Level:            zap.NewAtomicLevelAt(levelToZapLevel(level)),
			Development:      true,
			Encoding:         encodingConsole,
			EncoderConfig:    developmentEncoderConfig,
			OutputPaths:      outputStderr,
			ErrorOutputPaths: outputStderr,
		}
	} else {
		config = &zap.Config{
			Level:            zap.NewAtomicLevelAt(levelToZapLevel(level)),
			Encoding:         encodingJSON,
			EncoderConfig:    productionEncoderConfig,
			OutputPaths:      outputStderr,
			ErrorOutputPaths: outputStderr,
		}
	}

	logger, err := config.Build()
	if err != nil {
		logger = zap.NewNop()
	}

	return logger.Sugar()
}

// NewLoggerFromEnv creates a new logger from the environment. It consumes
// LOG_LEVEL for determining the level and LOG_MODE for determining the output
// parameters.
func NewLoggerFromEnv() *zap.SugaredLogger {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = defaultLogLevel
	}
	development := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_MODE"))) == "development"
	return NewLogger(level, development)
}

// DefaultLogger returns the default logger for the package.
func DefaultLogger() *zap.SugaredLogger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLoggerFromEnv()
	})
	return defaultLogger
}

// WithLogger creates a new context with the provided logger attached.
func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
func From(ctx context.Context) *zap.SugaredLogger {
	return FromContext(ctx)
}

// fromContext returns the logger stored in the context. If no such logger
// exists, a default logger is returned.
func FromContext(ctx context.Context) *zap.SugaredLogger {
	//check type of ctx
	if ginCtx, ok := ctx.(*gin.Context); ok && ginCtx != nil {
		// If ctx is a gin.Context,
		// we need to extract the logger from the gin context.
		if ginCtx.Request != nil {
			ctx = ginCtx.Request.Context()
		}
	}
	if logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return logger
	}
	return DefaultLogger()
}

const (
	timestamp  = "timestamp"
	severity   = "severity"
	logger     = "logger"
	caller     = "caller"
	message    = "message"
	stacktrace = "stacktrace"

	levelDebug     = "DEBUG"
	levelInfo      = "INFO"
	levelWarning   = "WARNING"
	levelError     = "ERROR"
	levelCritical  = "CRITICAL"
	levelAlert     = "ALERT"
	levelEmergency = "EMERGENCY"

	encodingConsole = "console"
	encodingJSON    = "json"
)

var outputStderr = []string{"stderr"}

var productionEncoderConfig = zapcore.EncoderConfig{
	TimeKey:          timestamp,
	LevelKey:         severity,
	NameKey:          logger,
	CallerKey:        caller,
	MessageKey:       message,
	StacktraceKey:    stacktrace,
	LineEnding:       zapcore.DefaultLineEnding,
	EncodeLevel:      levelEncoder(),
	EncodeTime:       timeEncoder(),
	EncodeDuration:   zapcore.SecondsDurationEncoder,
	EncodeCaller:     zapcore.ShortCallerEncoder,
	ConsoleSeparator: "|",
}

var developmentEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "",
	LevelKey:       "L",
	NameKey:        "N",
	CallerKey:      "C",
	FunctionKey:    zapcore.OmitKey,
	MessageKey:     "M",
	StacktraceKey:  "S",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalColorLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

// levelToZapLevel converts the given string to the appropriate zap level
// value.
func levelToZapLevel(s string) zapcore.Level {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case levelDebug:
		return zapcore.DebugLevel
	case levelInfo:
		return zapcore.InfoLevel
	case levelWarning:
		return zapcore.WarnLevel
	case levelError:
		return zapcore.ErrorLevel
	case levelCritical:
		return zapcore.DPanicLevel
	case levelAlert:
		return zapcore.PanicLevel
	case levelEmergency:
		return zapcore.FatalLevel
	}

	return zapcore.WarnLevel
}

// levelEncoder transforms a zap level to the associated stackdriver level.
func levelEncoder() zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString(levelDebug)
		case zapcore.InfoLevel:
			enc.AppendString(levelInfo)
		case zapcore.WarnLevel:
			enc.AppendString(levelWarning)
		case zapcore.ErrorLevel:
			enc.AppendString(levelError)
		case zapcore.DPanicLevel:
			enc.AppendString(levelCritical)
		case zapcore.PanicLevel:
			enc.AppendString(levelAlert)
		case zapcore.FatalLevel:
			enc.AppendString(levelEmergency)
		default:
			enc.AppendString(levelDebug)
		}
	}
}

// timeEncoder encodes the time as RFC3339 nano.
func timeEncoder() zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339Nano))
	}
}
