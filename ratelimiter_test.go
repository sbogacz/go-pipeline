package pipeline

import (
	"testing"

	"golang.org/x/time/rate"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRateLimiter(t *testing.T) {
	Convey("RateLimiter works as expected", t, func() {
		limiter := rate.NewLimiter(10, 1)

		in := make(chan interface{}, 100)
		for i := 0; i < 100; i++ {
			in <- i
		}
		close(in)

		out := RateLimiter(limiter).Run(in)

		j := 0
		for n := range out {
			So(n, ShouldEqual, j)
			j++
		}
		So(j, ShouldEqual, 100)
	})
	Convey("RateLimiter with other operations", t, func() {
		limiter := rate.NewLimiter(100, 1)

		in := make(chan interface{}, 100)
		for i := 0; i < 100; i++ {
			in <- i
		}
		close(in)

		c := Split(in, 2)
		f1 := NewFlow(RateLimiter(limiter), multiplier(3), RateLimiter(limiter), summer)
		f2 := NewFlow(ifOdd, summer)
		out := Combine(f1.Run(c[0]), f2.Run(c[1]))

		j := 0
		for n := range out {
			So(n, ShouldBeGreaterThan, 0)
			j++
		}
		So(j, ShouldEqual, 2)
	})
}
