package pipeline

import (
	"context"

	"golang.org/x/time/rate"
)

// RateLimiter returns a new Operator that only lets items through at the rate
// of the given `rate.Limiter`. Passing the limiter in allows you to share it
// across multiple instances of this Operator.
func RateLimiter(l *rate.Limiter) Operator {
	return func(in chan interface{}, out chan interface{}) {
		for n := range in {
			_ = l.Wait(context.Background())
			out <- n
		}
	}
}
