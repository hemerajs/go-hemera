// +build ignore

package main

import (
	"log"
	server "github.com/hemerajs/go-hemera"
)

type DynPattern struct {
	Topic string
	Cmd   string
	A     string
	B     string
	C     string
	D     string
	E     string
	F     string
	G     string
	H     string
	I     string
	J     string
}


func main() {

	ch := make(chan server.PatternSet)

	hr := server.NewRouter()
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "test3")

	go hr.Lookup(ch, DynPattern{Topic: "math", Cmd: "add"})

	p := <-ch
	log.Printf("%+v", p.Pattern)
}
