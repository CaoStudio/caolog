package formatter

import (
	"context"
	"encoding/json"
	"strings"
	"time"
)

// Details 日志详情结构体
type Details struct {
	Level   Level         `json:"level,omitempty"`
	Path    string        `json:"path,omitempty"`
	Time    time.Time     `json:"time,omitempty"`
	Message string        `json:"message,omitempty"`
	Value   []interface{} `json:"-"`
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

// TextFormatter 文本格式化器（默认格式）
type TextFormatter struct {
	// 时间格式化模板
	TimeFormat string
}

// Format 实现TextFormatter的格式化方法
func (f *TextFormatter) Format(ctx context.Context, details *Details) string {
	if f.TimeFormat == "" {
		f.TimeFormat = "2006-01-02 - 15:04:05"
	}

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
	return &TextFormatter{}
}

// JSONFormatter JSON格式化器
type JSONFormatter struct {
	// 是否美化输出
	Pretty bool
	// 自定义字段映射
	FieldMap map[string]string
}

// Format 实现JSONFormatter的格式化方法
func (f *JSONFormatter) Format(ctx context.Context, details *Details) string {
	// 创建JSON结构
	jsonData := make(map[string]interface{})

	// 映射字段
	if f.FieldMap != nil {
		// 使用自定义字段映射
		for key, field := range f.FieldMap {
			switch field {
			case "level":
				jsonData[key] = details.Level.String()
			case "time":
				jsonData[key] = details.Time.Format(time.RFC3339)
			case "path":
				jsonData[key] = details.Path
			case "message":
				jsonData[key] = details.Message
			}
		}
	} else {
		// 使用默认字段名
		jsonData["level"] = details.Level.String()
		jsonData["time"] = details.Time.Format(time.RFC3339)
		jsonData["path"] = details.Path
		jsonData["message"] = details.Message
	}

	// 序列化为JSON
	var output []byte
	var err error
	if f.Pretty {
		output, err = json.MarshalIndent(jsonData, "", "  ")
	} else {
		output, err = json.Marshal(jsonData)
	}

	if err != nil {
		return "{\"error\":\"JSON marshal failed: " + err.Error() + "\"}\n"
	}

	return string(output) + "\n"
}

// JSONFormatterFactory 创建JSON格式化器
func JSONFormatterFactory(pretty bool) Formatter {
	return &JSONFormatter{Pretty: pretty}
}

// FormatDetails 格式化日志详情
func FormatDetails(ctx context.Context, details *Details, formatter Formatter) string {
	if formatter == nil {
		formatter = DefaultFormatter()
	}
	return formatter.Format(ctx, details)
}