package logrus

import (
	"fmt"
	"os"
	"strings"

	"github.com/PS-07/go-microservices/src/api/config"
	"github.com/sirupsen/logrus"
)

var (
	// Log var
	Log *logrus.Logger
)

func init() {
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		level = logrus.DebugLevel
	}
	Log = &logrus.Logger{
		Level: level,
		Out:   os.Stdout,
	}
	if config.IsProduction() {
		Log.Formatter = &logrus.JSONFormatter{}
	} else {
		Log.Formatter = &logrus.TextFormatter{}
	}
}

// Debug func
func Debug(msg string, tags ...string) {
	if Log.Level < logrus.InfoLevel {
		return
	}
	Log.WithFields(parseFields(tags...)).Debug(msg)
}

// Info func
func Info(msg string, tags ...string) {
	if Log.Level < logrus.InfoLevel {
		return
	}
	Log.WithFields(parseFields(tags...)).Info(msg)
}

// Error func
func Error(msg string, err error, tags ...string) {
	if Log.Level < logrus.ErrorLevel {
		return
	}
	msg = fmt.Sprintf("%s - ERROR - %v", msg, err)
	Log.WithFields(parseFields(tags...)).Error(msg)
}

func parseFields(tags ...string) logrus.Fields {
	result := make(logrus.Fields, len(tags))
	for _, tag := range tags {
		els := strings.Split(tag, ":")
		result[strings.TrimSpace(els[0])] = strings.TrimSpace(els[1])
	}
	return result
}
