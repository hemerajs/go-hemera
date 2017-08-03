# Hemera - Go Client
Experimental [Hemera](https://github.com/hemerajs/hemera) client for the language Go.

## Install

```
go get github.com/nats-io/go-nats
go get github.com/nats-io/nuid
```

## Example

1. Subscribe on `topic:math,cmd:add`
```
$ go run examples/sub.go
```

2. Act to `topic:math,cmd:add`
```
$ hemera-cli
$ connect
$ act --pattern topic:math,cmd:add,a:2,b:44
$ result: 46
```

## TODO

- [ ] Implement Add and Act
- [ ] Create Context
- [ ] Implement Router
