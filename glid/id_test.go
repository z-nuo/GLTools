package glid

import (
	"strconv"
	"sync"
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

func TestGeneratorNextReturnsUniqueIDsConcurrently(t *testing.T) {
	gen, err := NewGenerator(1)
	if err != nil {
		t.Fatal(err)
	}

	const goroutines = 16
	const perGoroutine = 500
	ids := make(chan int64, goroutines*perGoroutine)
	errs := make(chan error, goroutines*perGoroutine)

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				id, err := gen.Next()
				if err != nil {
					errs <- err
					return
				}
				ids <- id
			}
		}()
	}
	wg.Wait()
	close(ids)
	close(errs)

	for err := range errs {
		t.Fatalf("Next() error = %v", err)
	}

	seen := make(map[int64]struct{}, goroutines*perGoroutine)
	for id := range ids {
		if _, ok := seen[id]; ok {
			t.Fatalf("duplicate id generated: %d", id)
		}
		seen[id] = struct{}{}
	}
	if len(seen) != goroutines*perGoroutine {
		t.Fatalf("generated %d unique IDs, want %d", len(seen), goroutines*perGoroutine)
	}
}

func TestGeneratorNextReturnsErrorWhenClockMovesBackward(t *testing.T) {
	gen, err := NewGenerator(1)
	if err != nil {
		t.Fatal(err)
	}
	gen.lastTimestamp = currentMilli() + 1_000

	if _, err := gen.Next(); err == nil {
		t.Fatal("Next() error = nil, want clock rollback error")
	}
}
