<p align="center">
<img src="https://github.com/StarpTech/hemera/raw/master/media/hemera-logo.png" alt="Hemera" style="max-width:100%;">
</p>

<p align="center">
A <a href="https://golang.org/">Go</a> microservices toolkit for the <a href="https://nats.io">NATS messaging system</a>
</p>

[![License MIT](https://img.shields.io/badge/License-MIT-blue.svg)](http://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/hemerajs/go-hemera.svg?branch=master)](http://travis-ci.org/hemerajs/go-hemera)

**Status:** Experimental

## Install

```
go get ./..
go get github.com/nats-io/gnatsd/server
```

### Example
```go

type MathPattern struct {
	Topic string
	Cmd   string
}

type RequestPattern struct {
	Topic    string
	Cmd      string
	A        int
	B        int
	Meta     server.Meta
	Delegate server.Delegate
}

type Response struct {
	Result int
}

nc, _ := nats.Connect(nats.DefaultURL)

hemera, _ := server.CreateHemera(nc, server.Timeout(2000), server.IndexingStrategy(DepthIndexing)...)

// Define the pattern of your action
pattern := MathPattern{Topic: "math", Cmd: "add"}
hemera.Add(pattern, func(req *RequestPattern, reply server.Reply, context *server.Context) {
	// Build response
	result := Response{Result: req.A + req.B}
	// Add meta informations
	context.Meta["key"] = "value"
	// Send it back
	reply.Send(result)
})

// Define the call of your RPC
requestPattern := RequestPattern{
	Topic: "math",
	Cmd: "add",
	A: 1,
	B: 2,
	Meta: server.Meta{ "Test": 1 },
	Delegate: server.Delegate{ "Test": 2 },
}

res := &Response{} // Pointer to struct
hemera.Act(requestPattern, res)

log.Printf("Response %v", res)
```

## Pattern matching
We implemented two indexing strategys
- `depth order` match the entry with the most properties first.
- `insertion order` match the entry with the least properties first. `(default)`

## TODO
- [X] Setup nats server for testing
- [X] Implement Add and Act
- [X] Create Context (trace, meta, delegate) structures
- [X] Use tree for pattern indexing
- [X] Support indexing by depth order
- [X] Support indexing by insetion order
- [X] Clean request pattern from none primitive values
- [X] Meta & Delegate support
- [X] Implement basic pattern matching (router)
- [ ] Implement router `remove` method

## Credits

- [Bloomrun](https://github.com/mcollina/bloomrun) the pattern matching library for NodeJs
