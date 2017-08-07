# Hemera - Go Client
[Hemera](https://github.com/hemerajs/hemera) client for the language Go.

[![License MIT](https://img.shields.io/badge/License-MIT-blue.svg)](http://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/hemerajs/go-hemera.svg?branch=master)](http://travis-ci.org/hemerajs/go-hemera)

**Status:** Experimental

## Info
This client is under development. The first step is to provide the basic RPC functionality. It should be able to define patterns and call them.
JSON is chosen as default protocol.

## Install

```
go get github.com/nats-io/go-nats
go get github.com/nats-io/nuid
go get github.com/fatih/structs
go get github.com/mitchellh/mapstructure
// Testing
go get github.com/stretchr/testify/assert
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
hemera, _ := server.Create(nc, server.Timeout(2000), ...)

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

## TODO
- [ ] Setup nats server for testing
- [X] Implement Add and Act
- [X] Infer Response in Act
- [X] Create Context (trace, meta, delegate)
- [ ] Clean request pattern from none primitive values
- [X] Implement basic pattern matching (router)
