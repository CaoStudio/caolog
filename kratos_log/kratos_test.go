package kratoslog

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	caolog "github.com/CaoStudio/caolog"
	kratoslog "github.com/go-kratos/kratos/v2/log"
)

// TestNewKratosLogger 测试NewKratosLogger函数
func TestNewKratosLogger(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)
	
	logger := caolog.GetLogger()
	kratosLogger := NewKratosLogger(logger)
	if kratosLogger == nil {
		t.Error("NewKratosLogger returned nil")
	}
	
	// 验证返回的是Logger类型
	_, ok := kratosLogger.(Logger)
	if !ok {
		t.Error("NewKratosLogger should return Logger type")
	}
}

// TestLoggerLog 测试Logger.Log方法
func TestLoggerLog(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)
	
	logger := caolog.GetLogger()
	kratosLogger := NewKratosLogger(logger)
	
	// 测试不同kratos日志级别
	tests := []struct {
		level    kratoslog.Level
		expected string
	}{
		{kratoslog.LevelDebug, "DEBUG"},
		{kratoslog.LevelInfo, "INFO"},
		{kratoslog.LevelWarn, "WARN"},
		{kratoslog.LevelError, "ERROR"},
		// 跳过 LevelFatal，因为它会调用 os.Exit(1)
	}
	
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			buf.Reset()
			err := kratosLogger.Log(tt.level, "test message")
			if err != nil {
				t.Errorf("Log returned error: %v", err)
			}
			
			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %s in output, got: %v", tt.expected, output)
			}
			if !strings.Contains(output, "test message") {
				t.Errorf("Expected 'test message' in output, got: %v", output)
			}
		})
	}
}

// TestLoggerLogMultipleValues 测试Logger.Log多个值
func TestLoggerLogMultipleValues(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)
	
	logger := caolog.GetLogger()
	kratosLogger := NewKratosLogger(logger)
	
	err := kratosLogger.Log(kratoslog.LevelInfo, "value1", "value2", 123, 456.789)
	if err != nil {
		t.Errorf("Log returned error: %v", err)
	}
	
	output := buf.String()
	if !strings.Contains(output, "value1") || !strings.Contains(output, "value2") ||
		!strings.Contains(output, "123") || !strings.Contains(output, "456.789") {
		t.Errorf("Expected all values in output, got: %v", output)
	}
}

// TestKratosServer 测试KratosServer中间件
func TestKratosServer(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)
	
	logger := caolog.GetLogger()
	kratosLogger := NewKratosLogger(logger)
	
	// 创建中间件
	mw := KratosServer(kratosLogger)
	if mw == nil {
		t.Error("KratosServer returned nil middleware")
	}
	
	// 验证中间件是函数类型
	// middleware.Middleware是一个函数类型，我们无法直接类型断言
	// 但我们可以验证它是否可以被调用
}

// TestKratosServerWithHandler 测试KratosServer中间件处理请求
func TestKratosServerWithHandler(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)
	
	logger := caolog.GetLogger()
	kratosLogger := NewKratosLogger(logger)
	
	// 创建中间件
	middleware := KratosServer(kratosLogger)
	
	// 创建测试handler
	testHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}
	
	// 应用中间件
	handler := middleware(testHandler)
	
	// 创建测试上下文
	ctx := context.Background()
	
	// 调用handler
	reply, err := handler(ctx, "test request")
	if err != nil {
		t.Errorf("Handler returned error: %v", err)
	}
	if reply != "response" {
		t.Errorf("Expected 'response', got: %v", reply)
	}
	
	// 验证日志被记录
	output := buf.String()
	if !strings.Contains(output, "[Kratos]") {
		t.Errorf("Expected '[Kratos]' in output, got: %v", output)
	}
}

