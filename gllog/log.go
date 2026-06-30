package gllog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Format 表示日志输出格式。
type Format string

const (
	// FormatJSON 输出 JSON 格式日志，适合生产环境和日志采集系统。
	FormatJSON Format = "json"
	// FormatConsole 输出可读性更好的控制台格式日志。
	FormatConsole Format = "console"
	// FormatText 是 FormatConsole 的兼容别名。
	FormatText Format = FormatConsole
)

// Level 表示日志级别。
type Level string

const (
	// LevelDebug 表示调试日志级别。
	LevelDebug Level = "debug"
	// LevelInfo 表示信息日志级别。
	LevelInfo Level = "info"
	// LevelWarn 表示警告日志级别。
	LevelWarn Level = "warn"
	// LevelError 表示错误日志级别。
	LevelError Level = "error"
)

// Rotate 表示日志文件按时间切分的粒度。
type Rotate string

const (
	// RotateNone 表示不按时间切分日志文件。
	RotateNone Rotate = ""
	// RotateHourly 表示按小时切分日志文件。
	RotateHourly Rotate = "hourly"
	// RotateDaily 表示按天切分日志文件。
	RotateDaily Rotate = "daily"
)

// Config 表示 zap 日志配置。
type Config struct {
	// Output 表示自定义日志输出目标，常用于测试或接入外部 writer。
	Output io.Writer
	// FilePath 表示日志文件路径；设置后支持按 Rotate 写入文件。
	FilePath string
	// Format 表示日志输出格式。
	Format Format
	// Level 表示日志级别，空值默认使用 info。
	Level Level
	// Rotate 表示日志文件按小时或按天切分，空值表示不切分。
	Rotate Rotate
	// AddCaller 表示是否记录调用方文件和行号。
	AddCaller bool
	// AddStacktrace 表示是否为 error 及以上级别记录堆栈。
	AddStacktrace bool
	// AlsoStdout 表示写入文件时是否同时输出到标准输出。
	AlsoStdout bool
}

var (
	defaultMu     sync.RWMutex
	defaultLogger = zap.NewNop()
)

// New 根据配置创建 zap.Logger。
func New(cfg Config) (*zap.Logger, error) {
	encoder, err := newEncoder(cfg.Format)
	if err != nil {
		return nil, err
	}
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	writer, err := newWriteSyncer(cfg)
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(encoder, writer, level)
	options := []zap.Option{zap.ErrorOutput(zapcore.Lock(os.Stderr))}
	if cfg.AddCaller {
		options = append(options, zap.AddCaller())
	}
	if cfg.AddStacktrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}
	return zap.New(core, options...), nil
}

// SetDefault 设置默认 zap.Logger。
func SetDefault(logger *zap.Logger) {
	if logger == nil {
		logger = zap.NewNop()
	}
	defaultMu.Lock()
	defer defaultMu.Unlock()
	defaultLogger = logger
	zap.ReplaceGlobals(logger)
}

// L 返回当前默认 zap.Logger。
func L() *zap.Logger {
	defaultMu.RLock()
	defer defaultMu.RUnlock()
	return defaultLogger
}

// S 返回当前默认 zap.SugaredLogger。
func S() *zap.SugaredLogger {
	return L().Sugar()
}

// Sync 刷新默认日志器的缓冲数据。
func Sync() error {
	return L().Sync()
}

func newEncoder(format Format) (zapcore.Encoder, error) {
	if format == "" {
		format = FormatJSON
	}
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeDuration = zapcore.StringDurationEncoder
	cfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	switch format {
	case FormatJSON:
		return zapcore.NewJSONEncoder(cfg), nil
	case FormatConsole:
		return zapcore.NewConsoleEncoder(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported log format: %s", format)
	}
}

func parseLevel(level Level) (zapcore.Level, error) {
	if level == "" {
		level = LevelInfo
	}
	var parsed zapcore.Level
	if err := parsed.UnmarshalText([]byte(level)); err != nil {
		return zapcore.InfoLevel, fmt.Errorf("unsupported log level: %s", level)
	}
	return parsed, nil
}

func newWriteSyncer(cfg Config) (zapcore.WriteSyncer, error) {
	var writers []zapcore.WriteSyncer
	if cfg.Output != nil {
		writers = append(writers, zapcore.AddSync(cfg.Output))
	}
	if cfg.FilePath != "" {
		writer := &timeRotateWriter{
			basePath: cfg.FilePath,
			rotate:   cfg.Rotate,
			now:      time.Now,
		}
		writers = append(writers, writer)
		if cfg.AlsoStdout {
			writers = append(writers, zapcore.Lock(os.Stdout))
		}
	}
	if len(writers) == 0 {
		writers = append(writers, zapcore.Lock(os.Stdout))
	}
	return zapcore.NewMultiWriteSyncer(writers...), nil
}

type timeRotateWriter struct {
	mu       sync.Mutex
	basePath string
	rotate   Rotate
	now      func() time.Time
	file     *os.File
	path     string
}

func (w *timeRotateWriter) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	file, err := w.currentFile()
	if err != nil {
		return 0, err
	}
	return file.Write(data)
}

func (w *timeRotateWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file == nil {
		return nil
	}
	return w.file.Sync()
}

func (w *timeRotateWriter) currentFile() (*os.File, error) {
	path := rotateFileName(w.basePath, w.rotate, w.now())
	if w.file != nil && w.path == path {
		return w.file, nil
	}
	if w.file != nil {
		_ = w.file.Close()
		w.file = nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	w.file = file
	w.path = path
	return file, nil
}

func rotateFileName(basePath string, rotate Rotate, now time.Time) string {
	if rotate == RotateNone {
		return basePath
	}
	ext := filepath.Ext(basePath)
	prefix := strings.TrimSuffix(basePath, ext)
	switch rotate {
	case RotateHourly:
		return prefix + "-" + now.Format("2006010215") + ext
	case RotateDaily:
		return prefix + "-" + now.Format("20060102") + ext
	default:
		return basePath
	}
}
