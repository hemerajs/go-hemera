package hemera

import (
	"log"

	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/fatih/structs"
	zheSkiplist "github.com/zhenjl/skiplist"
)

const (
	DepthStrategy  = "depth"
	InsertStrategy = "insert"
)

type PatternField interface{}
type PatternFields map[string]PatternField

type PatternSet struct {
	Pattern  interface{}
	Weight   int64
	Fields   PatternFields
	Callback interface{}
}

type PatternSets []*PatternSet

type Bucket struct {
	PatternSets *zheSkiplist.Skiplist
}

type Router struct {
	Map      *hashmap.Map
	Buckets  []*Bucket
	Strategy string
	counter  int64
}

//NewRouter creaet a new router with indexing strategy
func NewRouter(strategy string) Router {
	hm := hashmap.New()
	return Router{Map: hm, Strategy: strategy}
}

// Add Insert a new pattern
func (r *Router) Add(pattern, payload interface{}) {
	ps := convertToPatternSet(pattern)

	ps.Callback = payload

	for key, val := range ps.Fields {
		if _, ok := r.Map.Get(key); !ok {
			r.Map.Put(key, hashmap.New())
		}

		patternField, _ := r.Map.Get(key)
		patternValueMap := patternField.(*hashmap.Map)
		patternValueMapValue, ok := patternValueMap.Get(val)

		var bucket *Bucket

		if !ok {

			var cmp zheSkiplist.Comparator
			if r.Strategy == DepthStrategy {
				cmp = zheSkiplist.BuiltinGreaterThan
			} else {
				cmp = zheSkiplist.BuiltinLessThan
			}

			bucket = &Bucket{}
			bucket.PatternSets = zheSkiplist.New(cmp)
			r.Buckets = append(r.Buckets, bucket)
			patternValueMap.Put(val, bucket)
		} else {
			bucket = patternValueMapValue.(*Bucket)
		}

		if r.Strategy == DepthStrategy {
			bucket.PatternSets.Insert(ps.Weight, ps)
		} else {
			bucket.PatternSets.Insert(r.counter, ps)
			r.counter++
		}

	}

}

// List return all added patterns
func (r *Router) List() PatternSets {
	visited := hashset.New()
	list := PatternSets{}

	for _, bucket := range r.Buckets {

		var rIter *zheSkiplist.Iterator
		var err error

		if r.Strategy == DepthStrategy {
			rIter, err = bucket.PatternSets.SelectRange(int64(bucket.PatternSets.Count()), int64(0))
		} else {
			rIter, err = bucket.PatternSets.SelectRange(int64(0), int64(r.counter))
		}

		if err != nil {
			log.Fatalf("List: No item found in selected range Counter: %v", r.counter)
			continue
		}

		for rIter.Next() {
			p, ok := rIter.Value().(*PatternSet)

			if ok && !visited.Contains(p) {
				visited.Add(p)
				list = append(list, p)
			}
		}
	}

	return list
}

// FieldsArrayEquals check if a is a subset of b
func FieldsArrayEquals(a PatternFields, b PatternFields) bool {
	for key, field := range b {
		if a[key] != field {
			return false
		}
	}

	return true
}

// equals is shorthand for FieldsArrayEquals
func equals(a *PatternSet, b *PatternSet) bool {
	return FieldsArrayEquals(a.Fields, b.Fields)
}

// Lookup Search for a specific pattern and returns it
func (r *Router) Lookup(p interface{}) *PatternSet {

	ps := convertToPatternSet(p)

	for key, val := range ps.Fields {

		// return e.g value of "topic"
		patternField, ok := r.Map.Get(key)

		// when patternKey was not indexed
		if !ok {
			continue
		}

		// convert value to hashMap
		patternValueMap := patternField.(*hashmap.Map)

		// return bucket of e.g "topic" -> "math" -> bucket
		patternValueMapValue, ok := patternValueMap.Get(val)

		if ok {
			b, ok := patternValueMapValue.(*Bucket)

			if !ok {
				continue
			}

			var rIter *zheSkiplist.Iterator
			var err error

			if r.Strategy == DepthStrategy {
				rIter, err = b.PatternSets.SelectRange(ps.Weight, int64(0))
			} else {
				rIter, err = b.PatternSets.SelectRange(int64(0), ps.Weight)
			}

			if err != nil {
				log.Fatalf("Lookup: No item found in selected range Weight: %v - Counter: %v", ps.Weight, r.counter)
				continue
			}

			for rIter.Next() {

				a, ok := rIter.Value().(*PatternSet)

				if ok && equals(ps, a) {
					return a
				}
			}

		}
	}

	return nil

}

// convertToPatternSet convert a struct to a patternset
func convertToPatternSet(p interface{}) *PatternSet {
	fields := structs.Fields(p)

	ps := &PatternSet{}
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
