# Retry

重试，支持限定最大重试时间、最大重试次数的重试机制，同时支持指数退避算法

## 安装:

```
go get github.com/BiaoLiu/retry
```

## 使用:

限定最大重试时间

```
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

```

限定最大重试次数

```
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
```

同时限定最大重试时间与最大重试次数，以达到其中任意一个条件而终止

```
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
```


指数退避

```
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
```
