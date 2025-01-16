package caolog

import (
	"context"
	"github.com/bytedance/sonic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var (
	logger *Logger
	Level  zapcore.Level
	writer io.Writer
)

const (
	logDeep = 4
	tabByte = byte('\t')
)

type (
	Logger struct {
		*zap.Logger
		Options []Option
	}

	Details struct {
		Level zapcore.Level `json:"level,omitempty"`
		// 调用log的文件路径
		Path string `json:"path,omitempty"`
		// Time holds the value of the "time" field.
		Time time.Time `json:"time,omitempty"`
		// 日志内容
		Message string `json:"message,omitempty"`
		// 内容列表
		Value []interface{}
	}
)

const (
	DebugLevel  = zapcore.DebugLevel
	InfoLevel   = zapcore.InfoLevel
	WarnLevel   = zapcore.WarnLevel
	ErrorLevel  = zapcore.ErrorLevel
	DPanicLevel = zapcore.DPanicLevel
	PanicLevel  = zapcore.PanicLevel
	FatalLevel  = zapcore.FatalLevel
)

func init() {
	Level = DebugLevel
	logger = &Logger{
		Logger: zap.NewExample(),
	}
}

type Option func(ctx context.Context, details *Details)

// InitLogger
// level: debug,info,warn,error,panic,fatal
func InitLogger(level zapcore.Level, options ...Option) {
	Level = level
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}

	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    customLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("[2006-01-02 - 15:04:05]"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	core := zapcore.NewCore(encoder, zapcore.AddSync(io.Discard), level)

	logger = &Logger{
		Logger:  zap.New(core),
		Options: make([]Option, 0),
	}
	if len(options) > 0 {
		logger.with(options...)
	}
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
		//return FormatInt(int64(v.(int)))
		return strconv.FormatInt(int64(v.(int)), 10)
	case uint:
		//return FormatUint(uint64(v.(uint)))
		return strconv.FormatUint(uint64(v.(uint)), 10)
	case int8:
		//return FormatInt(int64(v.(int8)))
		return strconv.FormatInt(int64(v.(int8)), 10)
	case uint8:
		//return FormatUint(uint64(v.(uint8)))
		return strconv.FormatUint(uint64(v.(uint8)), 10)
	case int16:
		//return FormatInt(int64(v.(int16)))
		return strconv.FormatInt(int64(v.(int16)), 10)
	case uint16:
		//return FormatUint(uint64(v.(uint16)))
		return strconv.FormatUint(uint64(v.(uint16)), 10)
	case int32:
		//return FormatInt(int64(v.(int32)))
		return strconv.FormatInt(int64(v.(int32)), 10)
	case uint32:
		//return FormatUint(uint64(v.(uint32)))
		return strconv.FormatUint(uint64(v.(uint32)), 10)
	case int64:
		//return FormatInt(v.(int64))
		return strconv.FormatInt(v.(int64), 10)
	case uint64:
		//return FormatUint(v.(uint64))
		return strconv.FormatUint(v.(uint64), 10)
	case string:
		return v.(string)
	case []byte:
		return string(v.([]byte))
	case error:
		return v.(error).Error()
	//case []int32:
	//	return FormatBufferPool(v)
	default:
		//newValue, err := json.Marshal(v)
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
		if index <= len(value)-1 {
			builder.WriteString("\t")
			//builder.Write([]byte{tabByte})
		}
	}

	//msgBytes := make([]byte, bufferLen+len(value)-1)
	//flagNum := 0
	//for i, s := range cache {
	//	copy(msgBytes[flagNum:], s)
	//	flagNum += len(s)
	//	if i < len(cache)-1 {
	//		msgBytes[flagNum] = tabByte
	//		flagNum++
	//	}
	//}

	return builder.String()
	//return *(*string)(unsafe.Pointer(&msgBytes))
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

