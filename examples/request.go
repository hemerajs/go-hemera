// +build ignore
package main

import (
	"fmt"
	"log"
	"runtime"

	server "github.com/hemerajs/go-hemera"
	nats "github.com/nats-io/go-nats"
)

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

func main() {
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	hemera, _ := server.CreateHemera(nc)

	pattern := MathPattern{Topic: "math", Cmd: "add"}

	hemera.Add(pattern, func(req *RequestPattern, reply server.Reply, context *server.Context) {
		fmt.Printf("Request: %+v\n", req)
		result := Response{Result: req.A + req.B}
		reply.Send(result)
	})

	requestPattern := RequestPattern{
		Topic:    "math",
		Cmd:      "add",
		A:        1,
		B:        2,
		Meta:     server.Meta{"Test": 1},
		Delegate: server.Delegate{"Test": 2},
	}

	res := &Response{}
	ctx := hemera.Act(requestPattern, res)

	log.Printf("Response context: %+v\n", ctx)
	log.Printf("Response payload: %+v\n", res)

	runtime.Goexit()
}
