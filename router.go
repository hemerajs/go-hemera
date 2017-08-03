package hemera

import (
	"github.com/fatih/structs"
)

type PatternSet struct {
	fields  []string
	payload interface{}
}

type Router struct {
	items []PatternSet
}

func NewRouter() Router {
	items := make([]PatternSet, 100)
	return Router{items: items}
}

func Add(r *Router, p interface{}) error {

	fields := structs.Fields(p)
	ps := PatternSet{}

	for _, field := range fields {
		ps.fields = append(ps.fields, field.Name())
	}

	r.items = append(r.items, ps)
	return nil
}
