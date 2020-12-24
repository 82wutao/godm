package log

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Logger log master
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Enable(lvl LogLevel) bool
}

type simpleLogger struct {
	queue chan []byte

	async     bool
	level     LogLevel
	layouts   []LayoutElement
	appenders []Appender
}

func layout(layouts []LayoutElement, lvl LogLevel, msg string) ([]byte, time.Time) {
	buffer := make([]string, len(layouts))

	timestamp := time.Now().UTC()
	var (
		pc   uintptr
		file string
		line int
	)

	for i, e := range layouts {
		switch e {
		case LEVEL:
			buffer[i] = lvl.String()
		case PROCESS:
			buffer[i] = os.Args[0]
		case DATATIME:
			buffer[i] = timestamp.Format("15:04:05.000")
		case FILE:
			if pc == 0 {
				pc, file, line, _ = runtime.Caller(2)
			}
			buffer[i] = file
		case FUNC:
			if pc == 0 {
				pc, file, line, _ = runtime.Caller(2)
			}
			buffer[i] = runtime.FuncForPC(pc).Name()
		case LINE:
			if pc == 0 {
				pc, file, line, _ = runtime.Caller(2)
			}
			buffer[i] = strconv.Itoa(line)
		case MESSAGE:
			buffer[i] = msg
		case THREAD:
			buffer[i] = "thread"
		}
	}
	return []byte(strings.Join(buffer, " ")), time.Now()
}

func (l *simpleLogger) Debug(format string, args ...interface{}) {
	if !l.Enable(DebugLevel) {
		return
	}

	msg := fmt.Sprintf(format, args...)
	log, _ := layout(l.layouts, l.level, msg)

	l.queue <- log
}

func (l *simpleLogger) Info(format string, args ...interface{}) {
	if !l.Enable(InfoLevel) {
		return
	}

	msg := fmt.Sprintf(format, args...)
	log, _ := layout(l.layouts, l.level, msg)

	l.queue <- log
}
func (l *simpleLogger) Warn(format string, args ...interface{}) {
	if !l.Enable(WarnLevel) {
		return
	}

	msg := fmt.Sprintf(format, args...)
	log, _ := layout(l.layouts, l.level, msg)

	l.queue <- log
}
func (l *simpleLogger) Error(format string, args ...interface{}) {
	if !l.Enable(ErrorLevel) {
		return
	}

	msg := fmt.Sprintf(format, args...)
	log, _ := layout(l.layouts, l.level, msg)

	l.queue <- log
}
func (l *simpleLogger) Fatal(format string, args ...interface{}) {
	if !l.Enable(FatalLevel) {
		return
	}

	msg := fmt.Sprintf(format, args...)
	log, _ := layout(l.layouts, l.level, msg)

	l.queue <- log
}
func (l *simpleLogger) Enable(lvl LogLevel) bool {
	return l.level < lvl
}

func (l *simpleLogger) outputLog() {
	for {
		select {
		case i := <-l.queue:
			l.appenders[0].Flush(i)
			//TODO exit
		}
	}
}

// NewLogger new a logger instance
func NewLogger(async bool, lvl LogLevel, layout []LayoutElement, appenders ...Appender) Logger {
	// now()
	// curr timeblock start
	// existed          open that old file
	// not existed      trunk untimestamp file
	// tag curr timeblock start
	// tag next timeblocck start
	// open file descripter

	// compare time on a daemon thread{
	// compare time
	// queue msg waiting be send
	// }
	logger := &simpleLogger{
		queue: make(chan []byte),

		async:     async,
		level:     lvl,
		layouts:   layout,
		appenders: appenders,
	}
	go logger.outputLog()

	return logger
}
