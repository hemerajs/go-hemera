<p align="center">
<img src="https://github.com/StarpTech/hemera/raw/master/media/hemera-logo.png" alt="Hemera" style="max-width:100%;">
</p>

<p align="center">
<a href="http://travis-ci.org/hemerajs/go-hemera"><img src="https://camo.githubusercontent.com/e63eeeaa28adaf0d6aa7abd5ca9d2dd1f2f7293d/68747470733a2f2f7472617669732d63692e6f72672f68656d6572616a732f676f2d68656d6572612e7376673f6272616e63683d6d6173746572" alt="Build Status" data-canonical-src="https://travis-ci.org/hemerajs/go-hemera.svg?branch=master" style="max-width:100%;"></a>
<a href="http://opensource.org/licenses/MIT"><img src="https://camo.githubusercontent.com/311762166ef25238116d3cadd22fcb6091edab98/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f4c6963656e73652d4d49542d626c75652e737667" alt="License MIT" data-canonical-src="https://img.shields.io/badge/License-MIT-blue.svg" style="max-width:100%;"></a>
</p>

<p align="center">
A <a href="https://golang.org/">Go</a> microservices toolkit for the <a href="https://nats.io">NATS messaging system</a>
</p>

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
ctx := hemera.Act(requestPattern, res)

res = &Response{}
ctx = hemera.Act(requestPattern, res, ctx)

log.Printf("Response %+v", res)
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
