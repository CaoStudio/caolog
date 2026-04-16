package formatter

import (
	"context"
	"time"
)

// Details 日志详情结构体
type Details struct {
	Level   Level         `json:"level,omitempty"`
	Path    string        `json:"path,omitempty"`
	Time    time.Time     `json:"time,omitempty"`
	Message string        `json:"message,omitempty"`
	Value   []interface{} `json:"value,omitempty"`
}

// Level 日志级别
type Level int8

const (
	DebugLevel  Level = -1
	InfoLevel   Level = 0
	WarnLevel   Level = 1
	ErrorLevel  Level = 2
	DPanicLevel Level = 3
	PanicLevel  Level = 4
	FatalLevel  Level = 5
)

var levelStrings = map[Level]string{
	DebugLevel:  "DEBUG",
	InfoLevel:   "INFO",
	WarnLevel:   "WARN",
	ErrorLevel:  "ERROR",
	DPanicLevel: "DPANIC",
	PanicLevel:  "PANIC",
	FatalLevel:  "FATAL",
}

func (l Level) String() string {
	if l < DebugLevel || l > FatalLevel {
		return "UNKNOWN"
	}
	return levelStrings[l]
}

// Formatter 接口定义日志格式化器
type Formatter interface {
	// Format 格式化日志详情
	Format(ctx context.Context, details *Details) string
}
