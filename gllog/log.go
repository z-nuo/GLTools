package gllog

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

// Format 表示日志输出格式。
type Format string

const (
	// FormatJSON 输出 JSON 格式日志。
	FormatJSON Format = "json"
	// FormatText 输出文本格式日志。
	FormatText Format = "text"
)

// Config 表示日志配置。
type Config struct {
	Output    io.Writer
	Format    Format
	Level     slog.Level
	AddSource bool
}

// New 根据配置创建 slog.Logger。
func New(cfg Config) (*slog.Logger, error) {
	output := cfg.Output
	if output == nil {
		output = os.Stdout
	}

	opts := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     cfg.Level,
	}

	switch cfg.Format {
	case FormatJSON:
		return slog.New(slog.NewJSONHandler(output, opts)), nil
	case FormatText:
		return slog.New(slog.NewTextHandler(output, opts)), nil
	default:
		return nil, fmt.Errorf("unsupported log format: %s", cfg.Format)
	}
}

// SetDefault 设置默认日志器。
func SetDefault(logger *slog.Logger) {
	slog.SetDefault(logger)
}
