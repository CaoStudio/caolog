package plugin

import (
	"bytes"
	"errors"
	"net"
	"os"
	"strings"
	"testing"

	caolog "github.com/CaoStudio/caolog"
)

// TestNewRecovery 测试NewRecovery函数
func TestNewRecovery(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)

	logger := caolog.GetLogger()
	recovery := NewRecovery(*logger)
	if recovery.deep != 4 {
		t.Errorf("Expected deep=4, got %d", recovery.deep)
	}
}

// TestWithDeep 测试WithDeep方法
func TestWithDeep(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)

	logger := caolog.GetLogger()
	recovery := NewRecovery(*logger)
	recovery.WithDeep(5)
	if recovery.deep != 5 {
		t.Errorf("Expected deep=5, got %d", recovery.deep)
	}
}

// TestRecoveryNormal 测试Recovery正常恢复
func TestRecoveryNormal(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)

	logger := caolog.GetLogger()
	recovery := NewRecovery(*logger)

	// 模拟panic并恢复
	func() {
		defer recovery.Recovery()
		panic("test panic")
	}()

	// 验证日志被记录
	output := buf.String()
	if !strings.Contains(output, "Recovery from panic") {
		t.Errorf("Expected recovery log, got: %v", output)
	}
	if !strings.Contains(output, "test panic") {
		t.Errorf("Expected panic message in log, got: %v", output)
	}
}

// TestRecoveryBrokenPipe 测试Recovery处理broken pipe错误
func TestRecoveryBrokenPipe(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)

	logger := caolog.GetLogger()
	recovery := NewRecovery(*logger)

	// 模拟broken pipe错误
	opError := &net.OpError{
		Op:  "write",
		Err: &os.SyscallError{Err: errors.New("broken pipe")},
	}

	func() {
		defer recovery.Recovery()
		panic(opError)
	}()

	// 验证日志被记录
	output := buf.String()
	if !strings.Contains(output, "error:") {
		t.Errorf("Expected error log, got: %v", output)
	}
	if !strings.Contains(output, "broken pipe") {
		t.Errorf("Expected 'broken pipe' in log, got: %v", output)
	}
}

// TestRecoveryConnectionReset 测试Recovery处理connection reset错误
func TestRecoveryConnectionReset(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)

	logger := caolog.GetLogger()
	recovery := NewRecovery(*logger)

	// 模拟connection reset错误
	opError := &net.OpError{
		Op:  "write",
		Err: &os.SyscallError{Err: errors.New("connection reset by peer")},
	}

	func() {
		defer recovery.Recovery()
		panic(opError)
	}()

	// 验证日志被记录
	output := buf.String()
	if !strings.Contains(output, "error:") {
		t.Errorf("Expected error log, got: %v", output)
	}
	if !strings.Contains(output, "connection reset by peer") {
		t.Errorf("Expected 'connection reset by peer' in log, got: %v", output)
	}
}

// TestRecoveryErrorType 测试Recovery处理不同错误类型
func TestRecoveryErrorType(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)

	logger := caolog.GetLogger()
	recovery := NewRecovery(*logger)

	// 测试error类型
	func() {
		defer recovery.Recovery()
		panic(errors.New("test error"))
	}()

	output := buf.String()
	if !strings.Contains(output, "test error") {
		t.Errorf("Expected 'test error' in log, got: %v", output)
	}

	// 测试string类型
	buf.Reset()
	func() {
		defer recovery.Recovery()
		panic("string panic")
	}()

	output = buf.String()
	if !strings.Contains(output, "string panic") {
		t.Errorf("Expected 'string panic' in log, got: %v", output)
	}
}

// TestCRecoveryNormal 测试CRecovery正常恢复
func TestCRecoveryNormal(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)

	logger := caolog.GetLogger()
	recovery := NewRecovery(*logger)

	// 模拟panic并恢复
	func() {
		defer recovery.CRecovery()
		panic("test panic")
	}()

	// 验证日志被记录
	output := buf.String()
	if !strings.Contains(output, "Recovery from panic") {
		t.Errorf("Expected recovery log, got: %v", output)
	}
	if !strings.Contains(output, "test panic") {
		t.Errorf("Expected panic message in log, got: %v", output)
	}
}

// TestCRecoveryBrokenPipe 测试CRecovery处理broken pipe错误
func TestCRecoveryBrokenPipe(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)

	logger := caolog.GetLogger()
	recovery := NewRecovery(*logger)

	// 模拟broken pipe错误
	opError := &net.OpError{
		Op:  "write",
		Err: &os.SyscallError{Err: errors.New("broken pipe")},
	}

	func() {
		defer recovery.CRecovery()
		panic(opError)
	}()

	// 验证日志被记录
	output := buf.String()
	if !strings.Contains(output, "error:") {
		t.Errorf("Expected error log, got: %v", output)
	}
	if !strings.Contains(output, "broken pipe") {
		t.Errorf("Expected 'broken pipe' in log, got: %v", output)
	}
}

// TestRecoveryWithDifferentDeep 测试不同deep值
func TestRecoveryWithDifferentDeep(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)

	logger := caolog.GetLogger()
	recovery := NewRecovery(*logger)
	recovery.WithDeep(2)

	// 模拟panic并恢复
	func() {
		defer recovery.Recovery()
		panic("test panic")
	}()

	// 验证日志被记录
	output := buf.String()
	if !strings.Contains(output, "Recovery from panic") {
		t.Errorf("Expected recovery log, got: %v", output)
	}
}

// TestRecoveryNoPanic 测试没有panic的情况
func TestRecoveryNoPanic(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)

	logger := caolog.GetLogger()
	recovery := NewRecovery(*logger)

	// 没有panic
	func() {
		defer recovery.Recovery()
		// 正常执行
	}()

	// 验证没有日志被记录
	output := buf.String()
	if output != "" {
		t.Errorf("Expected no log when no panic, got: %v", output)
	}
}

// TestCRecoveryNoPanic 测试CRecovery没有panic的情况
func TestCRecoveryNoPanic(t *testing.T) {
	var buf bytes.Buffer
	caolog.InitLogger(caolog.DebugLevel)
	caolog.SetWriter(&buf)

	logger := caolog.GetLogger()
	recovery := NewRecovery(*logger)

	// 没有panic
	func() {
		defer recovery.CRecovery()
		// 正常执行
	}()

	// 验证没有日志被记录
	output := buf.String()
	if output != "" {
		t.Errorf("Expected no log when no panic, got: %v", output)
	}
}
