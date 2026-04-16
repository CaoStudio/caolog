package formatter

import (
	"context"
	"strings"
	"testing"
	"time"
)

// TestTextFormatter 测试文本格式化器
func TestTextFormatter(t *testing.T) {
	formatter := &TextFormatter{}
	details := &Details{
		Level:   InfoLevel,
		Path:    "test.go:123",
		Time:    time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC),
		Message: "test message",
	}

	result := formatter.Format(context.Background(), details)

	if !strings.Contains(result, "INFO") {
		t.Errorf("Expected INFO in output, got: %v", result)
	}
	if !strings.Contains(result, "test message") {
		t.Errorf("Expected 'test message' in output, got: %v", result)
	}
	if !strings.Contains(result, "test.go:123") {
		t.Errorf("Expected 'test.go:123' in output, got: %v", result)
	}
}

// TestJSONFormatter 测试JSON格式化器
func TestJSONFormatter(t *testing.T) {
	formatter := &JSONFormatter{}
	details := &Details{
		Level:   InfoLevel,
		Path:    "test.go:123",
		Time:    time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC),
		Message: "test message",
	}

	result := formatter.Format(context.Background(), details)

	if !strings.Contains(result, "INFO") {
		t.Errorf("Expected INFO in output, got: %v", result)
	}
	if !strings.Contains(result, "test message") {
		t.Errorf("Expected 'test message' in output, got: %v", result)
	}
}

// TestJSONFormatterPretty 测试美化输出的JSON格式化器
func TestJSONFormatterPretty(t *testing.T) {
	formatter := &JSONFormatter{Pretty: true}
	details := &Details{
		Level:   InfoLevel,
		Path:    "test.go:123",
		Time:    time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC),
		Message: "test message",
	}

	result := formatter.Format(context.Background(), details)

	if !strings.Contains(result, "\n") {
		t.Errorf("Expected pretty JSON to contain newlines, got: %v", result)
	}
}

// TestJSONFormatterFieldMap 测试自定义字段映射的JSON格式化器
func TestJSONFormatterFieldMap(t *testing.T) {
	formatter := &JSONFormatter{
		FieldMap: map[string]string{
			"severity":  "level",
			"timestamp": "time",
			"location":  "path",
			"msg":       "message",
		},
	}
	details := &Details{
		Level:   InfoLevel,
		Path:    "test.go:123",
		Time:    time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC),
		Message: "test message",
	}

	result := formatter.Format(context.Background(), details)

	if !strings.Contains(result, "severity") {
		t.Errorf("Expected 'severity' in output, got: %v", result)
	}
	if !strings.Contains(result, "INFO") {
		t.Errorf("Expected INFO in output, got: %v", result)
	}
}

// TestDefaultFormatter 测试默认格式化器
func TestDefaultFormatter(t *testing.T) {
	formatter := DefaultFormatter()
	if formatter == nil {
		t.Error("DefaultFormatter should not return nil")
	}
	if _, ok := formatter.(*TextFormatter); !ok {
		t.Error("DefaultFormatter should return *TextFormatter")
	}
}

// TestFormatDetails 测试FormatDetails函数
func TestFormatDetails(t *testing.T) {
	details := &Details{
		Level:   InfoLevel,
		Path:    "test.go:123",
		Time:    time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC),
		Message: "test message",
	}

	// 测试nil formatter
	result := FormatDetails(context.Background(), details, nil)
	if !strings.Contains(result, "INFO") {
		t.Errorf("Expected INFO in output, got: %v", result)
	}

	// 测试自定义formatter
	formatter := &TextFormatter{TimeFormat: "2006-01-02"}
	result = FormatDetails(context.Background(), details, formatter)
	if !strings.Contains(result, "2026-04-16") {
		t.Errorf("Expected custom time format in output, got: %v", result)
	}
}

// TestLevelString 测试Level.String方法
func TestLevelString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{WarnLevel, "WARN"},
		{ErrorLevel, "ERROR"},
		{DPanicLevel, "DPANIC"},
		{PanicLevel, "PANIC"},
		{FatalLevel, "FATAL"},
		{Level(100), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.level.String()
			if result != tt.expected {
				t.Errorf("Level(%d).String() = %v, want %v", tt.level, result, tt.expected)
			}
		})
	}
}

// TestJSONFormatterFactory 测试JSONFormatterFactory函数
func TestJSONFormatterFactory(t *testing.T) {
	formatter := JSONFormatterFactory(true)
	if formatter == nil {
		t.Error("JSONFormatterFactory should not return nil")
	}
	if jsonFormatter, ok := formatter.(*JSONFormatter); !ok {
		t.Error("JSONFormatterFactory should return *JSONFormatter")
	} else if !jsonFormatter.Pretty {
		t.Error("JSONFormatterFactory should set Pretty to true")
	}

	formatter = JSONFormatterFactory(false)
	if jsonFormatter, ok := formatter.(*JSONFormatter); !ok {
		t.Error("JSONFormatterFactory should return *JSONFormatter")
	} else if jsonFormatter.Pretty {
		t.Error("JSONFormatterFactory should set Pretty to false")
	}
}

// TestJSONFormatterValueField 测试JSON格式化器的Value字段序列化
func TestJSONFormatterValueField(t *testing.T) {
	t.Run("ValueFieldSerialized", func(t *testing.T) {
		formatter := &JSONFormatter{}
		details := &Details{
			Level:   InfoLevel,
			Path:    "test.go:123",
			Time:    time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC),
			Message: "test message",
			Value:   []interface{}{"traceID:", "abc123", "extra", 42},
		}

		result := formatter.Format(context.Background(), details)

		// 验证value字段在输出中
		if !strings.Contains(result, "\"value\"") {
			t.Errorf("Expected 'value' field in output, got: %v", result)
		}
		// 验证value数组内容
		if !strings.Contains(result, "traceID:") {
			t.Errorf("Expected 'traceID:' in value array, got: %v", result)
		}
	})

	t.Run("ValueFieldEmpty", func(t *testing.T) {
		formatter := &JSONFormatter{}
		details := &Details{
			Level:   InfoLevel,
			Path:    "test.go:123",
			Time:    time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC),
			Message: "test message",
			Value:   []interface{}{},
		}

		result := formatter.Format(context.Background(), details)

		// 空Value字段不应该出现在输出中
		if strings.Contains(result, "\"value\"") {
			t.Errorf("Empty value field should not appear in output, got: %v", result)
		}
	})

	t.Run("ValueFieldWithFieldMap", func(t *testing.T) {
		formatter := &JSONFormatter{
			FieldMap: map[string]string{
				"severity":  "level",
				"timestamp": "time",
				"location":  "path",
				"msg":       "message",
				"extra":     "value",
			},
		}
		details := &Details{
			Level:   InfoLevel,
			Path:    "test.go:123",
			Time:    time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC),
			Message: "test message",
			Value:   []interface{}{"traceID:", "abc123"},
		}

		result := formatter.Format(context.Background(), details)

		// 验证自定义字段名
		if !strings.Contains(result, "\"extra\"") {
			t.Errorf("Expected 'extra' field in output, got: %v", result)
		}
		if !strings.Contains(result, "traceID:") {
			t.Errorf("Expected 'traceID:' in extra array, got: %v", result)
		}
	})
}
