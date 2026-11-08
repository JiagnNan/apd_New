package log

import (
	"context"
	"strings"

	"go.uber.org/zap"
)

// 日志字段 context key，使用私有类型避免和业务 context key 冲突。
type ctxKeyFields struct{}

// 强制打印日志 context key，适用于临时调试或重要链路。
type ctxKeyPrint struct{}

// Field 复用 zap 字段类型，业务层无需直接导入 zap。
type Field = zap.Field

// 常用 zap 字段构造器别名，统一从 pkg/log 暴露。
var (
	String   = zap.String
	Int      = zap.Int
	Int32    = zap.Int32
	Int64    = zap.Int64
	Uint     = zap.Uint
	Uint32   = zap.Uint32
	Uint64   = zap.Uint64
	Float32  = zap.Float32
	Float64  = zap.Float64
	Bool     = zap.Bool
	Any      = zap.Any
	Time     = zap.Time
	Duration = zap.Duration
)

// 向 context 写入字符串日志字段。
func WithString(ctx context.Context, key, value string) context.Context {
	return WithField(ctx, String(key, value))
}

// 向 context 写入 int64 日志字段。
func WithInt64(ctx context.Context, key string, value int64) context.Context {
	return WithField(ctx, Int64(key, value))
}

// 往 ctx 里追加字段。
func WithField(ctx context.Context, fields ...Field) context.Context {
	for i := range fields {
		fields[i].Key = strings.ToLower(fields[i].Key)
	}

	exists, _ := ctx.Value(ctxKeyFields{}).([]Field)
	next := make([]Field, 0, len(exists)+len(fields))
	next = append(next, exists...)
	next = append(next, fields...)
	return context.WithValue(ctx, ctxKeyFields{}, next)
}

// 获取 context 中已绑定的日志字段。
func GetFields(ctx context.Context) []Field {
	fields, _ := ctx.Value(ctxKeyFields{}).([]Field)
	return fields
}

// 判断 context 中是否已绑定指定日志字段。
func HasField(ctx context.Context, key string) bool {
	key = strings.ToLower(key)
	for _, field := range GetFields(ctx) {
		if strings.ToLower(field.Key) == key {
			return true
		}
	}
	return false
}

// 标记当前 context 即使低于日志级别也要打印。
func WithPrint(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyPrint{}, true)
}

// 判断当前 context 是否要求强制打印日志。
func IsPrint(ctx context.Context) bool {
	v, _ := ctx.Value(ctxKeyPrint{}).(bool)
	return v
}

// 复制请求值和 trace 但解除取消信号，供请求结束后仍需执行的异步日志使用。
func Copy(ctx context.Context) context.Context {
	return context.WithoutCancel(ctx)
}
