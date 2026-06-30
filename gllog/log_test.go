package gllog

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestNewCreatesJSONZapLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	logger, err := New(Config{Output: buf, Format: FormatJSON, Level: LevelInfo})
	if err != nil {
		t.Fatal(err)
	}

	logger.Info("hello", zap.String("name", "gltools"))

	got := buf.String()
	if !strings.Contains(got, `"msg":"hello"`) {
		t.Fatalf("log output = %s", got)
	}
	if !strings.Contains(got, `"name":"gltools"`) {
		t.Fatalf("log output = %s", got)
	}
}

func TestNewCreatesConsoleZapLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	logger, err := New(Config{Output: buf, Format: FormatConsole, Level: LevelInfo})
	if err != nil {
		t.Fatal(err)
	}

	logger.Info("hello", zap.String("name", "gltools"))

	got := buf.String()
	if !strings.Contains(got, "hello") {
		t.Fatalf("log output = %s", got)
	}
	if !strings.Contains(got, "gltools") {
		t.Fatalf("log output = %s", got)
	}
}

func TestNewRejectsUnsupportedFormat(t *testing.T) {
	_, err := New(Config{Output: new(bytes.Buffer), Format: Format("xml")})
	if err == nil {
		t.Fatal("New() error = nil, want error")
	}
}

func TestSetDefaultUsesProvidedZapLogger(t *testing.T) {
	prev := L()
	t.Cleanup(func() {
		SetDefault(prev)
	})

	buf := new(bytes.Buffer)
	logger, err := New(Config{Output: buf, Format: FormatJSON, Level: LevelInfo})
	if err != nil {
		t.Fatal(err)
	}

	SetDefault(logger)
	L().Info("default logger")
	S().Infow("sugared logger", "name", "gltools")

	got := buf.String()
	if !strings.Contains(got, `"msg":"default logger"`) {
		t.Fatalf("default log output = %s", got)
	}
	if !strings.Contains(got, `"msg":"sugared logger"`) {
		t.Fatalf("default log output = %s", got)
	}
}

func TestRotateFileName(t *testing.T) {
	now := time.Date(2026, 6, 30, 15, 4, 5, 0, time.Local)

	tests := []struct {
		name   string
		rotate Rotate
		want   string
	}{
		{name: "none", rotate: RotateNone, want: filepath.Join("logs", "app.log")},
		{name: "hourly", rotate: RotateHourly, want: filepath.Join("logs", "app-2026063015.log")},
		{name: "daily", rotate: RotateDaily, want: filepath.Join("logs", "app-20260630.log")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rotateFileName(filepath.Join("logs", "app.log"), tt.rotate, now)
			if got != tt.want {
				t.Fatalf("rotateFileName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewWritesDailyRotatedFile(t *testing.T) {
	dir := t.TempDir()
	logger, err := New(Config{
		FilePath: filepath.Join(dir, "app.log"),
		Format:   FormatJSON,
		Level:    LevelInfo,
		Rotate:   RotateDaily,
	})
	if err != nil {
		t.Fatal(err)
	}

	logger.Info("file log", zap.String("name", "gltools"))
	if err := logger.Sync(); err != nil {
		t.Fatal(err)
	}

	matches, err := filepath.Glob(filepath.Join(dir, "app-*.log"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("rotated files = %v, want one file", matches)
	}
	data, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"msg":"file log"`) {
		t.Fatalf("file log output = %s", string(data))
	}
}

func TestWithTraceStoresTraceAndSpanInContext(t *testing.T) {
	ctx := WithTrace(context.Background(), "trace-001", "span-001")

	if got := TraceID(ctx); got != "trace-001" {
		t.Fatalf("TraceID() = %q, want %q", got, "trace-001")
	}
	if got := SpanID(ctx); got != "span-001" {
		t.Fatalf("SpanID() = %q, want %q", got, "span-001")
	}
}

func TestWithContextAddsTraceFieldsToLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	logger, err := New(Config{Output: buf, Format: FormatJSON, Level: LevelInfo})
	if err != nil {
		t.Fatal(err)
	}
	ctx := WithTrace(context.Background(), "trace-002", "span-002")

	WithContext(ctx, logger).Info("trace log")

	got := buf.String()
	if !strings.Contains(got, `"trace_id":"trace-002"`) {
		t.Fatalf("log output = %s", got)
	}
	if !strings.Contains(got, `"span_id":"span-002"`) {
		t.Fatalf("log output = %s", got)
	}
}

func TestContextLevelHelpersUseDefaultLogger(t *testing.T) {
	prev := L()
	t.Cleanup(func() {
		SetDefault(prev)
	})

	buf := new(bytes.Buffer)
	logger, err := New(Config{Output: buf, Format: FormatJSON, Level: LevelDebug})
	if err != nil {
		t.Fatal(err)
	}
	SetDefault(logger)
	ctx := WithTrace(context.Background(), "trace-003", "span-003")

	DebugContext(ctx, "debug log", zap.String("name", "debug"))
	InfoContext(ctx, "info log")
	WarnContext(ctx, "warn log")
	ErrorContext(ctx, "error log", nil, zap.String("name", "error"))

	got := buf.String()
	for _, want := range []string{
		`"msg":"debug log"`,
		`"msg":"info log"`,
		`"msg":"warn log"`,
		`"msg":"error log"`,
		`"trace_id":"trace-003"`,
		`"span_id":"span-003"`,
		`"name":"debug"`,
		`"name":"error"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("log output missing %s: %s", want, got)
		}
	}
}

func TestTraceFieldsSkipsEmptyValues(t *testing.T) {
	buf := new(bytes.Buffer)
	logger, err := New(Config{Output: buf, Format: FormatJSON, Level: LevelInfo})
	if err != nil {
		t.Fatal(err)
	}
	ctx := WithTrace(context.Background(), "trace-004", "")

	WithContext(ctx, logger).Info("trace log")

	got := buf.String()
	if !strings.Contains(got, `"trace_id":"trace-004"`) {
		t.Fatalf("log output = %s", got)
	}
	if strings.Contains(got, `"span_id"`) {
		t.Fatalf("log output should not contain empty span_id: %s", got)
	}
}
