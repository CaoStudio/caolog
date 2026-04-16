package core

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

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

// TestNewLogger 测试NewLogger函数
func TestNewLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	if logger.writer != &buf {
		t.Error("Logger writer not set correctly")
	}
	if logger.level != DebugLevel {
		t.Errorf("Expected level DebugLevel, got %v", logger.level)
	}
	
	// 测试nil writer处理
	loggerNil := NewLogger(nil, InfoLevel)
	if loggerNil.writer != os.Stdout {
		t.Error("Logger should default to os.Stdout when writer is nil")
	}
}

// TestWithLevel 测试WithLevel方法
func TestWithLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	newLogger := logger.WithLevel(InfoLevel)
	if newLogger.level != InfoLevel {
		t.Errorf("Expected level InfoLevel, got %v", newLogger.level)
	}
	if newLogger.writer != logger.writer {
		t.Error("Writer should be preserved when creating new logger with level")
	}
}

// TestWithWriter 测试WithWriter方法
func TestWithWriter(t *testing.T) {
	var buf1 bytes.Buffer
	var buf2 bytes.Buffer
	logger := NewLogger(&buf1, DebugLevel)
	
	newLogger := logger.WithWriter(&buf2)
	if newLogger.writer != &buf2 {
		t.Error("Writer not set correctly in new logger")
	}
	if newLogger.level != logger.level {
		t.Error("Level should be preserved when creating new logger with writer")
	}
}

// TestLoggerDebug 测试Logger.Debug方法
func TestLoggerDebug(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Debug()
	if event.level != DebugLevel {
		t.Errorf("Expected DebugLevel, got %v", event.level)
	}
	if event.logger != logger {
		t.Error("Event logger not set correctly")
	}
}

// TestLoggerInfo 测试Logger.Info方法
func TestLoggerInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	if event.level != InfoLevel {
		t.Errorf("Expected InfoLevel, got %v", event.level)
	}
}

// TestLoggerWarn 测试Logger.Warn方法
func TestLoggerWarn(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Warn()
	if event.level != WarnLevel {
		t.Errorf("Expected WarnLevel, got %v", event.level)
	}
}

// TestLoggerError 测试Logger.Error方法
func TestLoggerError(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Error()
	if event.level != ErrorLevel {
		t.Errorf("Expected ErrorLevel, got %v", event.level)
	}
}

// TestLoggerErr 测试Logger.Err方法
func TestLoggerErr(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	// 测试非nil错误
	err := errors.New("test error")
	event := logger.Err(err)
	if event.level != ErrorLevel {
		t.Errorf("Expected ErrorLevel for non-nil error, got %v", event.level)
	}
	
	// 测试nil错误
	eventNil := logger.Err(nil)
	if eventNil.level != InfoLevel {
		t.Errorf("Expected InfoLevel for nil error, got %v", eventNil.level)
	}
}

// TestLoggerFatal 测试Logger.Fatal方法
func TestLoggerFatal(t *testing.T) {
	// 由于Fatal会退出程序，我们只测试函数存在且可调用
	// 在实际测试中，应该使用mock或避免测试Fatal
	t.Skip("Fatal calls os.Exit(1), skipping actual test")
}

// TestLoggerPanic 测试Logger.Panic方法
func TestLoggerPanic(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Panic()
	if event.level != PanicLevel {
		t.Errorf("Expected PanicLevel, got %v", event.level)
	}
	if event.exit == nil {
		t.Error("Panic event should have exit function")
	}
}

// TestEventMsg 测试Event.Msg方法
func TestEventMsg(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	event.Msg("test message")
	
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected 'test message' in output, got: %v", output)
	}
	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected INFO in output, got: %v", output)
	}
}

// TestEventMsgLevelFiltering 测试Event.Msg级别过滤
func TestEventMsgLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, InfoLevel) // 只记录Info及以上级别
	
	// Debug级别应该被过滤
	event := logger.Debug()
	event.Msg("should not appear")
	if strings.Contains(buf.String(), "should not appear") {
		t.Error("Debug log should be filtered when level is InfoLevel")
	}
	
	// Info级别应该出现
	buf.Reset()
	event = logger.Info()
	event.Msg("should appear")
	if !strings.Contains(buf.String(), "should appear") {
		t.Error("Info log should appear when level is InfoLevel")
	}
}

// TestEventErr 测试Event.Err方法
func TestEventErr(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	err := errors.New("test error")
	event := logger.Error()
	event.Err(err).Msg("error occurred")
	
	output := buf.String()
	if !strings.Contains(output, "error") || !strings.Contains(output, "test error") {
		t.Errorf("Expected error field in output, got: %v", output)
	}
}

// TestEventStr 测试Event.Str方法
func TestEventStr(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	event.Str("key", "value").Msg("test")
	
	output := buf.String()
	if !strings.Contains(output, "key=value") {
		t.Errorf("Expected 'key=value' in output, got: %v", output)
	}
}

// TestEventInt 测试Event.Int方法
func TestEventInt(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	event.Int("count", 42).Msg("test")
	
	output := buf.String()
	if !strings.Contains(output, "count=42") {
		t.Errorf("Expected 'count=42' in output, got: %v", output)
	}
}

