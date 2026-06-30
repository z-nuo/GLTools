package glretry

import (
	"context"
	"errors"
	"time"
)

// Operation 表示可重试执行的操作。
type Operation func(context.Context) error

// Options 定义重试次数和等待策略。
type Options struct {
	attempts int
	delay    func(int) time.Duration
}

// FixedDelay 使用固定间隔创建重试配置。
func FixedDelay(attempts int, delay time.Duration) Options {
	return Options{
		attempts: attempts,
		delay: func(int) time.Duration {
			return delay
		},
	}
}

// ExponentialBackoff 使用指数退避创建重试配置。
func ExponentialBackoff(attempts int, baseDelay time.Duration, maxDelay time.Duration) Options {
	return Options{
		attempts: attempts,
		delay: func(attempt int) time.Duration {
			delay := baseDelay
			for i := 1; i < attempt; i++ {
				delay *= 2
				if maxDelay > 0 && delay >= maxDelay {
					return maxDelay
				}
			}
			if maxDelay > 0 && delay > maxDelay {
				return maxDelay
			}
			return delay
		},
	}
}

// Do 按配置重试执行 operation，成功时返回 nil。
func Do(ctx context.Context, opts Options, operation Operation) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if opts.attempts <= 0 {
		return errors.New("attempts must be greater than 0")
	}
	if operation == nil {
		return errors.New("operation must not be nil")
	}

	var err error
	for attempt := 1; attempt <= opts.attempts; attempt++ {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		err = operation(ctx)
		if err == nil {
			return nil
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		if attempt == opts.attempts {
			return err
		}
		if err := wait(ctx, opts.delayFor(attempt)); err != nil {
			return err
		}
	}

	return err
}

func (o Options) delayFor(attempt int) time.Duration {
	if o.delay == nil {
		return 0
	}
	return o.delay(attempt)
}

func wait(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
