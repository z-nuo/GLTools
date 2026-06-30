package glid

import (
	"strconv"
	"testing"
)

func TestGeneratorNextReturnsIncreasingIDs(t *testing.T) {
	gen, err := NewGenerator(1)
	if err != nil {
		t.Fatal(err)
	}
	first, err := gen.Next()
	if err != nil {
		t.Fatal(err)
	}
	second, err := gen.Next()
	if err != nil {
		t.Fatal(err)
	}
	if second <= first {
		t.Fatalf("ids should increase: %d <= %d", second, first)
	}
}

func TestNewGeneratorRejectsInvalidMachineID(t *testing.T) {
	for _, machineID := range []int64{-1, 1024} {
		if _, err := NewGenerator(machineID); err == nil {
			t.Fatalf("NewGenerator(%d) error = nil, want error", machineID)
		}
	}
}

func TestGeneratorNextStringReturnsDecimalID(t *testing.T) {
	gen, err := NewGenerator(1)
	if err != nil {
		t.Fatal(err)
	}
	id, err := gen.NextString()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		t.Fatalf("NextString() = %q, want decimal int64: %v", id, err)
	}
}
