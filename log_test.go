package caolog

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/CaoStudio/caolog/formatter"
)

// TestInitLogger 测试InitLogger函数
func TestInitLogger(t *testing.T) {
	// 测试初始化不同日志级别
	levels := []Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel, DPanicLevel, PanicLevel, FatalLevel}
	for _, level := range levels {
		t.Run(level.String(), func(t *testing.T) {
			InitLogger(level)
			if logger.level != level {
				t.Errorf("Expected logger level %v, got %v", level, logger.level)
			}
			if logLevel != level {
				t.Errorf("Expected global logLevel %v, got %v", level, logLevel)
			}
		})
	}

	// 测试带选项的初始化
	t.Run("WithOptions", func(t *testing.T) {
		optionCalled := false
		testOption := func(ctx context.Context, details *Details) {
			optionCalled = true
		}
		InitLogger(DebugLevel, testOption)
		if len(logger.Options) != 1 {
			t.Errorf("Expected 1 option, got %d", len(logger.Options))
		}
		// 调用option验证
		logger.Options[0](context.Background(), &Details{})
		if !optionCalled {
			t.Error("Option was not called")
		}
	})
}

// TestJSONFormatter 测试JSON格式化器
func TestJSONFormatter(t *testing.T) {
	t.Run("DefaultJSONFormatter", func(t *testing.T) {
		formatter := formatter.JSONFormatterFactory(false)
		details := Details{
			Level:   InfoLevel,
			Path:    "test.go:42",
			Time:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			Message: "test message",
		}
		
		output := formatter.Format(context.Background(), &details)
		
		// 验证输出是有效的JSON
		if !strings.Contains(output, "\"level\":\"INFO\"") {
			t.Errorf("Expected JSON to contain level field, got: %s", output)
		}
		if !strings.Contains(output, "\"message\":\"test message\"") {
			t.Errorf("Expected JSON to contain message field, got: %s", output)
		}
	})
}

// TestTextFormatter 测试文本格式化器
func TestTextFormatter(t *testing.T) {
	t.Run("DefaultTextFormatter", func(t *testing.T) {
		formatter := formatter.DefaultFormatter()
		details := Details{
			Level:   InfoLevel,
			Path:    "test.go:42",
			Time:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			Message: "test message",
		}
		
		output := formatter.Format(context.Background(), &details)
		
		// 验证输出格式
		if !strings.Contains(output, "[INFO]") {
			t.Errorf("Expected output to contain [INFO], got: %s", output)
		}
		if !strings.Contains(output, "test.go:42") {
			t.Errorf("Expected output to contain path, got: %s", output)
		}
		if !strings.Contains(output, "test message") {
			t.Errorf("Expected output to contain message, got: %s", output)
		}
	})
}

// TestSetFormatter 测试设置格式化器
func TestSetFormatter(t *testing.T) {
	t.Run("SetJSONFormatter", func(t *testing.T) {
		// 保存原始格式化器
		originalFormatter := GetLogger().Formatter
		
		// 设置JSON格式化器
		jsonFormatter := formatter.JSONFormatterFactory(false)
		SetFormatter(jsonFormatter)
		
		// 验证格式化器已设置
		logger := GetLogger()
		if logger.Formatter != jsonFormatter {
			t.Error("Formatter was not set correctly")
		}
		
		// 恢复原始格式化器
		SetFormatter(originalFormatter)
	})
}
