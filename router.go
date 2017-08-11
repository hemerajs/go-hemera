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
}

func NewRouter() Router {
	hm := hashmap.New()
	return Router{Map: hm}
}

func (r *Router) Add(args ...interface{}) {
	if len(args) == 0 {
		panic("hemera: Requires at least one argument")
	}

	ps := convertToPatternSet(args[0])

	if len(args) == 2 {
		ps.Callback = args[1]
	}

	for key, val := range ps.Fields {
		if _, ok := r.Map.Get(key); !ok {
			r.Map.Put(key, hashmap.New())
		}

		patternField, _ := r.Map.Get(key)
		patternValueMap := patternField.(*hashmap.Map)
		patternValueMapValue, ok := patternValueMap.Get(val)

		var bucket *Bucket

		if !ok {
			bucket = &Bucket{}
			bucket.PatternSets = zheSkiplist.New(zheSkiplist.BuiltinGreaterThan)
			r.Buckets = append(r.Buckets, bucket)
			patternValueMap.Put(val, bucket)
		} else {
			bucket = patternValueMapValue.(*Bucket)
		}

		// Add PatternSet to bucket
		bucket.PatternSets.Insert(ps.Weight, ps)

	}

}

func (r *Router) List() PatternSets {
	visited := hashset.New()
	list := PatternSets{}

	for _, bucket := range r.Buckets {
		rIter, _ := bucket.PatternSets.SelectRange(math.MaxInt64, 0)

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

func FieldsArrayEquals(a PatternFields, b PatternFields) bool {
	for key, field := range b {
		if a[key] != field {
			return false
		}
	}

	return true
}

func equals(a *PatternSet, b *PatternSet) bool {
	return FieldsArrayEquals(a.Fields, b.Fields)
}

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

			// search pattern with equal weight or less
			rIter, _ := b.PatternSets.SelectRange(ps.Weight, 0)

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
