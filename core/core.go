package core

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var builderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
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

// Logger 结构体
type Logger struct {
	writer io.Writer
	level  Level
}

// Event 结构体
type Event struct {
	logger *Logger
	level  Level
	msg    string
	fields []interface{}
	exit   func(string)
}

// NewLogger 创建一个新的 Logger
func NewLogger(writer io.Writer, level Level) *Logger {
	if writer == nil {
		writer = os.Stdout
	}
	return &Logger{
		writer: writer,
		level:  level,
	}
}

// WithLevel 创建一个新的 Logger，使用指定的日志级别
func (l *Logger) WithLevel(level Level) *Logger {
	return NewLogger(l.writer, level)
}

// WithWriter 创建一个新的 Logger，使用指定的 writer
func (l *Logger) WithWriter(writer io.Writer) *Logger {
	return NewLogger(writer, l.level)
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Debug() *Event {
	return l.newEvent(DebugLevel, nil)
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Info() *Event {
	return l.newEvent(InfoLevel, nil)
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Warn() *Event {
	return l.newEvent(WarnLevel, nil)
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Error() *Event {
	return l.newEvent(ErrorLevel, nil)
}

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Err(err error) *Event {
	if err != nil {
		return l.Error().Err(err)
	}
	return l.Info()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method, which terminates the program immediately.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Fatal() *Event {
	return l.newEvent(FatalLevel, func(msg string) {
		if closer, ok := l.writer.(io.Closer); ok {
			// Close the writer to flush any buffered message. Otherwise the message
			// will be lost as os.Exit() terminates the program immediately.
			closer.Close()
		}
		os.Exit(1)
	})
}

// Panic starts a new message with panic level. The panic() function
// is called by the Msg method, which stops the ordinary flow of a goroutine.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Panic() *Event {
	return l.newEvent(PanicLevel, func(msg string) { panic(msg) })
}

// newEvent creates a new event with the given level and optional exit function.
func (l *Logger) newEvent(level Level, exit func(string)) *Event {
	return &Event{
		logger: l,
		level:  level,
		exit:   exit,
	}
}

// Msg sends the event with the supplied message.
func (e *Event) Msg(msg string) {
	if e.logger.level > e.level {
		return
	}
	e.msg = msg
	e.output()
	if e.exit != nil {
		e.exit(msg)
	}
}

// Err adds an error field to the event.
func (e *Event) Err(err error) *Event {
	if err != nil {
		e.fields = append(e.fields, "error", err.Error())
	}
	return e
}

// Str adds a string field to the event.
func (e *Event) Str(key, value string) *Event {
	e.fields = append(e.fields, key, value)
	return e
}

// Int adds an int field to the event.
func (e *Event) Int(key string, value int) *Event {
	e.fields = append(e.fields, key, strconv.Itoa(value))
	return e
}

// Bool adds a bool field to the event.
func (e *Event) Bool(key string, value bool) *Event {
	e.fields = append(e.fields, key, strconv.FormatBool(value))
	return e
}

// Float64 adds a float64 field to the event.
func (e *Event) Float64(key string, value float64) *Event {
	e.fields = append(e.fields, key, strconv.FormatFloat(value, 'f', -1, 64))
	return e
}

// Float32 adds a float32 field to the event.
func (e *Event) Float32(key string, value float32) *Event {
	e.fields = append(e.fields, key, strconv.FormatFloat(float64(value), 'f', -1, 32))
	return e
}

// Int32 adds an int32 field to the event.
func (e *Event) Int32(key string, value int32) *Event {
	e.fields = append(e.fields, key, strconv.FormatInt(int64(value), 10))
	return e
}

// Int64 adds an int64 field to the event.
func (e *Event) Int64(key string, value int64) *Event {
	e.fields = append(e.fields, key, strconv.FormatInt(value, 10))
	return e
}

// Uint adds a uint field to the event.
func (e *Event) Uint(key string, value uint) *Event {
	e.fields = append(e.fields, key, strconv.FormatUint(uint64(value), 10))
	return e
}

// Uint32 adds a uint32 field to the event.
func (e *Event) Uint32(key string, value uint32) *Event {
	e.fields = append(e.fields, key, strconv.FormatUint(uint64(value), 10))
	return e
}

// Uint64 adds a uint64 field to the event.
func (e *Event) Uint64(key string, value uint64) *Event {
	e.fields = append(e.fields, key, strconv.FormatUint(value, 10))
	return e
}

// Dur adds a time.Duration field to the event.
func (e *Event) Dur(key string, value time.Duration) *Event {
	e.fields = append(e.fields, key, value.String())
	return e
}

// Any adds an arbitrary field to the event.
func (e *Event) Any(key string, value interface{}) *Event {
	e.fields = append(e.fields, key, getValue(value))
	return e
}

// Time adds a time field to the event.
func (e *Event) Time(key string, value time.Time) *Event {
	e.fields = append(e.fields, key, value.Format(time.RFC3339))
	return e
}

// getValue converts a value to a string representation.
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
		// For other types, use fmt.Sprintf as fallback
		return fmt.Sprintf("%v", v)
	}
}

// output writes the event to the logger's writer.
func (e *Event) output() {
	// 获取调用位置
	_, file, line, _ := runtime.Caller(2)
	file = file[penultimateIndexByteString(file, '/')+1:]
	lineStr := strconv.FormatInt(int64(line), 10)

	// 从池中获取 builder
	builder := builderPool.Get().(*strings.Builder)
	builder.Reset()
	defer builderPool.Put(builder)

	// 预分配空间
	builder.Grow(100)

	// 时间
	builder.WriteString("[")
	builder.WriteString(time.Now().Format("2006-01-02 - 15:04:05"))
	builder.WriteString("] ")

	// 级别
	builder.WriteString("[")
	builder.WriteString(e.level.String())
	builder.WriteString("] ")

	// 调用位置
	builder.WriteString(file)
	builder.WriteString(":")
	builder.WriteString(lineStr)
	builder.WriteString(" ")

	// 消息
	builder.WriteString(e.msg)

	// 字段
	for i := 0; i < len(e.fields); i += 2 {
		builder.WriteString(" ")
		builder.WriteString(e.fields[i].(string))
		builder.WriteString("=")
		if i+1 < len(e.fields) {
			builder.WriteString(e.fields[i+1].(string))
		}
	}

	builder.WriteString("\n")
	e.logger.writer.Write([]byte(builder.String()))
}

// penultimateIndexByteString 返回字符串中倒数第二个指定字节的索引。
func penultimateIndexByteString(s string, c byte) int {
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

// Package-level functions
var logger = NewLogger(os.Stdout, DebugLevel)

func Debug() *Event {
	return logger.Debug()
}

func Info() *Event {
	return logger.Info()
}

func Warn() *Event {
	return logger.Warn()
}

func Error() *Event {
	return logger.Error()
}

func Fatal() *Event {
	return logger.Fatal()
}

func Panic() *Event {
	return logger.Panic()
}

func Err(err error) *Event {
	return logger.Err(err)
}

// WithLogger sets the global logger instance.
func WithLogger(l *Logger) {
	logger = l
}

// WithLevel sets the global logger level.
func WithLevel(level Level) {
	logger = logger.WithLevel(level)
}

// WithWriter sets the global logger writer.
func WithWriter(writer io.Writer) {
	logger = logger.WithWriter(writer)
}
