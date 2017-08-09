package hemera

import (
	"github.com/fatih/structs"
)

type PatternField interface{}
type PatternFields map[string]PatternField

type PatternSet struct {
	Pattern  interface{}
	Weight   int
	Fields   PatternFields
	Callback interface{}
}

type Router struct {
	Map map[int][]PatternSet
}

func NewRouter() Router {
	m := make(map[int][]PatternSet, 10)
	return Router{Map: m}
}

func (r *Router) Len() int {
	total := 0
	for _, bucket := range r.Map {
		for range bucket {
			total++
		}
	}
	return total
}

func (r *Router) Add(args ...interface{}) {
	if len(args) == 0 {
		panic("hemera: Requires at least one argument")
	}

	ps := convertToPatternSet(args[0])

	if len(args) == 2 {
		ps.Callback = args[1]
	}

	r.Map[ps.Weight] = append(r.Map[ps.Weight], ps)
}

func FieldsArrayEquals(a PatternFields, b PatternFields) bool {
	for key, field := range b {
		if a[key] != field {
			return false
		}
	}

	return true
}

func equals(a PatternSet, b PatternSet) bool {
	return FieldsArrayEquals(a.Fields, b.Fields)
}

func (r *Router) Lookup(p interface{}) *PatternSet {

	if len(r.Map) == 0 {
		return nil
	}

	a := convertToPatternSet(p)

	for i := a.Weight; i > 0; i-- {
		if len(r.Map[i]) != 0 {
			bucket := r.Map[i]
			for _, pset := range bucket {
				if equals(a, pset) {
					return &pset
				}
			}
		}
	}

	return nil

}

func convertToPatternSet(p interface{}) PatternSet {
	fields := structs.Fields(p)

	ps := PatternSet{}
	ps.Fields = make(PatternFields, 20)
	ps.Pattern = p
	ps.Weight = 0

	// @TODO: only primitive values
	for _, field := range fields {
		if !field.IsZero() {
			ps.Fields[field.Name()] = field.Value()
			ps.Weight++
		}
	}

	return ps
}
