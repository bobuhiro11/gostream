# gostream

[![Build Status](https://travis-ci.org/nmi/gostream.svg?branch=master)](https://travis-ci.org/nmi/gostream)
[![Coverage Status](https://coveralls.io/repos/github/nmi/gostream/badge.svg?branch=master)](https://coveralls.io/github/nmi/gostream?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/nmi/gostream)](https://goreportcard.com/report/github.com/nmi/gostream)
[![GoDoc](https://godoc.org/github.com/nmi/gostream?status.svg)](https://godoc.org/github.com/nmi/gostream)

## Installation

```bash
$ go install github.com/nmi/gostream
```

## Usage

See `stream_test.go` for other functions.

```go
// repeat: 1, 2, 3, 1, 2, 3
repeat := gostream.GenerateRepeatStream(ctx, 1, 2, 3)

// samenum: 1000, 1000, 1000, ...
samenum := gostream.GenerateRepeatStream(ctx, 1000)

// random: 3, 3843029809, 11, ... (for example)
random := gostream.GenerateRandIntsStream(ctx)
```

## Benchmark

For some workloads, multiple streams run in parallel and scale with the number of CPUs.

```bash
$ for i in `seq 1 4`; do go test -bench=. -cpu=$i | grep "ns/op"; done
BenchmarkAppendStringsStream             1000000              1108 ns/op  # 1 CPU
BenchmarkAppendStringsStream-2           1000000              1439 ns/op  # 2 CPU
BenchmarkAppendStringsStream-3           1000000              1774 ns/op  # 3 CPU
BenchmarkAppendStringsStream-4           1000000              1992 ns/op  # 4 CPU
```
