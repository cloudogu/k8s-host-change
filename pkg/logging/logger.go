package logging

import (
	"fmt"
	"os"

	"github.com/bombsimon/logrusr/v2"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	ctrl "sigs.k8s.io/controller-runtime"
)

const logLevelEnvName = "LOG_LEVEL"

const (
	errorLevel int = iota
	warningLevel
	infoLevel
	debugLevel
)

type logSink interface {
	logr.LogSink
}

type libraryLogger struct {
	logger logSink
	name   string
}

func (l *libraryLogger) log(level int, args ...interface{}) {
	l.logger.Info(level, fmt.Sprintf("[%s] %s", l.name, fmt.Sprint(args...)))
}

func (l *libraryLogger) logf(level int, format string, args ...interface{}) {
	l.logger.Info(level, fmt.Sprintf("[%s] %s", l.name, fmt.Sprintf(format, args...)))
}

// Debug logs the given arguments with the debug log-level.
func (l *libraryLogger) Debug(args ...interface{}) {
	l.log(debugLevel, args...)
}

// Info logs the given arguments with the info log-level.
func (l *libraryLogger) Info(args ...interface{}) {
	l.log(infoLevel, args...)
}

// Warning logs the given arguments with the warning log-level.
func (l *libraryLogger) Warning(args ...interface{}) {
	l.log(warningLevel, args...)
}

// Error logs the given arguments with the error log-level.
func (l *libraryLogger) Error(args ...interface{}) {
	l.log(errorLevel, args...)
}

// Debugf formats the arguments into the given format string and logs it with the debug log-level.
func (l *libraryLogger) Debugf(format string, args ...interface{}) {
	l.logf(debugLevel, format, args...)
}

// Infof formats the arguments into the given format string and logs it with the info log-level.
func (l *libraryLogger) Infof(format string, args ...interface{}) {
	l.logf(infoLevel, format, args...)
}

// Warningf formats the arguments into the given format string and logs it with the warning log-level.
func (l *libraryLogger) Warningf(format string, args ...interface{}) {
	l.logf(warningLevel, format, args...)
}

// Errorf formats the arguments into the given format string and logs it with the error log-level.
func (l *libraryLogger) Errorf(format string, args ...interface{}) {
	l.logf(errorLevel, format, args...)
}

func getLogLevelFromEnv() (logrus.Level, error) {
	logLevel, found := os.LookupEnv(logLevelEnvName)
	if !found {
		return logrus.ErrorLevel, nil
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return logrus.ErrorLevel, fmt.Errorf("value of log environment variable [%s] is not a valid log level: %w", logLevelEnvName, err)
	}

	return level, nil
}

// ConfigureLogger sets the logrus logger as for all logging implementations from the controller-runtime.
func ConfigureLogger() error {
	level, err := getLogLevelFromEnv()
	if err != nil {
		return err
	}

	// create logrus logger that can be styled and formatted
	logrusLog := logrus.New()
	logrusLog.SetFormatter(&logrus.TextFormatter{})
	logrusLog.SetLevel(level)

	// convert logrus logger to logr logger
	logrusrLogger := logrusr.New(logrusLog)

	// set logr logger as controller logger
	ctrl.SetLogger(logrusrLogger)

	return nil
}
