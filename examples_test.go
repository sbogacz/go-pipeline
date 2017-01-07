package pipeline_test

import (
	"fmt"
	"strings"
	"time"

	pipeline "github.com/sbogacz/go-pipeline"

	"golang.org/x/time/rate"
)

func ExampleRateLimiter() {
	// create new rate limiter allowing 10 ops/sec
	limiter := rate.NewLimiter(10, 1)

	in := make(chan interface{}, 20)

	// create a RateLimiter operator and run it on input channel
	out := pipeline.RateLimiter(limiter).Run(in)
	startTime := time.Now()
	for i := 0; i < 100; i++ {
		in <- i
	}
	close(in)

	for _ = range out {
	} // this is just to flush the output channel
	// should have taken about 2 seconds
	fmt.Printf("After rate limiting, this took %s", startTime.Sub(time.Now()).String())
}

func ExampleFlow() {
	input0 := make(chan interface{})
	input1 := make(chan interface{})

	// multiplier takes an int and returns an Operator which multiplies the
	// input by the given int
	multiplier := func(x int) pipeline.Operator {
		return pipeline.Operator(func(in chan interface{}, out chan interface{}) {
			for m := range in {
				n := m.(int)
				out <- (int(n) * x)
			}
		})
	}

	// ifEven is an operator which filters out odds and passes evens through
	ifEven := pipeline.Operator(func(in chan interface{}, out chan interface{}) {
		for m := range in {
			n := m.(int)
			if n%2 == 0 {
				out <- n
			}
		}
	})

	// ifOdd is an operator which filters out evens and passes odds through
	ifOdd := pipeline.Operator(func(in chan interface{}, out chan interface{}) {
		for m := range in {
			n := m.(int)
			if n%2 == 1 {
				out <- n
			}
		}
	})

	// summer is an operator which aggregates input integers, and outputs the
	// total once the input channel closes
	summer := pipeline.Operator(func(in chan interface{}, out chan interface{}) {
		total := 0
		for m := range in {
			n := m.(int)
			total += n
		}
		out <- total
	})

	// for every odd, mulitply times two, and add the results
	oddFlow := pipeline.NewFlow(ifOdd, multiplier(2), summer)
	out0 := oddFlow.Run(input0)

	// for every even, multiply times three and add the results
	evenFlow := pipeline.NewFlow(ifEven, multiplier(3), summer)
	out1 := evenFlow.Run(input1)

	// use the Combine helper to merge the output of out0 and out1 into
	// a single output channel out
	out := summer.Run(pipeline.Combine(out0, out1))

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
	fmt.Printf("The total is %d\n", total) // Should total 110
}

func ExampleFlow_wordCount() {
	// Create a new intermediary type to operate on.
	// A tuple of the word and the number of times it occurred.
	type tuple struct {
		token string
		count int
	}

	// wordCount is an operator that takes in strings (words) and emits a tuple
	// of (word, 1)
	wordCount := pipeline.Operator(func(in chan interface{}, out chan interface{}) {
		for word := range in {
			out <- tuple{word.(string), 1}
		}
	})

	// countAggregator takes in tuples and aggregates their counts. Outputs
	// the word and count output as a string.
	countAggregator := pipeline.Operator(func(in chan interface{}, out chan interface{}) {
		counts := make(map[string]int)
		for t := range in {
			counts[t.(tuple).token] += t.(tuple).count
		}
		for word, count := range counts {
			out <- fmt.Sprintf("%s appears %d times", word, count)
		}
	})

	// Launch the word count Flow
	input := make(chan interface{})
	wordCountFlow := pipeline.NewFlow(wordCount, countAggregator)
	output := wordCountFlow.Run(input)

	// Feed in the input document
	document := "the quick fox jumps over the lazy brown dog fox fox"
	for _, word := range strings.Split(document, " ") {
		input <- word
	}
	// Signal that we are done submitting input
	close(input)

	// Read the output
	for result := range output {
		fmt.Println(result)
	}
}
