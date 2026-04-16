package formatter

import (
	"context"
	"strings"
	"sync"
)

// TextFormatter 文本格式化器（默认格式）
type TextFormatter struct {
	// 时间格式化模板
	TimeFormat string
	once       sync.Once
}

// NewTextFormatter 创建文本格式化器
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		TimeFormat: "2006-01-02 - 15:04:05",
	}
}

// Format 实现TextFormatter的格式化方法
func (f *TextFormatter) Format(ctx context.Context, details *Details) string {
	// 使用sync.Once确保线程安全地设置默认值
	f.once.Do(func() {
		if f.TimeFormat == "" {
			f.TimeFormat = "2006-01-02 - 15:04:05"
		}
	})

	var builder strings.Builder
	builder.Grow(100)

	// 时间
	builder.WriteString("[")
	builder.WriteString(details.Time.Format(f.TimeFormat))
	builder.WriteString("] ")

	// 级别
	builder.WriteString("[")
	builder.WriteString(details.Level.String())
	builder.WriteString("] ")

	// 调用位置
	builder.WriteString(details.Path)
	builder.WriteString(" ")

	// 消息
	builder.WriteString(details.Message)
	builder.WriteString("\n")

	return builder.String()
}

// DefaultFormatter 返回默认格式化器（文本格式）
func DefaultFormatter() Formatter {
	return NewTextFormatter()
}

// FormatDetails 格式化日志详情
func FormatDetails(ctx context.Context, details *Details, formatter Formatter) string {
	if formatter == nil {
		formatter = DefaultFormatter()
	}
	return formatter.Format(ctx, details)
}
