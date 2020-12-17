package log

import (
	"io"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials/processcreds"
)

type appender interface {
	flush(bytes []byte)
}

type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
}

///////////////

type RollConfig struct {
	LifeCycle   time.Duration
	NamePattern string
}
type AppenderConfig struct {
	FilePath      string
	rollConf      RollConfig
	FileDescriber io.Writer
}

type LayoutConfig struct {
	process
	level
	datatime-pattern
	thread
	filep
	func
	line
}
type LoggerConfig struct {
	Async     bool
	Level     LogLevel
	Layout    LayoutConfig
	Appenders []AppenderConfig
}

// multitarget out
// playout
// roll

/**
logger

layout
appender
	roll

level


**/
// NewLogger
func NewLogger(config *LoggerConfig) Logger {
	return nil
}