//func MakeDetails(deep int, level zapcore.Level, value ...interface{}) Details {
//	return logger.makeDetails(deep, level, value...)
//}

// makeDetails
func (l *Logger) makeDetails(deep int, level zapcore.Level, value ...interface{}) Details {
	_, file, line, _ := runtime.Caller(deep)

	file = file[PenultimateIndexByteString(file, '/')+1:]
	lineStr := FormatInt(int64(line))

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

type output func(msg string, fields ...zap.Field)

func With(options ...Option) {
	logger.with(options...)
}

// With
func (l *Logger) with(options ...Option) {
	l.Options = append(l.Options, options...)
}

func (l *Logger) withSpan(c context.Context, deep int, level zapcore.Level, output output, value ...interface{}) {

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
	//_ = builder.String()
	//println(builder.Len())

}

func (l *Logger) CDebug(c context.Context, deep int, args ...interface{}) {
	if Level > DebugLevel {
		return
	}
	l.withSpan(c, deep, DebugLevel, l.Logger.Debug, args...)
}
func (l *Logger) CInfo(c context.Context, deep int, args ...interface{}) {
	if Level > InfoLevel {
		return
	}
	l.withSpan(c, deep, InfoLevel, l.Logger.Info, args...)
}
func (l *Logger) CWarn(c context.Context, deep int, args ...interface{}) {
	if Level > WarnLevel {
		return
	}
	l.withSpan(c, deep, WarnLevel, l.Logger.Warn, args...)
}
func (l *Logger) CError(c context.Context, deep int, args ...interface{}) {
	if Level > ErrorLevel {
		return
	}
	l.withSpan(c, deep, ErrorLevel, l.Logger.Error, args...)
}
func (l *Logger) CDPanic(c context.Context, deep int, args ...interface{}) {
	if Level > DPanicLevel {
		return
	}
	l.withSpan(c, deep, DPanicLevel, l.Logger.DPanic, args...)
}
func (l *Logger) CPanic(c context.Context, deep int, args ...interface{}) {
	if Level > PanicLevel {
		return
	}
	l.withSpan(c, deep, PanicLevel, l.Logger.Panic, args...)
}
func (l *Logger) CFatal(c context.Context, deep int, args ...interface{}) {
	l.withSpan(c, deep, FatalLevel, l.Logger.Fatal, args...)
}

func (l *Logger) Debug(deep int, args ...interface{}) {
	l.CDebug(context.Background(), deep, args)
}
func (l *Logger) Info(deep int, args ...interface{}) {
	l.CInfo(context.Background(), deep, args...)
}
func (l *Logger) Warn(deep int, args ...interface{}) {
	l.CWarn(context.Background(), deep, args...)
}
func (l *Logger) Error(deep int, args ...interface{}) {
	l.CError(context.Background(), deep, args...)
}
func (l *Logger) DPanic(deep int, args ...interface{}) {
	l.CDebug(context.Background(), deep, args...)
}
func (l *Logger) Panic(deep int, args ...interface{}) {
	l.CPanic(context.Background(), deep, args...)
}
func (l *Logger) Fatal(deep int, args ...interface{}) {
	l.CFatal(context.Background(), deep, args...)
}

func Debug(c context.Context, args ...interface{}) {
	logger.CDebug(c, logDeep, args...)
}
func Info(c context.Context, args ...interface{}) {
	logger.CInfo(c, logDeep, args...)
}
func Warn(c context.Context, args ...interface{}) {
	logger.CWarn(c, logDeep, args...)
}
func Error(c context.Context, args ...interface{}) {
	logger.CError(c, logDeep, args...)
}
func Panic(c context.Context, args ...interface{}) {
	logger.CPanic(c, logDeep, args...)
}
func DPanic(c context.Context, args ...interface{}) {
	logger.CDPanic(c, logDeep, args...)
}
func Fatal(c context.Context, args ...interface{}) {
	logger.CFatal(c, logDeep, args...)
}
