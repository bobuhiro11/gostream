// Package gostream provides utility functions for channel processing
package gostream

import (
	"context"
	"math/rand"
	"sync"
)

// GenerateRepeatStream generates all elemtnes of the list infinitely
func GenerateRepeatStream(ctx context.Context, list ...interface{}) <-chan interface{} {
	s := make(chan interface{})
	go func() {
		defer close(s)
		for {
			for i := 0; i < len(list); i++ {
				select {
				case s <- list[i]:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return s
}

// GenerateRandIntsStream generates a random integer sequence infinitely.
func GenerateRandIntsStream(ctx context.Context) <-chan interface{} {
	s := make(chan interface{})
	go func() {
		defer close(s)
		for {
			select {
			case s <- rand.Int():
			case <-ctx.Done():
				return
			}
		}
	}()

	return s
}

// GenerateSerialIntsStream produces an infinite serial integer sequence in ascending order.
func GenerateSerialIntsStream(ctx context.Context) <-chan interface{} {
	s := make(chan interface{})
	go func() {
		defer close(s)
		for i := 0; ; i++ {
			select {
			case s <- i:
			case <-ctx.Done():
				return
			}
		}
	}()

	return s
}

// Map first applies the function f to the first element of all input streams.
// Also, Map does the same for the second and subsequent elements.
func Map(ctx context.Context,
	f func(args ...interface{}) interface{},
	inStreams ...<-chan interface{}) <-chan interface{} {
	s := make(chan interface{})
	go func() {
		defer close(s)
		for {
			a := []interface{}{}
			for i := 0; i < len(inStreams); i++ {
				select {
				case x := <-inStreams[i]:
					a = append(a, x)
				case <-ctx.Done():
					return
				}
			}
			s <- f(a...)
		}
	}()
	return s
}

// Take takes the top num items from the input stream.
func Take(ctx context.Context, inStream <-chan interface{}, num int) <-chan interface{} {
	s := make(chan interface{})
	go func() {
		defer close(s)
		for i := 0; i < num; i++ {
			select {
			case x := <-inStream:
				s <- x
			case <-ctx.Done():
				return
			}
		}
	}()
	return s
}

// Merge combines multiple streams into one stream.
func Merge(ctx context.Context, inStreams ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	s := make(chan interface{})

	inputHandler := func(inStream <-chan interface{}) {
		for x := range inStream {
			s <- x
		}
		wg.Done()
	}

	for _, inStream := range inStreams {
		go inputHandler(inStream)
		wg.Add(1)
	}

	go func() {
		defer close(s)
		wg.Wait()
	}()

	return s
}
