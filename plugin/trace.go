package plugin

import (
	"context"
	"errors"
	caolog "gitee.com/cao_5/cao-log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zapcore"
	"strings"
)

type Trace struct {
	tracerProvider trace.TracerProvider
	tracer         trace.Tracer
}

// NewTrace returns a new trace plugin.
func NewTrace() *Trace {
	provider := otel.GetTracerProvider()
	return &Trace{
		tracerProvider: provider,
		tracer:         provider.Tracer("log"),
	}
}

func (t *Trace) GetLogTag(in zapcore.Level) string {
	//DebugLevel:  "Log.DEBUG",
	//InfoLevel:   "Log.INFO",
	//WarnLevel:   "Log.WARN",
	//ErrorLevel:  "Log.ERROR",
	//DPanicLevel: "Log.DPANIC",
	//PanicLevel:  "Log.PANIC",
	//FatalLevel:  "Log.FATAL",
	switch in {
	case zapcore.DebugLevel:
		return "Log.DEBUG"
	case zapcore.InfoLevel:
		return "Log.INFO"
	case zapcore.WarnLevel:
		return "Log.WARN"
	case zapcore.ErrorLevel:
		return "Log.ERROR"
	case zapcore.DPanicLevel:
		return "Log.DPANIC"
	case zapcore.PanicLevel:
		return "Log.PANIC"
	case zapcore.FatalLevel:
		return "Log.FATAL"
	default:
		return "Log.UNKNOWN"
	}
}

func (t *Trace) Option(ctx context.Context, details *caolog.Details) {
	_, span := t.tracer.Start(ctx, t.GetLogTag(details.Level))
	defer span.End()

	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.HasTraceID() {
		return
	}

	//traceMsg := caolog.FormatBufferPool(details.Value...)

	attrs := make([]attribute.KeyValue, 2)
	attrs[0] = attribute.String("path", details.Path)
	if details.Level >= zapcore.ErrorLevel {
		err := errors.New(details.Message)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		attrs[1] = attribute.String("value", details.Message)
	}

	traceID := spanCtx.TraceID()
	details.Value = append([]interface{}{"traceID:", traceID.String()}, details.Value...)
	builder := strings.Builder{}
	builder.Grow(len(details.Message) + len(traceID.String()) + 1)
	builder.WriteString(traceID.String())
	builder.WriteString("\t")
	builder.WriteString(details.Message)
	details.Message = builder.String()

	span.SetAttributes(attrs...)
}
