package glrand

import "testing"

func TestNumericCodeReturnsDigitsWithRequestedLength(t *testing.T) {
	got, err := NumericCode(6)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 6 {
		t.Fatalf("len = %d, want 6", len(got))
	}
	for _, r := range got {
		if r < '0' || r > '9' {
			t.Fatalf("non numeric rune %q", r)
		}
	}
}

func TestStringReturnsRequestedLength(t *testing.T) {
	got, err := String(12)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 12 {
		t.Fatalf("len = %d, want 12", len(got))
	}
}

func TestStringFromCharsetUsesOnlyCharset(t *testing.T) {
	got, err := StringFromCharset(20, "ab")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 20 {
		t.Fatalf("len = %d, want 20", len(got))
	}
	for _, r := range got {
		if r != 'a' && r != 'b' {
			t.Fatalf("rune = %q, want one of charset", r)
		}
	}
}

func TestStringFromCharsetSupportsMultibyteCharset(t *testing.T) {
	got, err := StringFromCharset(10, "你好")
	if err != nil {
		t.Fatal(err)
	}
	if len([]rune(got)) != 10 {
		t.Fatalf("rune len = %d, want 10", len([]rune(got)))
	}
	for _, r := range got {
		if r != '你' && r != '好' {
			t.Fatalf("rune = %q, want one of charset", r)
		}
	}
}

func TestStringFromCharsetRejectsInvalidArguments(t *testing.T) {
	if _, err := StringFromCharset(0, "abc"); err == nil {
		t.Fatal("StringFromCharset() error = nil, want error")
	}
	if _, err := StringFromCharset(4, ""); err == nil {
		t.Fatal("StringFromCharset() error = nil, want error")
	}
}
