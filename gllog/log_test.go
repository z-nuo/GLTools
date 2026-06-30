package gllog

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestNewCreatesJSONLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	logger, err := New(Config{Output: buf, Format: FormatJSON, Level: slog.LevelInfo})
	if err != nil {
		t.Fatal(err)
	}

	logger.Info("hello", slog.String("name", "gltools"))

	got := buf.String()
	if !strings.Contains(got, `"msg":"hello"`) {
		t.Fatalf("log output = %s", got)
	}
	if !strings.Contains(got, `"name":"gltools"`) {
		t.Fatalf("log output = %s", got)
	}
}

func TestNewCreatesTextLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	logger, err := New(Config{Output: buf, Format: FormatText, Level: slog.LevelInfo})
	if err != nil {
		t.Fatal(err)
	}

	logger.Info("hello", slog.String("name", "gltools"))

	got := buf.String()
	if !strings.Contains(got, `msg=hello`) {
		t.Fatalf("log output = %s", got)
	}
	if !strings.Contains(got, `name=gltools`) {
		t.Fatalf("log output = %s", got)
	}
}

func TestNewRejectsUnsupportedFormat(t *testing.T) {
	_, err := New(Config{Output: new(bytes.Buffer), Format: Format("xml")})
	if err == nil {
		t.Fatal("New() error = nil, want error")
	}
}

func TestSetDefaultUsesProvidedLogger(t *testing.T) {
	prev := slog.Default()
	t.Cleanup(func() {
		slog.SetDefault(prev)
	})

	buf := new(bytes.Buffer)
	logger, err := New(Config{Output: buf, Format: FormatJSON, Level: slog.LevelInfo})
	if err != nil {
		t.Fatal(err)
	}

	SetDefault(logger)
	slog.Info("default logger")

	if !strings.Contains(buf.String(), `"msg":"default logger"`) {
		t.Fatalf("default log output = %s", buf.String())
	}
}
