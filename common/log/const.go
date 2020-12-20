package log

type LogLevel int

const (
	_ LogLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	}
	return ""
}

type LayoutElement int

const (
	_ LayoutElement = iota
	LEVEL
	PROCESS
	DATATIME
	THREAD
	FILE
	FUNC
	LINE
	MESSAGE
)
