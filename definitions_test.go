package pipeline

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPipeline(t *testing.T) {
	Convey("Make sure combining and applying operators works as expected", t, func() {
		input0 := make(chan interface{})
		input1 := make(chan interface{})
		// for every odd, mulitply times two, and add the results
		oddFlow := NewFlow(ifOdd, multiplier(2), summer)
		out0 := oddFlow.Run(input0)

		evenFlow := NewFlow(ifEven, multiplier(3), summer)
		out1 := evenFlow.Run(input1)

		out := summer.Run(Combine(out0, out1))

		go func() {
			for i := 0; i < 10; i++ {
				input0 <- i
				input1 <- i
			}
			close(input0)
			close(input1)
		}()
		total := <-out
		// 1 + 3 + 5 + 7 + 9 = 25... * 2 = 50
		// 2 + 4 + 6 + 8 = 20 * 3 = 60
		So(total, ShouldEqual, 110)
	})
	Convey("Make sure we didn't break math", t, func() {
		input0 := make(chan interface{})
		out := summer.Run(input0)
		input0 <- 1
		input0 <- 1
		close(input0)
		total := <-out
		So(total, ShouldEqual, 2)
	})
	Convey("Make sure splitting works", t, func() {
		input := make(chan interface{})
		go func() {
			for i := 0; i < 10; i++ {
				input <- i
			}
			close(input)
		}()

		out := Combine(Split(Combine(Split(input, 2)...), 2)...)
		total := 0
		for n := range out {
			total += n.(int)
		}
		So(total, ShouldEqual, 45*4)
	})
}

func multiplier(x int) Operator {
	return Operator(func(in chan interface{}, out chan interface{}) {
		for m := range in {
			n := m.(int)
			out <- (int(n) * x)
		}
	})
}

var ifOdd = Operator(func(in chan interface{}, out chan interface{}) {
	for m := range in {
		n := m.(int)
		if n%2 == 1 {
			out <- n
		}
	}
})

var ifEven = Operator(func(in chan interface{}, out chan interface{}) {
	for m := range in {
		n := m.(int)
		if n%2 == 0 {
			out <- n
		}
	}
})

var summer = Operator(func(in chan interface{}, out chan interface{}) {
	total := 0
	for m := range in {
		n := m.(int)
		total += n
	}
	out <- total
})
