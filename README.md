# go-pipeline

[![Build Status](https://travis-ci.org/sbogacz/go-pipeline.svg?branch=master)](https://travis-ci.org/sbogacz/go-pipeline) [![Coverage Status](https://coveralls.io/repos/github/sbogacz/go-pipeline/badge.svg?branch=master)](https://coveralls.io/github/sbogacz/go-pipeline?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/sbogacz/go-pipeline)](https://goreportcard.com/report/github.com/sbogacz/go-pipeline)

go-pipeline is a Go library to provide some channel "middleware"-like functionality.
It can lead to some clean code when processing various inputs that share a flow.

## Runner

The `Runner` interface provides much of the functionality in the package. It's
defined as
```go
type Runner interface {
	Run(chan interface{}) chan interface{}
}
```

The two implementing types provided are `Operator` and `Flow`, where the latter
is just a collection of `Operator`s, and both provide a `Run` method.

Example:

```go
func multiplier(x int) Operator {
	return Operator(func(in chan interface{}, out chan interface{}) {
		for m := range in {
			n := m.(int)
			out <- (int(n) * x)
		}
	})
}
```

## Rate Limiter

This package also provides a `RateLimiter` function which takes a rate limiter
from the "golang.org/x/time/rate" package, and returns an `Operator` which returns
a channel whose input is throttled by the provided rate limiter.

## Examples

More examples of how to use the pipeline package can be found in the test files

## Contributing

If you'd like to contribute to this project, make sure that you're running [go vet](https://golang.org/cmd/vet/)
and [go lint](https://github.com/golang/lint) before submitting a pull request. If
adding a feature or fixing a bug, please also add a test case verify the new functionality/new fix.
