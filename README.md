# Hemera - Go Client
[Hemera](https://github.com/hemerajs/hemera) client for the language Go.

[![License MIT](https://img.shields.io/badge/License-MIT-blue.svg)](http://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/hemerajs/go-hemera.svg?branch=master)](http://travis-ci.org/hemerajs/go-hemera)

**Status:** Experimental

## Info
The first step is to provide the basic RPC functionality. It should be able to define patterns and call them.
JSON is chosen as default protocol.

## Install

```
go get github.com/nats-io/go-nats
go get github.com/nats-io/nuid
go get github.com/fatih/structs
// Testing
go get github.com/stretchr/testify/assert
```

### Example

```go
// Connect to NATS
nc, _ := nats.Connect(nats.DefaultURL)
// Create hemera struct
hemera, _ := server.NewHemera(nc)
pattern := server.Pattern{"topic": "math", "cmd": "add"}

// Simple hemera add
hemera.Add(pattern, func(req server.Pattern, reply server.Reply) {
  fmt.Printf("Request: %+v\n", req)
  reply(payload | error)
})

// Pattern
request := server.Pattern{"topic": "math", "cmd": "add", "a": 1, "b": 2}

// Simple hemera act
hemera.Act(request, func(resp server.ClientResult) {
  fmt.Printf("Response: %+v\n", resp)
})
```

## TODO
- [X] Implement Add and Act
- [ ] Create Context (trace, meta, delegate)
- [ ] Implement Pattern matching (router)
