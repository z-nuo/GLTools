package glstrings

import "testing"

func TestIsBlank(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{name: "empty", in: "", want: true},
		{name: "spaces", in: " \t\n", want: true},
		{name: "text", in: "go", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBlank(tt.in); got != tt.want {
				t.Fatalf("IsBlank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	if got := Trim(" \tGo 工具\n"); got != "Go 工具" {
		t.Fatalf("Trim() = %q, want %q", got, "Go 工具")
	}
}

func TestTruncateUsesRunes(t *testing.T) {
	if got := Truncate("你好Go", 2); got != "你好" {
		t.Fatalf("Truncate() = %q, want %q", got, "你好")
	}
}

func TestTruncateKeepsShorterInput(t *testing.T) {
	if got := Truncate("Go", 4); got != "Go" {
		t.Fatalf("Truncate() = %q, want %q", got, "Go")
	}
}

func TestTruncateNonPositiveMax(t *testing.T) {
	if got := Truncate("Go", 0); got != "" {
		t.Fatalf("Truncate() = %q, want empty string", got)
	}
}

func TestSnakeToCamel(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "basic", in: "hello_world", want: "helloWorld"},
		{name: "collapses empty parts", in: "__hello__world__", want: "helloWorld"},
		{name: "single word", in: "go", want: "go"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnakeToCamel(tt.in); got != tt.want {
				t.Fatalf("SnakeToCamel() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCamelToSnake(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "basic", in: "helloWorld", want: "hello_world"},
		{name: "acronym boundary", in: "HTTPServer", want: "http_server"},
		{name: "already lower", in: "go", want: "go"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CamelToSnake(tt.in); got != tt.want {
				t.Fatalf("CamelToSnake() = %q, want %q", got, tt.want)
			}
		})
	}
}
