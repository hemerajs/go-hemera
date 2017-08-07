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
	Topic string `json:"topic"`
	Cmd string `json:"cmd"`
}

type Delegate struct {
		Query string `json:"query"`
}

type Meta struct {
		Token string `json:"token"`
}

type RequestPattern struct {
	Topic string `json:"topic" mapstructure:"topic"`
	Cmd string `json:"cmd" mapstructure:"cmd"`
	A int `json:"a" mapstructure:"a"`
	B int `json:"b" mapstructure:"b"`
}

type Response struct {
	Result int `json:"result"`
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	hemera, _ := server.Create(nc)

	pattern := MathPattern{Topic: "math", Cmd: "add"}

	hemera.Add(pattern, func(req *RequestPattern, reply server.Reply, context server.Context) {
		fmt.Printf("Request: %+v\n", req)
		result := Response{Result: req.A + req.B}
		reply.Send(result)
	})

	requestPattern := RequestPattern{Topic: "math", Cmd: "add", A: 1, B: 2}
	
	hemera.Act(requestPattern, func(resp *Response, err server.Error, context server.Context) {
		fmt.Printf("Response: %+v\n", resp)
	})

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on \n")

	runtime.Goexit()
}
