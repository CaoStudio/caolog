package caolog

import (
	"context"
	"github.com/bytedance/sonic"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/CaoStudio/caolog/formatter"
)

var (
	logger   *Logger
	logLevel formatter.Level
	writer   io.Writer
)

const (
	logDeep = 4
	tabByte = byte('\t')
)

// Level 日志级别
type Level = formatter.Level

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

type (
	Logger struct {
		writer    io.Writer
		level     Level
		Options   []Option
		Formatter formatter.Formatter
	}

	Details = formatter.Details
)

// Formatter 接口定义日志格式化器
type Formatter = formatter.Formatter

// TextFormatter 文本格式化器（默认格式）
type TextFormatter = formatter.TextFormatter

// JSONFormatter JSON格式化器
type JSONFormatter = formatter.JSONFormatter

type Option func(ctx context.Context, details *formatter.Details)

func init() {
	logLevel = DebugLevel
	logger = &Logger{
		writer:    os.Stdout,
		level:     DebugLevel,
		Formatter: formatter.DefaultFormatter(),
	}
}

// InitLogger 初始化日志
// level: debug,info,warn,error,panic,fatal
func InitLogger(level Level, options ...Option) {
	logLevel = level
	logger = &Logger{
		writer:    os.Stdout,
		level:     level,
		Options:   make([]Option, 0),
		Formatter: formatter.DefaultFormatter(),
	}
	if len(options) > 0 {
		logger.with(options...)
	}
}

// SetWriter 设置日志输出
func SetWriter(w io.Writer) {
	logger.writer = w
	writer = w
}

func GetLogger() *Logger {
	return logger
}

func GetWriter() io.Writer {
	return writer
}

func getValue(v interface{}) string {
	switch v.(type) {
	case float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v.(float32)), 'f', -1, 64)
	case int:
		return strconv.FormatInt(int64(v.(int)), 10)
	case uint:
		return strconv.FormatUint(uint64(v.(uint)), 10)
	case int8:
		return strconv.FormatInt(int64(v.(int8)), 10)
	case uint8:
		return strconv.FormatUint(uint64(v.(uint8)), 10)
	case int16:
		return strconv.FormatInt(int64(v.(int16)), 10)
	case uint16:
		return strconv.FormatUint(uint64(v.(uint16)), 10)
	case int32:
		return strconv.FormatInt(int64(v.(int32)), 10)
	case uint32:
		return strconv.FormatUint(uint64(v.(uint32)), 10)
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	case uint64:
		return strconv.FormatUint(v.(uint64), 10)
	case string:
		return v.(string)
	case []byte:
		return string(v.([]byte))
	case error:
		return v.(error).Error()
	default:
		newValue, err := sonic.Marshal(&v)
		if err != nil {
			return "Log Format Error:" + err.Error()
		}
		return *(*string)(unsafe.Pointer(&newValue))
	}
}

func FormatBufferPool[t any](value ...t) string {
	if value == nil {
		return ""
	}
	bufferLen := 0

	cache := make([]string, 0, len(value))
	for index, v := range value {
		cache = append(cache, getValue(v))
		bufferLen += len(cache[index])
	}

	builder := strings.Builder{}
	builder.Grow(bufferLen + len(value) - 1)
	for index, v := range cache {
		builder.WriteString(v)
		if index < len(value)-1 {
			builder.WriteString("\t")
		}
	}

	return builder.String()
}

func PenultimateIndexByteString(s string, c byte) int {
	var flag = true
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			if flag {
				flag = false
				continue
			}
			return i
		}
	}
	return -1
}

// makeDetails
func (l *Logger) makeDetails(deep int, level Level, value ...interface{}) Details {
	_, file, line, _ := runtime.Caller(deep)

	file = file[PenultimateIndexByteString(file, '/')+1:]
	lineStr := strconv.FormatInt(int64(line), 10)

	msgBytes := make([]byte, len(file)+len(lineStr)+1)
	copy(msgBytes, file)
	msgBytes[len(file)] = ':'
	copy(msgBytes[len(file)+1:], lineStr)
	fileLine := *(*string)(unsafe.Pointer(&msgBytes))

	return Details{
		Level:   level,
		Path:    fileLine,
		Time:    time.Now(),
		Message: FormatBufferPool(value...),
		Value:   value,
	}
}

type output func(msg string)

func With(options ...Option) {
	logger.with(options...)
}

// SetFormatter 设置日志格式化器
func SetFormatter(formatter Formatter) {
	logger.Formatter = formatter
}

// With
func (l *Logger) with(options ...Option) {
	l.Options = append(l.Options, options...)
}

