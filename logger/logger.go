package logger

import (
	"fmt"
	"log/syslog"
	"os"

	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

var logLevel = logrus.DebugLevel

type logFormatter struct {
	name string
}

func (lf *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	logEntry := fmt.Sprintf("[metadata={process='containeragent',function='containeragent',TMG_CLUSTER_NAME='%s',TMG_ZONE_NAME='%s',POD_IP='%s'}", os.Getenv("TMG_CLUSTER_NAME"), os.Getenv("TMG_ZONE_NAME"), os.Getenv("POD_IP")) + fmt.Sprintf("] [%-5s] [microgateway] %s - %s\n", getLevel(entry.Level), lf.name, entry.Message)

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

	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "go-mlca")
	if err == nil {
		logger.Hooks.Add(hook)
	}

	logger.SetFormatter(&logFormatter{name: loggerName})
	logger.SetLevel(logLevel)
	logger.SetOutput(os.Stdout)

	return logger
}
