package glconv

import "testing"

func TestToInt(t *testing.T) {
	got, err := ToInt(" 42 ")
	if err != nil {
		t.Fatalf("ToInt() unexpected error: %v", err)
	}
	if got != 42 {
		t.Fatalf("ToInt() = %d, want 42", got)
	}
}

func TestToIntInvalid(t *testing.T) {
	if _, err := ToInt("bad"); err == nil {
		t.Fatal("ToInt() error = nil, want error")
	}
}

func TestToIntDefault(t *testing.T) {
	if got := ToIntDefault("bad", 7); got != 7 {
		t.Fatalf("ToIntDefault() = %d, want 7", got)
	}
	if got := ToIntDefault("9", 7); got != 9 {
		t.Fatalf("ToIntDefault() = %d, want 9", got)
	}
}

func TestToInt64(t *testing.T) {
	got, err := ToInt64("922337203685477580")
	if err != nil {
		t.Fatalf("ToInt64() unexpected error: %v", err)
	}
	if got != 922337203685477580 {
		t.Fatalf("ToInt64() = %d, want 922337203685477580", got)
	}
}

func TestToFloat64(t *testing.T) {
	got, err := ToFloat64("3.14")
	if err != nil {
		t.Fatalf("ToFloat64() unexpected error: %v", err)
	}
	if got != 3.14 {
		t.Fatalf("ToFloat64() = %v, want 3.14", got)
	}
}

func TestToBool(t *testing.T) {
	got, err := ToBool(" true ")
	if err != nil {
		t.Fatalf("ToBool() unexpected error: %v", err)
	}
	if !got {
		t.Fatal("ToBool() = false, want true")
	}
}

func TestToBoolInvalid(t *testing.T) {
	if _, err := ToBool("maybe"); err == nil {
		t.Fatal("ToBool() error = nil, want error")
	}
}

func TestString(t *testing.T) {
	if got := String(123); got != "123" {
		t.Fatalf("String() = %q, want %q", got, "123")
	}
	if got := String(true); got != "true" {
		t.Fatalf("String() = %q, want %q", got, "true")
	}
}
