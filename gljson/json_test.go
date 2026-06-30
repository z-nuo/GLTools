package gljson

import (
	"strings"
	"testing"
)

type jsonSample struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestMarshal(t *testing.T) {
	got, err := Marshal(jsonSample{Name: "go", Age: 22})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"name":"go","age":22}` {
		t.Fatalf("Marshal() = %s", got)
	}
}

func TestUnmarshal(t *testing.T) {
	var got jsonSample
	if err := Unmarshal([]byte(`{"name":"go","age":22}`), &got); err != nil {
		t.Fatal(err)
	}
	if got != (jsonSample{Name: "go", Age: 22}) {
		t.Fatalf("Unmarshal() = %+v", got)
	}
}

func TestValid(t *testing.T) {
	if !Valid(`{"ok":true}`) {
		t.Fatal("Valid() = false, want true")
	}
	if Valid(`{"ok":`) {
		t.Fatal("Valid() = true, want false")
	}
}

func TestPretty(t *testing.T) {
	got, err := Pretty(jsonSample{Name: "go", Age: 22})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "\n") {
		t.Fatalf("Pretty() = %q, want newlines", got)
	}
	if !strings.Contains(got, `  "name": "go"`) {
		t.Fatalf("Pretty() = %q, want indented field", got)
	}
}
