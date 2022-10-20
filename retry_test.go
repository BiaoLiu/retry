package retry

import (
	"fmt"
	"testing"
	"time"
)

func TestMaxRetryCount(t *testing.T) {
	opts := []Option{
		WithMaxRetryCount(10),
	}
	r := NewRetry(opts...)

	for i := 0; i < 20; i++ {
		err := r.Do(func(firstRetryTime int64, retriedCount int64, delay time.Duration) error {
			fmt.Println("retriedCount...", retriedCount)
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("total retriedCount", r.RetriedCount())
}

func TestMaxRetryTime(t *testing.T) {
	opts := []Option{
		WithMaxRetryTime(1 * time.Second),
	}
	r := NewRetry(opts...)

	for i := 0; i < 50; i++ {
		err := r.Do(func(firstRetryTime int64, retriedCount int64, delay time.Duration) error {
			fmt.Println("retriedCount...", retriedCount)
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("total retriedCount", r.RetriedCount())
}

func TestMaxRetryTimeAndCount(t *testing.T) {
	opts := []Option{
		WithMaxRetryTime(1 * time.Second),
		WithMaxRetryCount(10),
	}
	r := NewRetry(opts...)

	for i := 0; i < 50; i++ {
		err := r.Do(func(firstRetryTime int64, retriedCount int64, delay time.Duration) error {
			fmt.Println("retriedCount...", retriedCount)
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("total retriedCount", r.RetriedCount())
}

func TestBackOfRetry(t *testing.T) {
	opts := []Option{
		WithFirstRetryTime(time.Now().Unix()),
		WithRetriedCount(0),
		WithMinDelay(5 * time.Second),
		WithMaxDelay(10 * time.Hour),
		//WithMaxRetryCount(3),
		WithFactor(2),
	}
	r := NewRetry(opts...)

	for i := 0; i < 10; i++ {
		err := r.Do(func(firstRetryTime int64, retriedCount int64, delay time.Duration) error {
			fmt.Println("delay: ", delay.Seconds())
			return nil
		})
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