func (l *Logger) withSpan(c context.Context, deep int, level Level, output output, value ...interface{}) {
	// 构建日志详情结构体
	detail := l.makeDetails(deep, level, value...)
	// 遍历options，执行option
	for _, option := range l.Options {
		option(c, &detail)
	}

	// 从buffer池中获取buffer，用于拼接日志详情
	builder := strings.Builder{}
	builder.Grow(max(30, len(detail.Path)+len(detail.Message)+1))

	builder.WriteString(detail.Path)
	for builder.Len() < 30 {
		builder.WriteByte(' ')
	}
	builder.WriteString(detail.Message)

	output(builder.String())
}

func (l *Logger) log(level Level, deep int, args ...interface{}) {
	if l.level > level {
		return
	}

	msg := FormatBufferPool(args...)

	// 获取调用位置
	_, file, line, _ := runtime.Caller(deep)
	file = file[PenultimateIndexByteString(file, '/')+1:]
	lineStr := strconv.FormatInt(int64(line), 10)

	// 构建文件行信息
	fileLine := file + ":" + lineStr

	// 创建日志详情
	details := Details{
		Level:   level,
		Path:    fileLine,
		Time:    time.Now(),
		Message: msg,
		Value:   args,
	}

	// 使用格式化器格式化日志
	output := formatter.FormatDetails(context.Background(), &details, l.Formatter)

	// 写入输出
	l.writer.Write([]byte(output))
}

func (l *Logger) Debug(deep int, args ...interface{}) {
	l.log(DebugLevel, deep, args...)
}

func (l *Logger) Info(deep int, args ...interface{}) {
	l.log(InfoLevel, deep, args...)
}

func (l *Logger) Warn(deep int, args ...interface{}) {
	l.log(WarnLevel, deep, args...)
}

func (l *Logger) Error(deep int, args ...interface{}) {
	l.log(ErrorLevel, deep, args...)
}

func (l *Logger) DPanic(deep int, args ...interface{}) {
	l.log(DPanicLevel, deep, args...)
}

func (l *Logger) Panic(deep int, args ...interface{}) {
	l.log(PanicLevel, deep, args...)
	panic(FormatBufferPool(args...))
}

func (l *Logger) Fatal(deep int, args ...interface{}) {
	l.log(FatalLevel, deep, args...)
	os.Exit(1)
}

// Context-aware logging
func (l *Logger) CDebug(c context.Context, deep int, args ...interface{}) {
	l.log(DebugLevel, deep, args...)
}

func (l *Logger) CInfo(c context.Context, deep int, args ...interface{}) {
	l.log(InfoLevel, deep, args...)
}

func (l *Logger) CWarn(c context.Context, deep int, args ...interface{}) {
	l.log(WarnLevel, deep, args...)
}

func (l *Logger) CError(c context.Context, deep int, args ...interface{}) {
	l.log(ErrorLevel, deep, args...)
}

func (l *Logger) CDPanic(c context.Context, deep int, args ...interface{}) {
	l.log(DPanicLevel, deep, args...)
}

func (l *Logger) CPanic(c context.Context, deep int, args ...interface{}) {
	l.log(PanicLevel, deep, args...)
	panic(FormatBufferPool(args...))
}

func (l *Logger) CFatal(c context.Context, deep int, args ...interface{}) {
	l.log(FatalLevel, deep, args...)
	os.Exit(1)
}

// Package-level functions
func Debug(args ...interface{}) {
	logger.Debug(logDeep, args...)
}

func Info(args ...interface{}) {
	logger.Info(logDeep, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(logDeep, args...)
}

func Error(args ...interface{}) {
	logger.Error(logDeep, args...)
}

func DPanic(args ...interface{}) {
	logger.DPanic(logDeep, args...)
}

func Panic(args ...interface{}) {
	logger.Panic(logDeep, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(logDeep, args...)
}

// Context-aware package-level functions
func CDebug(c context.Context, args ...interface{}) {
	logger.CDebug(c, logDeep, args...)
}

func CInfo(c context.Context, args ...interface{}) {
	logger.CInfo(c, logDeep, args...)
}

func CWarn(c context.Context, args ...interface{}) {
	logger.CWarn(c, logDeep, args...)
}

func CError(c context.Context, args ...interface{}) {
	logger.CError(c, logDeep, args...)
}

func CDPanic(c context.Context, args ...interface{}) {
	logger.CDPanic(c, logDeep, args...)
}

func CPanic(c context.Context, args ...interface{}) {
	logger.CPanic(c, logDeep, args...)
}

func CFatal(c context.Context, args ...interface{}) {
	logger.CFatal(c, logDeep, args...)
}
