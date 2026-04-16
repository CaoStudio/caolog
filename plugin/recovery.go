package plugin

import (
	"context"
	caolog "github.com/CaoStudio/caolog"
	"net"
	"os"
	"runtime/debug"
	"strings"
)

type Recovery struct {
	logger caolog.Logger
	deep   int
}

// NewRecovery returns a new recovery plugin.
func NewRecovery(caologger caolog.Logger) *Recovery {
	return &Recovery{
		logger: caologger,
		deep:   4,
	}
}

// WithDeep with recovery deep.
func (r *Recovery) WithDeep(deep int) {
	r.deep = deep
}

// Recovery recover掉项目可能出现的panic，并记录相关日志
func (r *Recovery) Recovery() {
	if err := recover(); err != nil {
		// Check for a broken connection, as it is not really a
		// condition that warrants a panic stack trace.
		var brokenPipe bool
		if ne, ok := err.(*net.OpError); ok {
			if se, ok := ne.Err.(*os.SyscallError); ok {
				if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
					brokenPipe = true
				}
			}
		}
		if brokenPipe {
			var builder strings.Builder
			builder.WriteString("error: ")
			if e, ok := err.(error); ok {
				builder.WriteString(e.Error())
			} else {
				builder.WriteString(caolog.FormatBufferPool(err))
			}
			r.logger.Error(r.deep, builder.String())
			// If the connection is dead, we can't write a status to it.
			return
		}

		var builder strings.Builder
		builder.WriteString("[Recovery from panic]\nerror: ")
		if e, ok := err.(error); ok {
			builder.WriteString(e.Error())
		} else {
			builder.WriteString(caolog.FormatBufferPool(err))
		}
		builder.WriteString("\nstack: ")
		builder.WriteString(string(debug.Stack()))
		r.logger.Error(r.deep, builder.String())
	}
}

// CRecovery recover掉项目可能出现的panic，并记录相关日志
func (r *Recovery) CRecovery() {
	if err := recover(); err != nil {
		// Check for a broken connection, as it is not really a
		// condition that warrants a panic stack trace.
		var brokenPipe bool
		if ne, ok := err.(*net.OpError); ok {
			if se, ok := ne.Err.(*os.SyscallError); ok {
				if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
					brokenPipe = true
				}
			}
		}
		if brokenPipe {
			var builder strings.Builder
			builder.WriteString("error: ")
			if e, ok := err.(error); ok {
				builder.WriteString(e.Error())
			} else {
				builder.WriteString(caolog.FormatBufferPool(err))
			}
			r.logger.CError(context.Background(), r.deep, builder.String())
			// If the connection is dead, we can't write a status to it.
			return
		}

		var builder strings.Builder
		builder.WriteString("[Recovery from panic]\nerror: ")
		if e, ok := err.(error); ok {
			builder.WriteString(e.Error())
		} else {
			builder.WriteString(caolog.FormatBufferPool(err))
		}
		builder.WriteString("\nstack: ")
		builder.WriteString(string(debug.Stack()))
		r.logger.CError(context.Background(), r.deep, builder.String())
	}
}
