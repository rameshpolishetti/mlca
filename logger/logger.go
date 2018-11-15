package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var logLevel = logrus.DebugLevel

type logFormatter struct {
	name string
}

func (lf *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	logEntry := fmt.Sprintf("%s %-5s [microgateway][%s] - %s \n", entry.Time.Format("2006-01-02 15:04:05.000"), getLevel(entry.Level), lf.name, entry.Message)
	return []byte(logEntry), nil
}

func getLevel(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "DEBUG"
	case logrus.InfoLevel:
		return "INFO"
	case logrus.ErrorLevel:
		return "ERROR"
	case logrus.WarnLevel:
		return "WARN"
	case logrus.PanicLevel:
		return "PANIC"
	case logrus.FatalLevel:
		return "FATAL"
	}
	return "UNKNOWN"
}

// GetLogger returs logger
func GetLogger(loggerName string) *logrus.Logger {
	logger := logrus.New()
	// logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetFormatter(&logFormatter{name: loggerName})
	logger.SetLevel(logLevel)
	logger.SetOutput(os.Stdout)

	return logger
}
