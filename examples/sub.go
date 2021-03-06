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

type RequestPattern struct {
	Topic string `json:"topic" mapstructure:"topic"`
	Cmd string `json:"cmd" mapstructure:"cmd"`
	A int `json:"a" mapstructure:"a"`
	B int `json:"b" mapstructure:"b"`
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	hemera, _ := server.CreateHemera(nc)

	pattern := MathPattern{ Topic: "math", Cmd: "add" }

	hemera.Add(pattern, func(req *RequestPattern, reply server.Reply, context server.Context) {
		fmt.Printf("Request: %+v\n", req)
		reply.Send(req.A + req.B)
	})

	log.Printf("Listening on \n")

	runtime.Goexit()
}
