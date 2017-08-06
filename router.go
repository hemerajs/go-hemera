package hemera

import (
	"github.com/fatih/structs"
	"github.com/google/btree"
)

type PatternField struct {
	Key   string
	Value interface{}
}
type PatternFields map[string]PatternField

type PatternSet struct {
	Pattern  interface{}
	Weight   int
	Fields   PatternFields
	Callback interface{}
}

func (p PatternSet) Less(item btree.Item) bool {
	return p.Weight < item.(PatternSet).Weight
}

type Router struct {
	Tree *btree.BTree
}

func NewRouter() Router {
	tree := btree.New(7)
	return Router{Tree: tree}
}

func (r *Router) Add(args ...interface{}) error {
	if len(args) == 0 {
		panic("hemera: Requires at least one argument")
	}

	ps := convertToPatternSet(args[0])

	if len(args) == 2 {
		ps.Callback = args[1]
	}

	r.Tree.ReplaceOrInsert(ps)

	return nil
}

func FieldsArrayEquals(a PatternFields, b PatternFields) bool {
	for _, field := range b {
		if a[field.Key].Value != field.Value {
			return false
		}
	}

	return true
}

func equals(a PatternSet, b PatternSet) bool {
	return FieldsArrayEquals(a.Fields, b.Fields)
}

func (r *Router) Lookup(result chan<- interface{}, p interface{}) {
	ps := convertToPatternSet(p)

	r.Tree.DescendLessOrEqual(ps, func(i btree.Item) bool {
		ips := i.(PatternSet)

		if equals(ps, ips) {
			result <- ips
			return false
		}

		return true
	})

	r.Tree.AscendGreaterOrEqual(ps, func(i btree.Item) bool {
		ips := i.(PatternSet)

		// When add pattern is not a subset
		if ips.Weight > ps.Weight {
			result <- ErrPatternNotFound
			return false
		} else if equals(ps, ips) {
			result <- ips
			return false
		}

		return true
	})
}

func convertToPatternSet(p interface{}) PatternSet {
	fields := structs.Fields(p)

	ps := PatternSet{}
	ps.Fields = make(PatternFields, 10)
	ps.Pattern = p
	ps.Weight = 0

	for _, field := range fields {
		if !field.IsZero() {
			pf := PatternField{Key: field.Name(), Value: field.Value()}
			ps.Fields[field.Name()] = pf
			ps.Weight++
		}
	}

	return ps
}
