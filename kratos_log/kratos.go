package kratoslog

import (
	"context"
	"fmt"
	caolog "github.com/CaoStudio/cao-log"
	"github.com/go-kratos/kratos/v2/errors"
	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"time"
)

type Logger struct {
	*caolog.Logger
}

// NewKratosLogger	returns a new kratos logger
func NewKratosLogger(logger *caolog.Logger) kratoslog.Logger {
	return Logger{
		Logger: logger,
	}
}

var kratosDeep = 4

func (l Logger) Log(level kratoslog.Level, msg ...interface{}) error {
	switch level {
	case kratoslog.LevelDebug:
		l.Logger.Debug(kratosDeep, msg...)
	case kratoslog.LevelInfo:
		l.Logger.Info(kratosDeep, msg...)
	case kratoslog.LevelWarn:
		l.Logger.Warn(kratosDeep, msg...)
	case kratoslog.LevelError:
		l.Logger.Error(kratosDeep, msg...)
	case kratoslog.LevelFatal:
		l.Logger.Fatal(kratosDeep, msg...)
	}
	return nil
}

// KratosServer is an server logging middleware.
func KratosServer(logger kratoslog.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			level, _ := extractError(err)
			_ = kratoslog.WithContext(ctx, logger).Log(level,
				"msg",
				fmt.Sprintf("[Kratos] %4d | %12s | %8s | %7s | %s\t%s",
					code,
					time.Since(startTime).String(),
					"server",
					kind,
					operation,
					reason),
			)
			return
		}
	}
}

// extractError returns the string of the error
func extractError(err error) (kratoslog.Level, string) {
	if err != nil {
		return kratoslog.LevelError, fmt.Sprintf("%+v", err.Error())
	}
	return kratoslog.LevelInfo, ""
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}
