package caolog

import (
	"bytes"
	"strings"
	"testing"

	"github.com/CaoStudio/caolog/formatter"
)

// TestNewCommonLogger 测试NewCommonLogger函数
func TestNewCommonLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		writer:    &buf,
		level:     DebugLevel,
		Formatter: formatter.DefaultFormatter(),
	}

	commonLogger := NewCommonLogger(logger)

	// 验证CommonLogger包含Logger
	if commonLogger.Logger.writer != logger.writer {
		t.Error("CommonLogger should contain the same logger")
	}
	if commonLogger.Logger.level != logger.level {
		t.Error("CommonLogger should contain the same logger level")
	}
}

// TestCommonLoggerDebug 测试CommonLogger.Debug方法
func TestCommonLoggerDebug(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		writer:    &buf,
		level:     DebugLevel,
		Formatter: formatter.DefaultFormatter(),
	}

	commonLogger := NewCommonLogger(logger)
	commonLogger.Debug("debug test")

	output := buf.String()
	if !strings.Contains(output, "DEBUG") {
		t.Errorf("Expected DEBUG in output, got: %v", output)
	}
	if !strings.Contains(output, "debug test") {
		t.Errorf("Expected 'debug test' in output, got: %v", output)
	}
}

// TestCommonLoggerInfo 测试CommonLogger.Info方法
func TestCommonLoggerInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		writer:    &buf,
		level:     DebugLevel,
		Formatter: formatter.DefaultFormatter(),
	}

	commonLogger := NewCommonLogger(logger)
	commonLogger.Info("info test")

	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected INFO in output, got: %v", output)
	}
}

// TestCommonLoggerWarn 测试CommonLogger.Warn方法
func TestCommonLoggerWarn(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		writer:    &buf,
		level:     DebugLevel,
		Formatter: formatter.DefaultFormatter(),
	}

	commonLogger := NewCommonLogger(logger)
	commonLogger.Warn("warn test")

	output := buf.String()
	if !strings.Contains(output, "WARN") {
		t.Errorf("Expected WARN in output, got: %v", output)
	}
}

// TestCommonLoggerError 测试CommonLogger.Error方法
func TestCommonLoggerError(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		writer:    &buf,
		level:     DebugLevel,
		Formatter: formatter.DefaultFormatter(),
	}

	commonLogger := NewCommonLogger(logger)
	commonLogger.Error("error test")

	output := buf.String()
	if !strings.Contains(output, "ERROR") {
		t.Errorf("Expected ERROR in output, got: %v", output)
	}
}

// TestCommonLoggerDPanic 测试CommonLogger.DPanic方法
func TestCommonLoggerDPanic(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		writer:    &buf,
		level:     DebugLevel,
		Formatter: formatter.DefaultFormatter(),
	}

	commonLogger := NewCommonLogger(logger)
	commonLogger.DPanic("dpanic test")

	output := buf.String()
	if !strings.Contains(output, "DPANIC") {
		t.Errorf("Expected DPANIC in output, got: %v", output)
	}
}

// TestCommonLoggerPanic 测试CommonLogger.Panic方法
func TestCommonLoggerPanic(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		writer:    &buf,
		level:     DebugLevel,
		Formatter: formatter.DefaultFormatter(),
	}

	commonLogger := NewCommonLogger(logger)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic, but didn't get one")
		}
	}()

	commonLogger.Panic("panic test")
}

// TestCommonLoggerFatal 测试CommonLogger.Fatal方法
// 注意：Fatal会调用os.Exit(1)，所以这里只测试它不会panic
func TestCommonLoggerFatal(t *testing.T) {
	// 由于Fatal会退出程序，我们只测试函数存在且可调用
	// 在实际测试中，应该使用mock或避免测试Fatal
	t.Skip("Fatal calls os.Exit(1), skipping actual test")
}

// TestCommonLoggerMultipleCalls 测试CommonLogger多次调用
func TestCommonLoggerMultipleCalls(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		writer:    &buf,
		level:     DebugLevel,
		Formatter: formatter.DefaultFormatter(),
	}

	commonLogger := NewCommonLogger(logger)

	// 多次调用不同方法
	commonLogger.Debug("debug1")
	commonLogger.Info("info1")
	commonLogger.Warn("warn1")
	commonLogger.Error("error1")
	commonLogger.DPanic("dpanic1")

	output := buf.String()

	// 验证所有日志级别都出现了
	if !strings.Contains(output, "DEBUG") {
		t.Error("DEBUG should appear in output")
	}
	if !strings.Contains(output, "INFO") {
		t.Error("INFO should appear in output")
	}
	if !strings.Contains(output, "WARN") {
		t.Error("WARN should appear in output")
	}
	if !strings.Contains(output, "ERROR") {
		t.Error("ERROR should appear in output")
	}
	if !strings.Contains(output, "DPANIC") {
		t.Error("DPANIC should appear in output")
	}
}

// TestCommonLoggerLevelFiltering 测试CommonLogger日志级别过滤
func TestCommonLoggerLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		writer:    &buf,
		level:     InfoLevel, // 只记录Info及以上级别
		Formatter: formatter.DefaultFormatter(),
	}

	commonLogger := NewCommonLogger(logger)

	// Debug级别应该被过滤
	buf.Reset()
	commonLogger.Debug("should not appear")
	if strings.Contains(buf.String(), "should not appear") {
		t.Error("Debug log should be filtered when level is InfoLevel")
	}

	// Info级别应该出现
	buf.Reset()
	commonLogger.Info("should appear")
	if !strings.Contains(buf.String(), "should appear") {
		t.Error("Info log should appear when level is InfoLevel")
	}
}

// TestCommonLoggerWithMultipleValues 测试CommonLogger记录多个值
func TestCommonLoggerWithMultipleValues(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		writer:    &buf,
		level:     DebugLevel,
		Formatter: formatter.DefaultFormatter(),
	}

	commonLogger := NewCommonLogger(logger)
	commonLogger.Info("value1", "value2", 123, 456.789)

	output := buf.String()
	if !strings.Contains(output, "value1") || !strings.Contains(output, "value2") ||
		!strings.Contains(output, "123") || !strings.Contains(output, "456.789") {
		t.Errorf("Expected all values in output, got: %v", output)
	}
}
