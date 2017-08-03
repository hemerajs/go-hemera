# Hemera - Go Client
[Hemera](https://github.com/hemerajs/hemera) client for the language Go.

**Status:** Experimental

## Install

```
go get github.com/nats-io/go-nats
go get github.com/nats-io/nuid
```

### Add

```go
// Connect to NATS
nc, _ := nats.Connect(nats.DefaultURL)
// Create hemera struct
hemera := server.Hemera{Conn: nc}
pattern := server.Pattern{"topic": "math", "cmd": "add"}

// Simple hemera add
hemera.Add(pattern, func(req server.Pattern, reply server.Reply) {
  fmt.Printf("Request: %+v\n", req)
  reply(payload | error)
})
```

### Act

```go
// Connect to NATS
nc, _ := nats.Connect(nats.DefaultURL)
// Create hemera struct
hemera := server.Hemera{Conn: nc}
pattern := server.Pattern{"topic": "math", "cmd": "add", "a": 1, "b": 2}

// Simple hemera act
hemera.Act(requestPattern, func(resp server.ClientResult) {
  fmt.Printf("Response: %+v\n", resp)
})
```

## TODO

- [ ] Implement Add and Act
- [ ] Create Context
- [ ] Handle trace, meta and delegate informations
- [ ] Implement Pattern matching (router)