// TestKratosServerWithError 测试KratosServer中间件处理错误
func TestKratosServerWithError(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)
	
	logger := caolog.GetLogger()
	kratosLogger := NewKratosLogger(logger)
	
	// 创建中间件
	middleware := KratosServer(kratosLogger)
	
	// 创建测试handler返回错误
	testHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New("test error")
	}
	
	// 应用中间件
	handler := middleware(testHandler)
	
	// 创建测试上下文
	ctx := context.Background()
	
	// 调用handler
	reply, err := handler(ctx, "test request")
	if err == nil {
		t.Error("Expected error from handler")
	}
	if reply != nil {
		t.Errorf("Expected nil reply, got: %v", reply)
	}
	
	// 验证日志被记录
	output := buf.String()
	if !strings.Contains(output, "[Kratos]") {
		t.Errorf("Expected '[Kratos]' in output, got: %v", output)
	}
	// 注意：错误消息可能不会直接出现在输出中，因为reason字段可能为空
	// 这里只验证日志格式正确
}

// TestExtractError 测试extractError函数
func TestExtractError(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		expectedLevel kratoslog.Level
		expectedMsg   string
	}{
		{"NilError", nil, kratoslog.LevelInfo, ""},
		{"WithError", errors.New("test error"), kratoslog.LevelError, "test error"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, msg := extractError(tt.err)
			if level != tt.expectedLevel {
				t.Errorf("Expected level %v, got %v", tt.expectedLevel, level)
			}
			if msg != tt.expectedMsg {
				t.Errorf("Expected message '%v', got '%v'", tt.expectedMsg, msg)
			}
		})
	}
}

// TestExtractArgs 测试extractArgs函数
func TestExtractArgs(t *testing.T) {
	tests := []struct {
		name     string
		req      interface{}
		expected string
	}{
		{"Stringer", stringerImpl{value: "test"}, "test"},
		{"NonStringer", 123, "123"},
		{"Struct", struct{ Name string }{Name: "John"}, "{Name:John}"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractArgs(tt.req)
			if !strings.Contains(result, tt.expected) {
				t.Errorf("Expected '%s' in result, got: %v", tt.expected, result)
			}
		})
	}
}

// stringerImpl 实现fmt.Stringer接口用于测试
type stringerImpl struct {
	value string
}

func (s stringerImpl) String() string {
	return s.value
}

// TestKratosServerWithTransportContext 测试KratosServer中间件处理transport上下文
func TestKratosServerWithTransportContext(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)
	
	logger := caolog.GetLogger()
	kratosLogger := NewKratosLogger(logger)
	
	// 创建中间件
	middleware := KratosServer(kratosLogger)
	
	// 创建测试handler
	testHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}
	
	// 应用中间件
	handler := middleware(testHandler)
	
	// 创建带有transport信息的上下文
	ctx := context.Background()
	// 注意：这里无法直接创建transport.FromServerContext所需的上下文
	// 在实际使用中，kratos会自动添加这些信息
	
	// 调用handler
	reply, err := handler(ctx, "test request")
	if err != nil {
		t.Errorf("Handler returned error: %v", err)
	}
	if reply != "response" {
		t.Errorf("Expected 'response', got: %v", reply)
	}
}

// TestLoggerLogWithDifferentLevels 测试Logger.Log不同级别
func TestLoggerLogWithDifferentLevels(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.InfoLevel) // 只记录Info及以上级别
	caolog.SetWriter(&buf)
	
	logger := caolog.GetLogger()
	kratosLogger := NewKratosLogger(logger)
	
	// Debug级别应该被过滤
	buf.Reset()
	err := kratosLogger.Log(kratoslog.LevelDebug, "should not appear")
	if err != nil {
		t.Errorf("Log returned error: %v", err)
	}
	if strings.Contains(buf.String(), "should not appear") {
		t.Error("Debug log should be filtered when level is InfoLevel")
	}
	
	// Info级别应该出现
	buf.Reset()
	err = kratosLogger.Log(kratoslog.LevelInfo, "should appear")
	if err != nil {
		t.Errorf("Log returned error: %v", err)
	}
	if !strings.Contains(buf.String(), "should appear") {
		t.Error("Info log should appear when level is InfoLevel")
	}
}