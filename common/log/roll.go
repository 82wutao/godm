package log

import "time"

type RollConfig struct {
	LifeCycle       time.Duration
	DateTimePattern string
}
