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

var _lvl_mapping map[LogLevel]string = map[LogLevel]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	FatalLevel: "FATAL",
}

// String 日志级别string
func (l LogLevel) String() string {
	return _lvl_mapping[l]
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
