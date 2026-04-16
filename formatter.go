package caolog

import (
	"context"

	"github.com/CaoStudio/caolog/formatter"
)

// FormatDetails 格式化日志详情
func FormatDetails(ctx context.Context, details *formatter.Details, f formatter.Formatter) string {
	if f == nil {
		f = formatter.DefaultFormatter()
	}
	return f.Format(ctx, details)
}
