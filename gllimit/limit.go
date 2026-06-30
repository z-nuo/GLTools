package gllimit

import (
	"context"

	"golang.org/x/time/rate"
)

// Limiter 封装令牌桶限流器。
type Limiter struct {
	limiter *rate.Limiter
}

// New 创建每秒 ratePerSecond 个请求、突发容量为 burst 的限流器。
func New(ratePerSecond float64, burst int) *Limiter {
	return &Limiter{
		limiter: rate.NewLimiter(rate.Limit(ratePerSecond), burst),
	}
}

// Allow 在当前可用令牌充足时返回 true。
func (l *Limiter) Allow() bool {
	return l.limiter.Allow()
}

// Wait 阻塞等待直到获得令牌或 ctx 取消。
func (l *Limiter) Wait(ctx context.Context) error {
	return l.limiter.Wait(ctx)
}
