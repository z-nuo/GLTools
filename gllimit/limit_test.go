package gllimit

import (
	"context"
	"errors"
	"testing"
)

func TestLimiterAllowsBurstOnly(t *testing.T) {
	limiter := New(1, 1)
	if !limiter.Allow() {
		t.Fatal("first request should be allowed")
	}
	if limiter.Allow() {
		t.Fatal("second immediate request should be rejected")
	}
}

func TestLimiterWaitStopsWhenContextCanceled(t *testing.T) {
	limiter := New(1, 1)
	if !limiter.Allow() {
		t.Fatal("first request should be allowed")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := limiter.Wait(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Wait() error = %v, want %v", err, context.Canceled)
	}
}
