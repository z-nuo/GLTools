package glerror

import (
	"errors"
	"testing"
)

func TestWrapPreservesCodeAndUnderlyingError(t *testing.T) {
	base := errors.New("db failed")
	err := Wrap(50001, "database error", base)
	if err.Code != 50001 {
		t.Fatalf("Code = %d", err.Code)
	}
	if err.Message != "database error" {
		t.Fatalf("Message = %q", err.Message)
	}
	if !errors.Is(err, base) {
		t.Fatal("wrapped error should match base")
	}
	if !IsCode(err, 50001) {
		t.Fatal("IsCode should be true")
	}
}

func TestNewAndFrom(t *testing.T) {
	err := New(40001, "bad request")
	if err.Code != 40001 {
		t.Fatalf("Code = %d", err.Code)
	}
	if err.Message != "bad request" {
		t.Fatalf("Message = %q", err.Message)
	}
	got, ok := From(err)
	if !ok {
		t.Fatal("From() ok = false, want true")
	}
	if got != err {
		t.Fatal("From() should return original CodeError")
	}
	if IsCode(err, 40002) {
		t.Fatal("IsCode should be false for a different code")
	}
}
