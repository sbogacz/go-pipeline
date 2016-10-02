package pipeline

import "sync"

// Runner interface exposes functions that take in a chan interface{}
// and outputs to a chan interface{}
type Runner interface {
	Run(chan interface{}) chan interface{}
}

// Operator simply defines a type to simplify the following definition
type Operator func(chan interface{}, chan interface{})

// Run takes an input channel, and a series of operators, and uses the output
// of each successive operator as the input for the next
func (o Operator) Run(in chan interface{}) chan interface{} {
	out := make(chan interface{})
	go func() {
		o(in, out)
		close(out)
	}()
	return out
}

// Flow is a series of Operators that can be applied in sequence
type Flow []Operator

// NewFlow is syntactic sugar to create a Flow
func NewFlow(ops ...Operator) Flow {
	return Flow(ops)
}

// Run takes an input channel and runs the operator
func (f Flow) Run(in chan interface{}) chan interface{} {
	for _, m := range f {
		in = m.Run(in)
	}
	return in
}

// Split takes a single input channel, and broadcasts each item to each handler
// function's channel.
func Split(in chan interface{}, n int) []chan interface{} {
	outChans := make([]chan interface{}, n)
	for i := 0; i < n; i++ {
		outChans[i] = make(chan interface{})
	}
	go func() {
		for n := range in {
			for _, out := range outChans {
				out <- n
			}
		}
		for _, out := range outChans {
			close(out)
		}
	}()
	return outChans
}

// Combine takes a variable number of channels and combines their output into
// a single channel, that can still be used with an operator
func Combine(ins ...chan interface{}) chan interface{} {
	out := make(chan interface{})
	wg := &sync.WaitGroup{}
	wg.Add(len(ins))
	for _, in := range ins {
		go func(in chan interface{}) {
			for n := range in {
				out <- n
			}
			wg.Done()
		}(in)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