// TestEventBool 测试Event.Bool方法
func TestEventBool(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	event.Bool("active", true).Msg("test")
	
	output := buf.String()
	if !strings.Contains(output, "active=true") {
		t.Errorf("Expected 'active=true' in output, got: %v", output)
	}
}

// TestEventFloat64 测试Event.Float64方法
func TestEventFloat64(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	event.Float64("pi", 3.14159).Msg("test")
	
	output := buf.String()
	if !strings.Contains(output, "pi=3.14159") {
		t.Errorf("Expected 'pi=3.14159' in output, got: %v", output)
	}
}

// TestEventFloat32 测试Event.Float32方法
func TestEventFloat32(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	event.Float32("value", 3.14).Msg("test")
	
	output := buf.String()
	if !strings.Contains(output, "value=3.14") {
		t.Errorf("Expected 'value=3.14' in output, got: %v", output)
	}
}

// TestEventInt32 测试Event.Int32方法
func TestEventInt32(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	event.Int32("value", 2147483647).Msg("test")
	
	output := buf.String()
	if !strings.Contains(output, "value=2147483647") {
		t.Errorf("Expected 'value=2147483647' in output, got: %v", output)
	}
}

// TestEventInt64 测试Event.Int64方法
func TestEventInt64(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	event.Int64("value", 9223372036854775807).Msg("test")
	
	output := buf.String()
	if !strings.Contains(output, "value=9223372036854775807") {
		t.Errorf("Expected 'value=9223372036854775807' in output, got: %v", output)
	}
}

// TestEventUint 测试Event.Uint方法
func TestEventUint(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	event.Uint("value", 42).Msg("test")
	
	output := buf.String()
	if !strings.Contains(output, "value=42") {
		t.Errorf("Expected 'value=42' in output, got: %v", output)
	}
}

// TestEventUint32 测试Event.Uint32方法
func TestEventUint32(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	event.Uint32("value", 4294967295).Msg("test")
	
	output := buf.String()
	if !strings.Contains(output, "value=4294967295") {
		t.Errorf("Expected 'value=4294967295' in output, got: %v", output)
	}
}

// TestEventUint64 测试Event.Uint64方法
func TestEventUint64(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	event.Uint64("value", 18446744073709551615).Msg("test")
	
	output := buf.String()
	if !strings.Contains(output, "value=18446744073709551615") {
		t.Errorf("Expected 'value=18446744073709551615' in output, got: %v", output)
	}
}

// TestEventDur 测试Event.Dur方法
func TestEventDur(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	event.Dur("duration", 5*time.Second).Msg("test")
	
	output := buf.String()
	if !strings.Contains(output, "duration=") {
		t.Errorf("Expected 'duration=' in output, got: %v", output)
	}
}

// TestEventAny 测试Event.Any方法
func TestEventAny(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	event.Any("data", map[string]interface{}{"key": "value"}).Msg("test")
	
	output := buf.String()
	if !strings.Contains(output, "data=") {
		t.Errorf("Expected 'data=' in output, got: %v", output)
	}
}

// TestEventTime 测试Event.Time方法
func TestEventTime(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	event.Time("timestamp", time.Now()).Msg("test")
	
	output := buf.String()
	if !strings.Contains(output, "timestamp=") {
		t.Errorf("Expected 'timestamp=' in output, got: %v", output)
	}
}

// TestEventMultipleFields 测试Event多个字段
func TestEventMultipleFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	event.Str("user", "john").
		Int("age", 30).
		Bool("active", true).
		Msg("user logged in")
	
	output := buf.String()
	if !strings.Contains(output, "user=john") {
		t.Errorf("Expected 'user=john' in output, got: %v", output)
	}
	if !strings.Contains(output, "age=30") {
		t.Errorf("Expected 'age=30' in output, got: %v", output)
	}
	if !strings.Contains(output, "active=true") {
		t.Errorf("Expected 'active=true' in output, got: %v", output)
	}
	if !strings.Contains(output, "user logged in") {
		t.Errorf("Expected 'user logged in' in output, got: %v", output)
	}
}

// TestPackageLevelFunctions 测试包级函数
func TestPackageLevelFunctions(t *testing.T) {
	var buf bytes.Buffer
	WithWriter(&buf)
	WithLevel(DebugLevel)
	
	// 测试Debug
	buf.Reset()
	Debug().Msg("debug test")
	if !strings.Contains(buf.String(), "DEBUG") {
		t.Error("Debug function should work")
	}
	
	// 测试Info
	buf.Reset()
	Info().Msg("info test")
	if !strings.Contains(buf.String(), "INFO") {
		t.Error("Info function should work")
	}
	
	// 测试Warn
	buf.Reset()
	Warn().Msg("warn test")
	if !strings.Contains(buf.String(), "WARN") {
		t.Error("Warn function should work")
	}
	
	// 测试Error
	buf.Reset()
	Error().Msg("error test")
	if !strings.Contains(buf.String(), "ERROR") {
		t.Error("Error function should work")
	}
	
	// 测试Err
	buf.Reset()
	Err(errors.New("test error")).Msg("error occurred")
	if !strings.Contains(buf.String(), "ERROR") {
		t.Error("Err function should work")
	}
}

