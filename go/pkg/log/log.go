package log

import (
	"os"

	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
)

var (
	// L is an alias for the standard logger.
	L = NewLogger()
)

// Logger declares the methods a logger must support.
type Logger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Println(args ...interface{})
	Printf(format string, args ...interface{})
	WithFields(fields logrus.Fields) *logrus.Entry
	WithField(key string, value interface{}) *logrus.Entry
	WithError(err error) *logrus.Entry
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
}

// NewLogger initiates new Logger object.
func NewLogger() Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(joonix.NewFormatter())
	logger.SetLevel(logrus.InfoLevel)
	return logger
}
