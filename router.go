package hemera

import (
	"sort"

	"github.com/fatih/structs"
)

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

type PatternField string
type PatternFields []PatternField

type PatternSet struct {
	fields  PatternFields
	payload interface{}
}

type Router struct {
	items []PatternSet
}

func NewRouter() Router {
	items := make([]PatternSet, 100)
	return Router{items: items}
}

func (r *Router) Add(p interface{}) error {

	fields := structs.Fields(p)
	ps := PatternSet{}

	for _, field := range fields {
		pf := PatternField(field.Name())
		ps.fields = append(ps.fields, pf)
	}

	sort.Sort(ps.fields)

	r.items = append(r.items, ps)

	return nil
}

func (r *Router) Len() int {
	return len(r.items)
}
