# Hemera - Go Client
[Hemera](https://github.com/hemerajs/hemera) client for the language Go.

[![License MIT](https://img.shields.io/badge/License-MIT-blue.svg)](http://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/hemerajs/go-hemera.svg?branch=master)](http://travis-ci.org/hemerajs/go-hemera)

**Status:** Experimental

## Install

```
go get ./..
```

### Example
```go

// Define the pattern of your server method
type MathPattern struct {
	Topic string `json:"topic"`
	Cmd string `json:"cmd"`
}

// Define the pattern of your RPC
type RequestPattern struct {
	Topic string `json:"topic" mapstructure:"topic"`
	Cmd string `json:"cmd" mapstructure:"cmd"`
	A int `json:"a" mapstructure:"a"`
	B int `json:"b" mapstructure:"b"`
}

// Define the struct of your response
type Response struct {
	Result int `json:"result"`
}

// Connect to NATS
nc, _ := nats.Connect(nats.DefaultURL)

// Create hemera struct with options
hemera, _ := server.CreateHemera(nc, server.Timeout(2000), ...)

// Define your server method
pattern := MathPattern{ Topic: "math", Cmd: "add" }
hemera.Add(pattern, func(req *RequestPattern, reply server.Reply) {
  fmt.Printf("Request: %+v\n", req)
  result := Response{Result: req.A + req.B}
  reply.Send(result)
})

// Call your server method
requestPattern := RequestPattern{ Topic: "math", Cmd: "add", A: 1, B: 2 }
hemera.Act(requestPattern, func(resp *Response, err server.Error) {
  fmt.Printf("Response: %+v\n", resp)
})
```

## Pattern matching
We implement `depth order` this will match the entry with the most properties first. We can measure this depth by counting the fields of a struct.

```go
type Foo struct {
	A int
	B int
}
```
This struct has a weight of `2`. This information is indexed with a [skiplist](http://drum.lib.umd.edu/bitstream/handle/1903/544/CS-TR-2286.1.pdf?sequence=2) structure to ensure that we have an average O(log k) efficiency.

## Benchmark
Lookup on a dataset with 100000 Pattern
```
BenchmarkLookupWeightDepth7-4                200           8430554 ns/op
BenchmarkLookupWeightDepth6-4                200           7683160 ns/op
BenchmarkLookupWeightDepth5-4                200           7030294 ns/op
BenchmarkLookupWeightDepth4-4                200           9502888 ns/op
BenchmarkLookupWeightDepth3-4                300           3840225 ns/op
BenchmarkLookupWeightDepth2-4               2000            818048 ns/op
BenchmarkLookupWeightDepth1-4             300000              6489 ns/op
BenchmarkListDepth-4                      500000              3900 ns/op
BenchmarkAddDepth-4                       200000              7778 ns/op
PASS
```


## TODO
- [X] Setup nats server for testing
- [X] Implement Add and Act
- [X] Infer Response in Act
- [X] Create Context (trace, meta, delegate)
- [X] Use tree for pattern indexing
- [ ] Clean request pattern from none primitive values
- [X] Implement basic pattern matching (router)
