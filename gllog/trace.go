package gllog

import (
	"context"

	"go.uber.org/zap"
)

type traceContextKey struct{}

type traceContext struct {
	traceID string
	spanID  string
}

// WithTrace 将 trace_id 和 span_id 写入 context，便于后续日志自动携带链路字段。
//
// traceID 通常表示一次请求或一次跨服务调用链路的全局 ID。
// spanID 通常表示当前服务、当前方法或当前步骤的局部 ID。
func WithTrace(ctx context.Context, traceID string, spanID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, traceContextKey{}, traceContext{
		traceID: traceID,
		spanID:  spanID,
	})
}

// TraceID 从 context 中读取 trace_id；不存在时返回空字符串。
func TraceID(ctx context.Context) string {
	trace := traceFromContext(ctx)
	return trace.traceID
}

// SpanID 从 context 中读取 span_id；不存在时返回空字符串。
func SpanID(ctx context.Context) string {
	trace := traceFromContext(ctx)
	return trace.spanID
}

// TraceFields 将 context 中的链路信息转换为 zap 字段。
//
// 空 trace_id 或空 span_id 会被自动忽略，避免日志中出现无意义字段。
func TraceFields(ctx context.Context) []zap.Field {
	trace := traceFromContext(ctx)
	fields := make([]zap.Field, 0, 2)
	if trace.traceID != "" {
		fields = append(fields, zap.String("trace_id", trace.traceID))
	}
	if trace.spanID != "" {
		fields = append(fields, zap.String("span_id", trace.spanID))
	}
	return fields
}

// WithContext 返回带有 context 链路字段的 zap.Logger。
//
// logger 传 nil 时使用默认日志器，context 中没有链路字段时原样返回 logger。
func WithContext(ctx context.Context, logger *zap.Logger) *zap.Logger {
	if logger == nil {
		logger = L()
	}
	fields := TraceFields(ctx)
	if len(fields) == 0 {
		return logger
	}
	return logger.With(fields...)
}

// DebugContext 使用默认日志器输出 debug 级别日志，并自动携带 context 链路字段。
func DebugContext(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx, L()).Debug(msg, fields...)
}

// InfoContext 使用默认日志器输出 info 级别日志，并自动携带 context 链路字段。
func InfoContext(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx, L()).Info(msg, fields...)
}

// WarnContext 使用默认日志器输出 warn 级别日志，并自动携带 context 链路字段。
func WarnContext(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx, L()).Warn(msg, fields...)
}

// ErrorContext 使用默认日志器输出 error 级别日志，并自动携带 context 链路字段。
//
// err 非 nil 时会被写入 zap.Error(err)，业务字段仍可通过 fields 继续追加。
func ErrorContext(ctx context.Context, msg string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append([]zap.Field{zap.Error(err)}, fields...)
	}
	WithContext(ctx, L()).Error(msg, fields...)
}

func traceFromContext(ctx context.Context) traceContext {
	if ctx == nil {
		return traceContext{}
	}
	trace, ok := ctx.Value(traceContextKey{}).(traceContext)
	if !ok {
		return traceContext{}
	}
	return trace
}
