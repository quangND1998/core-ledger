package logger

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/sirupsen/logrus"
)

var sensitiveKeys = []string{
	"password",
	"key",
	"token",
	"secret",
	"credential",
	"private",
	"api_key",
	"access_token",
}

func maskSensitiveString(msg string) string {
	for _, key := range sensitiveKeys {
		patterns := []string{
			key + `=\S+`,
			key + `:\s*\S+`,
			key + `\s+\S+`,
		}
		for _, pattern := range patterns {
			re := regexp.MustCompile(`(?i)` + pattern)
			msg = re.ReplaceAllString(msg, key+"=********")
		}
	}
	return msg
}

func maskSensitiveValue(v any) any {
	if s, ok := v.(string); ok {
		return maskSensitiveString(s)
	}
	return v
}

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
)

func init() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&customLog{})
	logrus.SetOutput(os.Stdout)
}

type CustomLogger interface {
	Debug(msg ...any)
	Error(msg ...any)
	Info(msg ...any)
	Warn(msg ...any)
}

type customLog struct {
	name   string
	logger *logrus.Entry
}

func (l *customLog) Format(entry *logrus.Entry) ([]byte, error) {
	// Mask toàn bộ dữ liệu nhạy cảm
	entry.Message = maskSensitiveString(entry.Message)
	for k, v := range entry.Data {
		entry.Data[k] = maskSensitiveValue(v)
	}

	timestamp := fmt.Sprintf("%s[%s]%s", ColorCyan, entry.Time.Format("2006-01-02 15:04:05.000"), ColorReset)

	var levelColor string
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = ColorGreen
	case logrus.WarnLevel:
		levelColor = ColorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = ColorRed
	default:
		levelColor = ColorBlue
	}
	level := fmt.Sprintf("%s[%s]%s", levelColor, entry.Level.String(), ColorReset)

	provider := ""
	if p, ok := entry.Data["provider"]; ok {
		provider = fmt.Sprintf("%s[%s]%s", ColorBlue, p, ColorReset)
		delete(entry.Data, "provider")
	}

	var fields string
	for k, v := range entry.Data {
		fields += fmt.Sprintf(" %s=%v", k, v)
	}

	caller := ""
	if entry.Caller != nil {
		caller += path.Base(entry.Caller.File) + ":" + strconv.Itoa(entry.Caller.Line)
	}

	msg := fmt.Sprintf("%s%s[%v] %s %s%s\n", level, timestamp, caller, provider, entry.Message, fields)
	return []byte(msg), nil
}

func NewSystemLog(name string) CustomLogger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(true)
	logEntry := logrus.NewEntry(logger)
	ct := &customLog{
		name:   name,
		logger: logEntry,
	}
	logger.SetFormatter(ct)
	return ct
}

func (l *customLog) Debug(msg ...any) {
	logrus.Debug("[", l.name, "] ", msg)
}

func (l *customLog) Error(msg ...any) {
	logrus.Error("[", l.name, "] ", msg)
}

func (l *customLog) Info(msg ...any) {
	logrus.Info("[", l.name, "] ", msg)
}

func (l *customLog) Warn(msg ...any) {
	logrus.Warn("[", l.name, "] ", msg)
}
