package router

import (
	"math"
	"reflect"
	"sort"
	"strings"

	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/fatih/structs"
)

type PatternFieldValue interface{}
type PatternFields map[string]PatternFieldValue

type PatternSet struct {
	Pattern interface{}
	Weight  int
	Fields  PatternFields
	Payload interface{}
}

type PatternSets []*PatternSet

type Bucket struct {
	PatternSets PatternSets
	Weight      int
}

type Router struct {
	Map         *hashmap.Map
	Buckets     []*Bucket
	IsDeep      bool
	insertCount int
}

//NewRouter creaet a new router
func NewRouter(IsDeep bool) *Router {
	hm := hashmap.New()

	return &Router{Map: hm, IsDeep: IsDeep}
}

// Add Insert a new pattern
func (r *Router) Add(pattern, payload interface{}) {
	ps := r.convertToPatternSet(pattern)
	ps.Payload = payload

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

			// fmt.Printf("Create bucket Key: '%+v' Value: '%+v'\n", key, val)
			bucket = &Bucket{}

			if r.IsDeep {
				bucket.Weight = 0
			} else {
				bucket.Weight = math.MaxInt32
			}

			patternValueMap.Put(val, bucket)

			r.Buckets = append(r.Buckets, bucket)
		} else {
			bucket = patternValueMapValue.(*Bucket)
		}

		if r.IsDeep {
			if bucket.Weight < ps.Weight {
				bucket.Weight = ps.Weight
			}
		} else {
			if bucket.Weight > ps.Weight {
				bucket.Weight = ps.Weight
			}
		}

		// pattern to bucket
		bucket.PatternSets = append(bucket.PatternSets, ps)

		//sort buckets of pattern
		if r.IsDeep {
			sort.Slice(bucket.PatternSets, func(i int, j int) bool {
				return bucket.PatternSets[i].Weight > bucket.PatternSets[j].Weight
			})
		} else {
			sort.Slice(bucket.PatternSets, func(i int, j int) bool {
				return bucket.PatternSets[i].Weight < bucket.PatternSets[j].Weight
			})
		}
	}

}

func (r *Router) List() PatternSets {
	list := PatternSets{}
	visited := hashset.New()

	for _, key := range r.Map.Keys() {

		val, _ := r.Map.Get(key)
		pv := val.(*hashmap.Map)

		for _, o := range pv.Keys() {

			b, ok := pv.Get(o)

			if ok {
				ps := b.(*Bucket)

				for _, p := range ps.PatternSets {

					if !visited.Contains(p) {
						visited.Add(p)
						list = append(list, p)
					}

				}
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

	ps := r.convertToPatternSet(p)

	buckets := []*Bucket{}

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
				panic("Value is not from type *Bucket")
			}

			buckets = append(buckets, b)

			//sort buckets
			if r.IsDeep {
				sort.Slice(buckets, func(i int, j int) bool {
					return buckets[i].Weight > buckets[j].Weight
				})
			} else {
				sort.Slice(buckets, func(i int, j int) bool {
					return buckets[i].Weight < buckets[j].Weight
				})
			}

		}
	}

	/* for _, x := range buckets {
		fmt.Printf("Bucket %+v\n", x)
		for _, p := range x.PatternSets {
			fmt.Printf("		Set: %+v\n", p)
		}

	} */

	var matched bool

	for _, bucket := range buckets {
		for _, pattern := range bucket.PatternSets {

			matched = equals(ps, pattern)

			if matched {
				return pattern
			}

		}
	}

	return nil

}

// convertToPatternSet convert a struct to a patternset
func (r *Router) convertToPatternSet(p interface{}) *PatternSet {
	fields := structs.Fields(p)

	ps := &PatternSet{}
	ps.Fields = make(PatternFields)
	ps.Pattern = p
	ps.Weight = 0

	for _, field := range fields {
		if !strings.HasSuffix(field.Name(), "_") && !field.IsZero() {
			fieldKind := field.Kind()

			switch fieldKind {
			case reflect.Int8:
				fallthrough
			case reflect.Int16:
				fallthrough
			case reflect.Int32:
				fallthrough
			case reflect.Int64:
				fallthrough
			case reflect.Int:
				fallthrough
			case reflect.String:
				fallthrough
			case reflect.Bool:
				fallthrough
			case reflect.Float32:
				fallthrough
			case reflect.Float64:
				ps.Fields[field.Name()] = field.Value()
				ps.Weight++
			}
		}
	}

	// sort by insertion order
	if !r.IsDeep {
		r.insertCount++
		ps.Weight = r.insertCount
	}

	return ps
}
