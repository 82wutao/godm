package log

// LogLevel 日志级别
type LogLevel int

//  日志级别 枚举
const (
	_ LogLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// String 日志级别string
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

// LayoutElement 日志组成元素
type LayoutElement int

// 日志组成元素 元素种类枚举
const (
	_ LayoutElement = iota
	LEVEL
	PROCESS
	DATATIME
	//THREAD
	FILE
	FUNC
	LINE
	MESSAGE
)
