package log

import (
	"io"
	"os"
)

// appender sint
type Appender interface {
	Flush(bytes []byte)
}

type consoleAppender struct {
	descripter io.Writer
}

func (console *consoleAppender) Flush(buffer []byte) {
	offset := 0
	length := len(buffer)

	for readed, err := console.descripter.Write(buffer[offset:]); err == nil && (offset+readed) < length; readed, err = console.descripter.Write(buffer[offset:]) {
		offset = offset + readed
	}

}

type fileAppender struct {
	descripter io.WriteCloser
	logfile    string
	roll       *RollConfig
}

func (file *fileAppender) Flush(buffer []byte) {
	offset := 0
	length := len(buffer)

	for readed, err := file.descripter.Write(buffer[offset:]); err == nil && (offset+readed) < length; readed, err = file.descripter.Write(buffer[offset:]) {
		offset = offset + readed
	}
}

//NewStdoutAppender new a appender that flush data to stdout
func NewStdoutAppender() Appender {
	return &consoleAppender{descripter: os.Stdout}
}

//NewStderrAppender new a appender that flush data to stderr
func NewStderrAppender() Appender {
	return &consoleAppender{descripter: os.Stderr}
}

//NewFileAppender new a appender that flush data to file system
func NewFileAppender(filePath string, roll *RollConfig) Appender {
	return &fileAppender{logfile: filePath, roll: roll}
}
