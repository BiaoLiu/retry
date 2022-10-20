package retry

import (
	"errors"
	"time"

	"github.com/rfyiamcool/backoff"
)

var (
	ErrMaxRetryCount = errors.New("重试次数已达到最大限制")
	ErrMaxRetryTime  = errors.New("重试时间已达到最大限制")
)

type Retryable struct {
	FirstRetryTime int64 `json:"firstRetryTime"`
	RetriedCount   int64 `json:"retriedCount"`
}

type RetryableFunc func(firstRetryTime int64, retriedCount int64, delay time.Duration) error

type Option func(o *Retry)

func WithFirstRetryTime(t int64) Option {
	return func(o *Retry) {
		o.firstRetryTime = t
	}
}

func WithRetriedCount(n int64) Option {
	return func(o *Retry) {
		o.retriedCount = n
	}
}

func WithMaxRetryCount(n int64) Option {
	return func(o *Retry) {
		o.maxRetryCount = n
	}
}

func WithMaxRetryTime(d time.Duration) Option {
	return func(o *Retry) {
		o.maxRetryTime = d
	}
}

func WithMinDelay(d time.Duration) Option {
	return func(o *Retry) {
		o.minDelay = d
	}
}

func WithMaxDelay(d time.Duration) Option {
	return func(o *Retry) {
		o.maxDelay = d
	}
}

func WithFactor(f float64) Option {
	return func(o *Retry) {
		o.factor = f
	}
}

func WithJitterFlag(j bool) Option {
	return func(o *Retry) {
		o.jitter = j
	}
}

// Retry .
type Retry struct {
	firstRetryTime int64         // 首次重试时间
	retriedCount   int64         // 已重试次数
	maxRetryCount  int64         // 最大重试次数
	maxRetryTime   time.Duration // 最大重试时间
	minDelay       time.Duration // 最小重试间隔
	maxDelay       time.Duration // 最大重试间隔
	factor         float64       // 指数
	jitter         bool          // 防抖动，开启将增加随机时间
}

// NewRetry new a Retry.
func NewRetry(opts ...Option) *Retry {
	retry := Retry{
		maxRetryCount: 0,
		maxRetryTime:  0,
		minDelay:      3 * time.Second,
		factor:        1,
	}
	for _, opt := range opts {
		opt(&retry)
	}
	if retry.firstRetryTime <= 0 {
		retry.firstRetryTime = time.Now().Unix()
	}
	if retry.retriedCount < 0 {
		retry.retriedCount = 0
	}
	return &retry
}

func (r *Retry) FirstRetryTime() int64 {
	return r.firstRetryTime
}

func (r *Retry) RetriedCount() int64 {
	return r.retriedCount
}

// Do retry
func (r *Retry) Do(retryableFunc RetryableFunc) error {
	b := backoff.NewBackOff(
		backoff.WithMinDelay(r.minDelay),
		backoff.WithMaxDelay(r.maxDelay),
		backoff.WithFactor(r.factor),
		backoff.WithJitterFlag(r.jitter),
	)
	var i int64 = 0
	for i = 0; i < r.retriedCount; i++ {
		b.Duration()
	}
	delay := b.Duration()

	if r.maxRetryTime > 0 {
		if time.Now().Add(-r.maxRetryTime).Unix() > r.firstRetryTime {
			return ErrMaxRetryTime
		}
	}
	if r.maxRetryCount > 0 {
		if r.retriedCount >= r.maxRetryCount {
			return ErrMaxRetryCount
		}
	}
	r.retriedCount += 1
	err := retryableFunc(r.firstRetryTime, r.retriedCount, delay)
	return err
}
