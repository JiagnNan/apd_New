package log

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 是项目日志配置，支持日志级别、输出格式和堆栈开关。
type Config struct {
	Level      string `json:"level,optional"`
	Encoding   string `json:"encoding,optional"`
	Stacktrace string `json:"stacktrace,optional"`
}

// 全局日志器和级别开关，业务代码通过 pkg/log 暴露的方法写日志。
var (
	globalLevel  = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	globalLogger = newLogger(Config{
		Level:      getenv("LOG_LEVEL", DebugLevel),
		Encoding:   os.Getenv("LOG_ENCODING"),
		Stacktrace: os.Getenv("LOG_STACKTRACE"),
	})
)

// 初始化全局日志级别。
func init() {
	SetLevel(getenv("LOG_LEVEL", DebugLevel))
}

// 创建带名称的 zap logger。
func NewLogger(name string, cfg Config) *zap.Logger {
	return newLogger(cfg).Named(name)
}

// 设置全局日志配置。
func SetConfig(cfg Config) {
	globalLogger = newLogger(cfg)
	SetLevel(cfg.Level)
}

// 设置全局日志级别。
func SetLevel(level string) {
	if level == "" {
		level = DebugLevel
	}
	_ = globalLevel.UnmarshalText([]byte(strings.ToLower(level)))
}

// 刷新日志缓冲区。
func Sync() {
	_ = globalLogger.Sync()
}

// 强制打印 info 级别日志，不受全局级别过滤影响。
func Printx(ctx context.Context, template string, args ...any) {
	globalLogger.Info(formatMsg(template, args), fieldsFromCtx(ctx)...)
}

// 打印带 context 字段的 debug 日志。
func Debugx(ctx context.Context, template string, args ...any) {
	if enabled(ctx, zapcore.DebugLevel) {
		globalLogger.Debug(formatMsg(template, args), fieldsFromCtx(ctx)...)
	}
}

// 打印带 context 字段的 info 日志。
func Infox(ctx context.Context, template string, args ...any) {
	if enabled(ctx, zapcore.InfoLevel) {
		globalLogger.Info(formatMsg(template, args), fieldsFromCtx(ctx)...)
	}
}

// 打印带 context 字段的 warn 日志。
func Warnx(ctx context.Context, template string, args ...any) {
	if enabled(ctx, zapcore.WarnLevel) {
		globalLogger.Warn(formatMsg(template, args), fieldsFromCtx(ctx)...)
	}
}

// 打印带 context 字段的 error 日志。
func Errorx(ctx context.Context, template string, args ...any) {
	if enabled(ctx, zapcore.ErrorLevel) {
		globalLogger.Error(formatMsg(template, args), fieldsFromCtx(ctx)...)
	}
}

// 打印带 context 字段的 panic 日志。
func Panicx(ctx context.Context, template string, args ...any) {
	globalLogger.Panic(formatMsg(template, args), fieldsFromCtx(ctx)...)
}

// 打印不带 context 字段的 debug 日志。
func Debugf(template string, args ...any) {
	if enabled(context.Background(), zapcore.DebugLevel) {
		globalLogger.Debug(formatMsg(template, args))
	}
}

// 打印不带 context 字段的 info 日志。
func Infof(template string, args ...any) {
	if enabled(context.Background(), zapcore.InfoLevel) {
		globalLogger.Info(formatMsg(template, args))
	}
}

// 打印不带 context 字段的 warn 日志。
func Warnf(template string, args ...any) {
	if enabled(context.Background(), zapcore.WarnLevel) {
		globalLogger.Warn(formatMsg(template, args))
	}
}

// 打印不带 context 字段的 error 日志。
func Errorf(template string, args ...any) {
	if enabled(context.Background(), zapcore.ErrorLevel) {
		globalLogger.Error(formatMsg(template, args))
	}
}

// 打印不带 context 字段的 panic 日志。
func Panicf(template string, args ...any) {
	globalLogger.Panic(formatMsg(template, args))
}

// 默认使用生产级配置
// 默认输出 JSON
// 日志输出到 stdout
// 错误输出到 stderr
// 时间格式用 ISO8601
// 根据 Stacktrace 决定是否打印堆栈
// 模式下日志级别带颜色
// 自动打印调用文件和行号
func newLogger(cfg Config) *zap.Logger {
	encoding := strings.ToLower(cfg.Encoding)
	if encoding == "" {
		encoding = "json"
	}

	zcfg := zap.NewProductionConfig()
	zcfg.Level = globalLevel
	zcfg.Encoding = encoding
	zcfg.OutputPaths = []string{"stdout"}
	zcfg.ErrorOutputPaths = []string{"stderr"}
	zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zcfg.DisableStacktrace = !parseBool(cfg.Stacktrace)
	if encoding == "console" {
		zcfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := zcfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	return logger
}

// 判断指定级别在当前上下文和全局日志级别下是否应该打印。
func enabled(ctx context.Context, level zapcore.Level) bool {
	return IsPrint(ctx) || globalLevel.Enabled(level)
}

// 从 context 中读取日志字段。
func fieldsFromCtx(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}
	return GetFields(ctx)
}

// 1. 去掉日志末尾的换行符
// 2. 如果有参数，用 fmt.Sprintf 格式化
// 3. 如果 template 是空，就直接 fmt.Sprint(args...)
func formatMsg(template string, args []any) string {
	template = strings.TrimSuffix(template, "\n")
	if len(args) == 0 {
		return template
	}
	if template != "" {
		return fmt.Sprintf(template, args...)
	}
	return fmt.Sprint(args...)
}

// 读取环境变量，未设置时返回默认值。
func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// 解析布尔字符串，非法值按 false 处理。
func parseBool(v string) bool {
	ret, _ := strconv.ParseBool(v)
	return ret
}
