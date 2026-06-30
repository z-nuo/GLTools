package glconfig

import (
	"os"
	"path/filepath"
	"testing"
)

type sampleConfig struct {
	Name  string `json:"name" yaml:"name"`
	Port  int    `json:"port" yaml:"port"`
	Debug bool   `json:"debug" yaml:"debug"`
}

func TestLoadJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"name":"api","port":8080,"debug":true}`), 0o644); err != nil {
		t.Fatal(err)
	}

	var got sampleConfig
	if err := LoadJSON(path, &got); err != nil {
		t.Fatal(err)
	}

	want := sampleConfig{Name: "api", Port: 8080, Debug: true}
	if got != want {
		t.Fatalf("LoadJSON() = %+v, want %+v", got, want)
	}
}

func TestLoadYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	content := []byte("name: api\nport: 8080\ndebug: true\n")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}

	var got sampleConfig
	if err := LoadYAML(path, &got); err != nil {
		t.Fatal(err)
	}

	want := sampleConfig{Name: "api", Port: 8080, Debug: true}
	if got != want {
		t.Fatalf("LoadYAML() = %+v, want %+v", got, want)
	}
}

func TestEnvReturnsValueOrDefault(t *testing.T) {
	t.Setenv("GLTOOLS_NAME", "api")

	if got := Env("GLTOOLS_NAME", "default"); got != "api" {
		t.Fatalf("Env() = %q, want %q", got, "api")
	}
	if got := Env("GLTOOLS_MISSING", "default"); got != "default" {
		t.Fatalf("Env() = %q, want %q", got, "default")
	}
}

func TestEnvIntReturnsParsedValueOrDefault(t *testing.T) {
	t.Setenv("GLTOOLS_PORT", "8080")
	t.Setenv("GLTOOLS_BAD_PORT", "bad")

	if got := EnvInt("GLTOOLS_PORT", 80); got != 8080 {
		t.Fatalf("EnvInt() = %d, want %d", got, 8080)
	}
	if got := EnvInt("GLTOOLS_BAD_PORT", 80); got != 80 {
		t.Fatalf("EnvInt() = %d, want %d", got, 80)
	}
	if got := EnvInt("GLTOOLS_MISSING_PORT", 80); got != 80 {
		t.Fatalf("EnvInt() = %d, want %d", got, 80)
	}
}

func TestEnvBoolReturnsParsedValueOrDefault(t *testing.T) {
	t.Setenv("GLTOOLS_DEBUG", "true")
	t.Setenv("GLTOOLS_BAD_DEBUG", "maybe")

	if got := EnvBool("GLTOOLS_DEBUG", false); got != true {
		t.Fatalf("EnvBool() = %t, want %t", got, true)
	}
	if got := EnvBool("GLTOOLS_BAD_DEBUG", false); got != false {
		t.Fatalf("EnvBool() = %t, want %t", got, false)
	}
	if got := EnvBool("GLTOOLS_MISSING_DEBUG", true); got != true {
		t.Fatalf("EnvBool() = %t, want %t", got, true)
	}
}
