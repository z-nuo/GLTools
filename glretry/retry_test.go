package glretry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDoSucceedsAfterRetries(t *testing.T) {
	attempts := 0
	err := Do(context.Background(), FixedDelay(3, time.Millisecond), func(ctx context.Context) error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 3 {
		t.Fatalf("attempts = %d, want 3", attempts)
	}
}

func TestDoStopsWhenContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	attempts := 0
	err := Do(ctx, FixedDelay(3, time.Millisecond), func(ctx context.Context) error {
		attempts++
		return errors.New("temporary")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want %v", err, context.Canceled)
	}
	if attempts != 0 {
		t.Fatalf("attempts = %d, want 0", attempts)
	}
}
