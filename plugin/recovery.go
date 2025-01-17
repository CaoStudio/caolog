package plugin

import (
	"context"
	"fmt"
	caolog "github.com/CaoStudio/cao-log"
	"go.uber.org/zap"
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

// Recovery recover掉项目可能出现的panic，并使用zap记录相关日志
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
			r.logger.Logger.Error(
				fmt.Sprintln(
					zap.Any("error", err),
				),
			)
			// If the connection is dead, we can't write a status to it.
			return
		}

		r.logger.Error(
			r.deep,
			fmt.Sprintln(
				"[Recovery from panic]\n",
				zap.Any("error", err).String,
				zap.String("stack", string(debug.Stack())).String,
			),
		)
	}
}

// CRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
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
			r.logger.Logger.Error(
				fmt.Sprintln(
					zap.Any("error", err),
				),
			)
			// If the connection is dead, we can't write a status to it.
			return
		}

		r.logger.CError(
			context.Background(),
			r.deep,
			fmt.Sprintln(
				"[Recovery from panic]\n",
				zap.Any("error", err).String,
				zap.String("stack", string(debug.Stack())).String,
			),
		)
	}
}