// TestWithLogger 测试WithLogger函数
func TestWithLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, InfoLevel)
	
	WithLogger(logger)
	
	// 验证全局logger被设置
	// 注意：由于logger是包级变量，我们需要通过包级函数验证
	Info().Msg("test")
	if !strings.Contains(buf.String(), "INFO") {
		t.Error("WithLogger should set global logger")
	}
}

// TestPackageLevelWithLevel 测试包级WithLevel函数
func TestPackageLevelWithLevel(t *testing.T) {
	var buf bytes.Buffer
	WithWriter(&buf)
	WithLevel(DebugLevel)
	
	// 测试级别过滤
	buf.Reset()
	Debug().Msg("debug test")
	if !strings.Contains(buf.String(), "debug test") {
		t.Error("Debug should appear when level is DebugLevel")
	}
	
	// 改变级别
	WithLevel(InfoLevel)
	buf.Reset()
	Debug().Msg("debug test 2")
	if strings.Contains(buf.String(), "debug test 2") {
		t.Error("Debug should be filtered when level is InfoLevel")
	}
}

// TestPackageLevelWithWriter 测试包级WithWriter函数
func TestPackageLevelWithWriter(t *testing.T) {
	var buf1 bytes.Buffer
	var buf2 bytes.Buffer
	
	WithWriter(&buf1)
	Info().Msg("test1")
	if !strings.Contains(buf1.String(), "test1") {
		t.Error("First writer should receive log")
	}
	
	WithWriter(&buf2)
	Info().Msg("test2")
	if !strings.Contains(buf2.String(), "test2") {
		t.Error("Second writer should receive log")
	}
	if strings.Contains(buf1.String(), "test2") {
		t.Error("First writer should not receive second log")
	}
}

// TestGetValue 测试getValue函数
func TestGetValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"float64", 3.14, "3.14"},
		{"float32", float32(3.14), "3.140000104904175"},  // float32转换为float64会有精度损失
		{"int", 42, "42"},
		{"uint", uint(42), "42"},
		{"int8", int8(127), "127"},
		{"uint8", uint8(255), "255"},
		{"int16", int16(32767), "32767"},
		{"uint16", uint16(65535), "65535"},
		{"int32", int32(2147483647), "2147483647"},
		{"uint32", uint32(4294967295), "4294967295"},
		{"int64", int64(9223372036854775807), "9223372036854775807"},
		{"uint64", uint64(18446744073709551615), "18446744073709551615"},
		{"string", "hello", "hello"},
		{"[]byte", []byte("hello"), "hello"},
		{"error", errors.New("test error"), "test error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getValue(tt.input)
			if result != tt.expected {
				t.Errorf("getValue(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}

	// 测试其他类型（使用fmt.Sprintf）
	t.Run("OtherType", func(t *testing.T) {
		type TestStruct struct {
			Name string
			Age  int
		}
		result := getValue(TestStruct{Name: "John", Age: 30})
		if !strings.Contains(result, "John") || !strings.Contains(result, "30") {
			t.Errorf("getValue for struct failed, got: %v", result)
		}
	})
}

// TestPenultimateIndexByteString 测试penultimateIndexByteString函数
func TestPenultimateIndexByteString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		char     byte
		expected int
	}{
		{"Simple", "hello/world/test.go", '/', 5},  // 倒数第二个 '/' 在索引 5
		{"NoMatch", "hello", '/', -1},
		{"SingleMatch", "test/", '/', -1},  // 只有一个 '/'，应该返回 -1
		{"MultipleMatches", "a/b/c/d", '/', 3},  // 倒数第二个 '/' 在索引 3
		{"EmptyString", "", '/', -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := penultimateIndexByteString(tt.input, tt.char)
			if result != tt.expected {
				t.Errorf("penultimateIndexByteString(%v, %v) = %v, want %v", tt.input, tt.char, result, tt.expected)
			}
		})
	}
}

// TestEventOutput 测试Event.output方法
func TestEventOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	event := logger.Info()
	event.msg = "test message"
	event.output()
	
	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected INFO in output, got: %v", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected 'test message' in output, got: %v", output)
	}
	// 检查输出格式是否包含文件名和行号（格式：filename:line）
	// 由于 runtime.Caller 的调用栈可能不同，我们只检查格式是否正确
	if !strings.Contains(output, ".go:") {
		t.Errorf("Expected file name with .go extension in output, got: %v", output)
	}
}

// TestEventMsgWithExit 测试Event.Msg带exit函数
func TestEventMsgWithExit(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, DebugLevel)
	
	exitCalled := false
	exitFunc := func(msg string) {
		exitCalled = true
		if msg != "panic message" {
			t.Errorf("Expected 'panic message', got: %v", msg)
		}
	}
	
	event := &Event{
		logger: logger,
		level:  PanicLevel,
		msg:    "panic message",
		exit:   exitFunc,
	}
	
	// 调用Msg，应该触发exit函数
	event.Msg("panic message")
	if !exitCalled {
		t.Error("Exit function was not called")
	}
}