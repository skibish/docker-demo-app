package log

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

// Logger describes logger structure
type Logger struct {
	lgrs *logrus.Logger
}

// NewNullLogger return logger that do not write anywhere
func NewNullLogger() *Logger {
	logger, _ := test.NewNullLogger()
	return &Logger{
		lgrs: logger,
	}
}

// NewLogger creates logger instance
func NewLogger() *Logger {
	var logger = logrus.New()
	logger.Out = os.Stdout
	logger.Formatter = &logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	}

	return &Logger{
		lgrs: logger,
	}
}

// RequestLog log requests
func (l *Logger) RequestLog(r *http.Request) {
	l.lgrs.WithFields(logrus.Fields{
		"method": r.Method,
		"type":   "request",
		"path":   r.URL.Path,
	}).Info("Request")
}

// ResponseLog log responses
func (l *Logger) ResponseLog(r *http.Request) {
	l.lgrs.WithFields(logrus.Fields{
		"method": r.Method,
		"type":   "response",
		"path":   r.URL.Path,
	}).Info("Response")
}

// Print prints INFO level message
func (l *Logger) Print(args ...interface{}) {
	l.lgrs.Info(args...)
}

// Printf prints INFO level message
func (l *Logger) Printf(format string, args ...interface{}) {
	l.lgrs.Infof(format, args...)
}

// Error prints ERROR level message and exists
func (l *Logger) Error(args ...interface{}) {
	l.lgrs.Error(args...)

}

// Errorf prints ERROR level message and exists
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.lgrs.Errorf(format, args...)
}

// Fatal prints FATAL level message and exists
func (l *Logger) Fatal(args ...interface{}) {
	l.lgrs.Fatal(args...)

}

// Fatalf prints FATAL level message and exists
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.lgrs.Fatalf(format, args...)
}
