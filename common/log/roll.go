package log

import "time"

// RollConfig 日志文件卷动配置
type RollConfig struct {
	LifeCycle       time.Duration
	DateTimePattern string
}
