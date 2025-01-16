package caolog

// CommonLogger fmgo内部通用日志接口
type CommonLogger struct {
	Logger
}

var commonDeep = 3

// NewCommonLogger 创建一个通用日志对象
func NewCommonLogger(logger *Logger) CommonLogger {
	return CommonLogger{*logger}
}

func (l CommonLogger) Debug(args ...interface{}) {
	l.Logger.Debug(commonDeep, args...)
}
func (l CommonLogger) Info(args ...interface{}) {
	l.Logger.Info(commonDeep, args...)
}
func (l CommonLogger) Warn(args ...interface{}) {
	l.Logger.Warn(commonDeep, args...)
}
func (l CommonLogger) Error(args ...interface{}) {
	// 预留错误处理
	l.Logger.Error(commonDeep, args...)
}
func (l CommonLogger) DPanic(args ...interface{}) {
	l.Logger.DPanic(commonDeep, args...)
}
func (l CommonLogger) Panic(args ...interface{}) {
	l.Logger.Panic(commonDeep, args...)
}
func (l CommonLogger) Fatal(args ...interface{}) {
	l.Logger.Fatal(commonDeep, args...)
}
