package logger

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/sirupsen/logrus"
)

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
	//logrus.SetFormatter(&logrus.TextFormatter{
	//	ForceColors:               true,
	//	FullTimestamp:             true,
	//	EnvironmentOverrideColors: true,
	//})
}

type CustomLogger interface {
	Debug(msg ...any)
	Error(msa ...any)
	Info(msg ...any)
	Warn(msg ...any)
}

type customLog struct {
	name   string
	logger *logrus.Entry
}

func (l *customLog) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := fmt.Sprintf("%s[%s]%s", ColorCyan, entry.Time.Format("2006-01-02 15:04:05.000"), ColorReset)

	// Color by log level
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

	// Provider field (optional)
	provider := ""
	if p, ok := entry.Data["provider"]; ok {
		provider = fmt.Sprintf("%s[%s]%s", ColorBlue, p, ColorReset)
		delete(entry.Data, "provider")
	}
	var fields string
	for k, v := range entry.Data {
		fields += fmt.Sprintf(" %s=%v", k, v)
	}

	// caller
	caller := ""
	if entry.Caller != nil {
		caller += path.Base(entry.Caller.File) + ":" + strconv.Itoa(entry.Caller.Line)
	}

	// Final message
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
