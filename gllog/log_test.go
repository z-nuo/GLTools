package gllog

import (
	"bytes"
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
