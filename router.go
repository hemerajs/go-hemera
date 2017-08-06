package hemera

import (
	"sort"

	"github.com/fatih/structs"
	"github.com/google/btree"
)

type PatternField string
type PatternFields []PatternField

type PatternSet struct {
	Pattern interface{}
	Weight  int
	Fields  PatternFields
	Payload interface{}
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
		ps.Payload = args[1]
	}

	r.Tree.ReplaceOrInsert(ps)

	return nil
}

func FieldsArrayEquals(a PatternFields, b PatternFields) bool {
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func equals(a PatternSet, b PatternSet) bool {
	if b.Weight > a.Weight {
		return false
	}

	return FieldsArrayEquals(a.Fields, b.Fields)
}

func (r *Router) Lookup(result chan<- PatternSet, p interface{}) {

	ps := convertToPatternSet(p)

	r.Tree.AscendGreaterOrEqual(ps, func(i btree.Item) bool {
		ips := i.(PatternSet)
		if equals(ps, ips) {
			result <- ips
			return false
		}

		return true
	})
}

func convertToPatternSet(p interface{}) PatternSet {
	fields := structs.Fields(p)

	ps := PatternSet{}
	ps.Pattern = p
	ps.Weight = 0

	for _, field := range fields {

		if !field.IsZero() {
			pf := PatternField(field.Name())
			ps.Fields = append(ps.Fields, pf)
			ps.Weight++
		}
	}

	sort.Sort(ps.Fields)

	return ps
}

func (slice PatternFields) Len() int {
	return len(slice)
}

// Less ascending alphabetical sort order
func (slice PatternFields) Less(i, j int) bool {
	return slice[i] < slice[j]
}

func (slice PatternFields) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
