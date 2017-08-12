package hemera

import (
	"math"

	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/fatih/structs"
	zheSkiplist "github.com/zhenjl/skiplist"
)

type PatternField interface{}
type PatternFields map[string]PatternField

type PatternSet struct {
	Pattern  interface{}
	Weight   int
	Fields   PatternFields
	Callback interface{}
}

type PatternSets []*PatternSet

type Bucket struct {
	PatternSets *zheSkiplist.Skiplist
}

type Router struct {
	Map     *hashmap.Map
	Buckets []*Bucket
	cmp     zheSkiplist.Comparator
}

//NewRouter creaet a new router
func NewRouter() Router {
	hm := hashmap.New()
	return Router{Map: hm, cmp: zheSkiplist.BuiltinGreaterThan}
}

// Add Insert a new pattern
func (r *Router) Add(pattern, payload interface{}) {
	ps := convertToPatternSet(pattern)
	ps.Callback = payload

	for key, val := range ps.Fields {
		// create map to save key -> value pair
		if _, ok := r.Map.Get(key); !ok {
			r.Map.Put(key, hashmap.New())
		}

		patternField, _ := r.Map.Get(key)
		patternValueMap := patternField.(*hashmap.Map)
		patternValueMapValue, ok := patternValueMap.Get(val)

		var bucket *Bucket

		if !ok {
			bucket = &Bucket{}
			bucket.PatternSets = zheSkiplist.New(r.cmp)
			r.Buckets = append(r.Buckets, bucket)
			patternValueMap.Put(val, bucket)
		} else {
			bucket = patternValueMapValue.(*Bucket)
		}

		// index by weight
		bucket.PatternSets.Insert(ps.Weight, ps)
	}

}

// List return all added patterns
func (r *Router) List() PatternSets {
	visited := hashset.New()
	list := PatternSets{}

	for _, bucket := range r.Buckets {

		var rIter *zheSkiplist.Iterator
		var err error

		rIter, err = bucket.PatternSets.SelectRange(math.MaxInt16, 0)

		if err != nil {
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

// FieldsArrayEquals check if b is subset of a
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

		// when key was not indexed
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

			rIter, err = b.PatternSets.SelectRange(ps.Weight, 0)

			// no item found in range
			if err != nil {
				continue
			}

			for rIter.Next() {
				a, ok := rIter.Value().(*PatternSet)

				if !ok {
					continue
				}

				var matched bool

				// only subset match
				if a.Weight >= ps.Weight {
					matched = equals(a, ps)
				} else {
					matched = equals(ps, a)
				}

				if matched {
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
	ps.Fields = make(PatternFields)
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
