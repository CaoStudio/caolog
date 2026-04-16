package plugin

import (
	"bytes"
	"context"
	"strings"
	"testing"

	caolog "github.com/CaoStudio/caolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// TestNewTrace 测试NewTrace函数
func TestNewTrace(t *testing.T) {
	tracePlugin := NewTrace()
	if tracePlugin == nil {
		t.Error("NewTrace returned nil")
	}
	if tracePlugin.tracerProvider == nil {
		t.Error("tracerProvider should not be nil")
	}
	if tracePlugin.tracer == nil {
		t.Error("tracer should not be nil")
	}
}

// TestGetLogTag 测试GetLogTag方法
func TestGetLogTag(t *testing.T) {
	tracePlugin := NewTrace()
	
	tests := []struct {
		level    caolog.Level
		expected string
	}{
		{caolog.DebugLevel, "Log.DEBUG"},
		{caolog.InfoLevel, "Log.INFO"},
		{caolog.WarnLevel, "Log.WARN"},
		{caolog.ErrorLevel, "Log.ERROR"},
		{caolog.DPanicLevel, "Log.DPANIC"},
		{caolog.PanicLevel, "Log.PANIC"},
		{caolog.FatalLevel, "Log.FATAL"},
		{caolog.Level(100), "Log.UNKNOWN"},
	}
	
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tracePlugin.GetLogTag(tt.level)
			if result != tt.expected {
				t.Errorf("GetLogTag(%v) = %v, want %v", tt.level, result, tt.expected)
			}
		})
	}
}

// TestOption 测试Option方法
func TestOption(t *testing.T) {
	tracePlugin := NewTrace()
	
	// 创建一个带有trace ID的上下文
	ctx := context.Background()
	
	// 创建一个Details结构
	details := &caolog.Details{
		Level:   caolog.InfoLevel,
		Path:    "test.go:123",
		Message: "test message",
		Value:   []interface{}{"value1", "value2"},
	}
	
	// 调用Option方法
	tracePlugin.Option(ctx, details)
	
	// 验证details被修改（添加了traceID）
	// 注意：由于没有实际的trace span，traceID可能不会被添加
	// 但方法应该不会panic
}

// TestOptionWithErrorLevel 测试Option方法处理错误级别
func TestOptionWithErrorLevel(t *testing.T) {
	tracePlugin := NewTrace()
	
	ctx := context.Background()
	
	details := &caolog.Details{
		Level:   caolog.ErrorLevel,
		Path:    "test.go:123",
		Message: "error message",
		Value:   []interface{}{},
	}
	
	// 调用Option方法
	tracePlugin.Option(ctx, details)
	
	// 验证方法不会panic
}

// TestOptionWithSpanContext 测试Option方法处理span上下文
func TestOptionWithSpanContext(t *testing.T) {
	tracePlugin := NewTrace()
	
	// 创建一个带有span的上下文
	ctx := context.Background()
	
	// 创建一个Details结构
	details := &caolog.Details{
		Level:   caolog.InfoLevel,
		Path:    "test.go:123",
		Message: "test message",
		Value:   []interface{}{},
	}
	
	// 调用Option方法
	tracePlugin.Option(ctx, details)
	
	// 验证方法不会panic
}
// TestOptionWithTraceID 测试Option方法处理带有trace ID的上下文
func TestOptionWithTraceID(t *testing.T) {
	tracePlugin := NewTrace()
	
	// 创建一个带有trace ID的上下文
	// 使用OpenTelemetry创建一个span，需要配置TracerProvider
	ctx := context.Background()
	
	// 创建一个TracerProvider，确保trace ID被生成
	// 需要设置一个采样器，否则trace ID可能不会被生成
	provider := otel.GetTracerProvider()
	tracer := provider.Tracer("test")
	ctx, span := tracer.Start(ctx, "test-span")
	defer span.End()
	
	// 检查span context是否有trace ID
	spanCtx := span.SpanContext()
	t.Logf("HasTraceID: %v", spanCtx.HasTraceID())
	t.Logf("TraceID: %v", spanCtx.TraceID())
	
	// 如果没有trace ID，跳过测试
	if !spanCtx.HasTraceID() {
		t.Skip("Skipping test because span context does not have trace ID")
	}
	
	// 创建一个Details结构
	details := &caolog.Details{
		Level:   caolog.InfoLevel,
		Path:    "test.go:123",
		Message: "test message",
		Value:   []interface{}{},
	}
	
	// 调用Option方法
	tracePlugin.Option(ctx, details)
	
	// 验证details被修改（添加了traceID）
	if len(details.Value) == 0 {
		t.Errorf("Expected details.Value to be modified with traceID, got: %v", details.Value)
	}
	
	// 验证traceID被添加到details.Value中
	foundTraceID := false
	for _, v := range details.Value {
		if str, ok := v.(string); ok && strings.Contains(str, "traceID:") {
			foundTraceID = true
			break
		}
	}
	if !foundTraceID {
		t.Errorf("Expected traceID to be added to details.Value, got: %v", details.Value)
	}
}

// TestOptionWithManualTraceID 测试Option方法处理手动创建的trace ID
func TestOptionWithManualTraceID(t *testing.T) {
	tracePlugin := NewTrace()
	
	// 手动创建一个带有trace ID的上下文
	// 由于OpenTelemetry的默认配置可能不会生成trace ID，
	// 我们手动创建一个span context
	traceIDStr := "0123456789abcdef0123456789abcdef"
	spanIDStr := "0123456789abcdef"
	
	// 创建TraceID和SpanID
	traceID, err := trace.TraceIDFromHex(traceIDStr)
	if err != nil {
		t.Fatalf("Failed to create trace ID: %v", err)
	}
	spanID, err := trace.SpanIDFromHex(spanIDStr)
	if err != nil {
		t.Fatalf("Failed to create span ID: %v", err)
	}
	
	// 创建一个span context
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	})
	
	// 创建一个上下文，包含span context
	ctx := trace.ContextWithSpanContext(context.Background(), spanCtx)
	
	// 创建一个Details结构
	details := &caolog.Details{
		Level:   caolog.InfoLevel,
		Path:    "test.go:123",
		Message: "test message",
		Value:   []interface{}{},
	}
	
	// 调用Option方法
	tracePlugin.Option(ctx, details)
	
	// 验证details被修改（添加了traceID）
	if len(details.Value) == 0 {
		t.Errorf("Expected details.Value to be modified with traceID, got: %v", details.Value)
	}
	
	// 验证traceID被添加到details.Value中
	foundTraceID := false
	for _, v := range details.Value {
		if str, ok := v.(string); ok && strings.Contains(str, "traceID:") {
			foundTraceID = true
			break
		}
	}
	if !foundTraceID {
		t.Errorf("Expected traceID to be added to details.Value, got: %v", details.Value)
	}
}

// TestOptionWithMultipleCalls 测试Option方法多次调用
func TestOptionWithMultipleCalls(t *testing.T) {
	tracePlugin := NewTrace()
	
	ctx := context.Background()
	
	// 多次调用Option
	for i := 0; i < 5; i++ {
		details := &caolog.Details{
			Level:   caolog.InfoLevel,
			Path:    "test.go:123",
			Message: "test message",
			Value:   []interface{}{},
		}
		
		tracePlugin.Option(ctx, details)
	}
}

// TestNewTraceWithCustomProvider 测试使用自定义TracerProvider
func TestNewTraceWithCustomProvider(t *testing.T) {
	// 创建一个自定义的TracerProvider
	provider := otel.GetTracerProvider()
	
	// 验证NewTrace使用全局provider
	tracePlugin := NewTrace()
	if tracePlugin.tracerProvider != provider {
		t.Error("NewTrace should use global tracer provider")
	}
}

// TestTracePluginConcurrency 测试Trace插件并发安全
func TestTracePluginConcurrency(t *testing.T) {
	tracePlugin := NewTrace()
	
	// 并发调用Option
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			ctx := context.Background()
			details := &caolog.Details{
				Level:   caolog.InfoLevel,
				Path:    "test.go:123",
				Message: "test message",
				Value:   []interface{}{},
			}
			tracePlugin.Option(ctx, details)
			done <- true
		}()
	}
	
	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestTracePluginWithCaologIntegration 测试与caolog的集成
func TestTracePluginWithCaologIntegration(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)
	
	logger := caolog.GetLogger()
	
	tracePlugin := NewTrace()
	
	// 使用caolog的With方法添加trace插件
	// 注意：这里需要访问caolog的内部方法，实际使用中应该通过caolog.With
	// 由于测试限制，我们直接调用Option
	ctx := context.Background()
	details := &caolog.Details{
		Level:   caolog.InfoLevel,
		Path:    "test.go:123",
		Message: "test message",
		Value:   []interface{}{},
	}
	
	tracePlugin.Option(ctx, details)
	
	// 验证方法不会panic
	_ = logger
}