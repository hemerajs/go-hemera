package hemera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateRouter(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(DepthStrategy)

	assert.NotEqual(hr, nil, "they should not nil")

}

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

type TestIntPattern struct {
	Topic string
	Cmd   string
	A     int
	B     int
}

var hrouterDepth = NewRouter(DepthStrategy)
var hrouterInsert = NewRouter(InsertStrategy)

/**
* Pattern weight order
 */

func TestAddPatternDepth(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(DepthStrategy)
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "dede")
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "deded")

	assert.Equal(len(hr.List()), 2, "Should contain 2 elements")

}

func TestMatchedLookupDepth(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(DepthStrategy)
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "payment"}, "test2")
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "test3")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})

	assert.Equal(p.Callback, "test3", "Should be `test3`")

}

func TestMatchedLookupWhenEqualWeightDepth(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(DepthStrategy)
	hr.Add(DynPattern{Topic: "math"}, "test")

	p := hr.Lookup(DynPattern{Topic: "math"})

	assert.Equal(p.Callback, "test", "Should be `test`")

}

func TestMatchedLookupSubsetDepth(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(DepthStrategy)
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test1")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})

	assert.Equal(p.Callback, "test", "Should be `test`")

}

func TestMatchedLookupNotExistKeyDepth(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(DepthStrategy)
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "test1")

	p := hr.Lookup(TestIntPattern{Topic: "math", Cmd: "add", A: 1, B: 1})

	assert.Equal(p.Callback, "test1", "Should be `test1`")

}

func TestMatchedLookupLastDepth(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(DepthStrategy)
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test1")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add", A: "1"})

	assert.Equal(p.Callback, "test1", "Should be `test`")

}

func TestMatchedLookupWhenSubsetDepth(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(DepthStrategy)
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "payment"}, "test2")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})

	assert.Equal(p.Callback, "test", "Should be `test`")

}

func TestUnMatchedLookupNoPartialMatchSupportDepth(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(DepthStrategy)
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})

	assert.Empty(p, "Pattern not found", "Should pattern not found")

}

func TestUnMatchedLookupWhenTreeEmptyDepth(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(DepthStrategy)

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add222"})

	assert.Empty(p, "Pattern not found", "Should pattern not found")

}

func BenchmarkLookupWeightDepth7(b *testing.B) {

	for n := 0; n < b.N; n++ {
		hrouterDepth.Lookup(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo", D: "11", E: "d23"})
	}

}

func BenchmarkLookupWeightDepth6(b *testing.B) {

	for n := 0; n < b.N; n++ {
		hrouterDepth.Lookup(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo", D: "dedede"})
	}

}

func BenchmarkLookupWeightDepth5(b *testing.B) {

	for n := 0; n < b.N; n++ {
		hrouterDepth.Lookup(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo"})
	}

}

func BenchmarkLookupWeightDepth3(b *testing.B) {

	for n := 0; n < b.N; n++ {
		hrouterDepth.Lookup(DynPattern{Topic: "math", Cmd: "add", A: "1"})
	}

}

func BenchmarkLookupWeightDepth2(b *testing.B) {

	for n := 0; n < b.N; n++ {
		hrouterDepth.Lookup(DynPattern{Topic: "math", Cmd: "add"})
	}

}

func BenchmarkLookupWeightDepth1(b *testing.B) {

	for n := 0; n < b.N; n++ {
		hrouterDepth.Lookup(DynPattern{Topic: "math"})
	}

}

func BenchmarkListDepth(b *testing.B) {
	hr := NewRouter(DepthStrategy)

	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo"}, "test4")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo"}, "test4")

	for n := 0; n < b.N; n++ {
		hr.List()
	}

}

func BenchmarkAddDepth(b *testing.B) {

	hr := NewRouter(DepthStrategy)

	for n := 0; n < b.N; n++ {
		hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test")
	}

}

/**
* Insertion order
 */

func TestAddPatternInsert(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(InsertStrategy)
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "dede")
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "deded")

	assert.Equal(len(hr.List()), 2, "Should contain 2 elements")

}

func TestMatchedLookupInsert(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(InsertStrategy)
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "payment"}, "test2")
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "test3")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})

	assert.Equal(p.Callback, "test", "Should be `test`")

}

func TestMatchedLookupWhenEqualWeightInsert(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(InsertStrategy)
	hr.Add(DynPattern{Topic: "math"}, "test")

	p := hr.Lookup(DynPattern{Topic: "math"})

	assert.Equal(p.Callback, "test", "Should be `test`")

}

func TestMatchedLookupSubsetInsert(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(InsertStrategy)
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test1")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})

	assert.Equal(p.Callback, "test", "Should be `test`")

}

func TestMatchedLookupNotExistKeyInsert(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(InsertStrategy)
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "test1")

	p := hr.Lookup(TestIntPattern{Topic: "math", Cmd: "add", A: 1, B: 1})

	assert.Equal(p.Callback, "test", "Should be `test`")

}

func TestMatchedLookupLastInsert(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(InsertStrategy)
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test1")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add", A: "1"})

	assert.Equal(p.Callback, "test", "Should be `test`")

}

func TestMatchedLookupWhenSubsetInsert(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(InsertStrategy)
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "payment"}, "test2")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})

	assert.Equal(p.Callback, "test", "Should be `test`")

}

func TestUnMatchedLookupNoPartialMatchSupportInsert(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(InsertStrategy)
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})

	assert.Empty(p, "Pattern not found", "Should pattern not found")

}

func TestUnMatchedLookupWhenTreeEmptyInsert(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(InsertStrategy)

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add222"})

	assert.Empty(p, "Pattern not found", "Should pattern not found")

}

func BenchmarkLookupWeightInsert7(b *testing.B) {

	for n := 0; n < b.N; n++ {
		hrouterInsert.Lookup(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo", D: "11", E: "d23"})
	}

}

func BenchmarkLookupWeightInsert6(b *testing.B) {

	for n := 0; n < b.N; n++ {
		hrouterInsert.Lookup(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo", D: "dedede"})
	}

}

func BenchmarkLookupWeightInsert5(b *testing.B) {

	for n := 0; n < b.N; n++ {
		hrouterInsert.Lookup(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo"})
	}

}

func BenchmarkLookupWeighInsert3(b *testing.B) {

	for n := 0; n < b.N; n++ {
		hrouterInsert.Lookup(DynPattern{Topic: "math", Cmd: "add", A: "1"})
	}

}

func BenchmarkLookupWeightInsert2(b *testing.B) {

	for n := 0; n < b.N; n++ {
		hrouterInsert.Lookup(DynPattern{Topic: "math", Cmd: "add"})
	}

}

func BenchmarkLookupWeightInsert1(b *testing.B) {

	for n := 0; n < b.N; n++ {
		hrouterInsert.Lookup(DynPattern{Topic: "math"})
	}

}

func BenchmarkListInsert(b *testing.B) {
	hr := NewRouter(InsertStrategy)

	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo"}, "test4")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo"}, "test4")

	for n := 0; n < b.N; n++ {
		hr.List()
	}

}

func BenchmarkAddInsert(b *testing.B) {

	hr := NewRouter(InsertStrategy)

	for n := 0; n < b.N; n++ {
		hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test")
	}

}

func init() {

	for n := 0; n < 10000; n++ {
		hrouterDepth.Add(DynPattern{Topic: "payment"}, "test2")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add"}, "test3")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")
		hrouterDepth.Add(DynPattern{Topic: "payment"}, "test2")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add"}, "test3")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo"}, "test4")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo"}, "test4")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo", D: "dedede"}, "test4")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo", D: "11", E: "d23"}, "test4")
		hrouterDepth.Add(DynPattern{Topic: "payment"}, "test2")
	}

	for n := 0; n < 10000; n++ {
		hrouterInsert.Add(DynPattern{Topic: "payment"}, "test2")
		hrouterInsert.Add(DynPattern{Topic: "math", Cmd: "add"}, "test3")
		hrouterInsert.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
		hrouterInsert.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")
		hrouterInsert.Add(DynPattern{Topic: "payment"}, "test2")
		hrouterInsert.Add(DynPattern{Topic: "math", Cmd: "add"}, "test3")
		hrouterInsert.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
		hrouterInsert.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo"}, "test4")
		hrouterInsert.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")
		hrouterInsert.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo"}, "test4")
		hrouterInsert.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo", D: "dedede"}, "test4")
		hrouterInsert.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo", D: "11", E: "d23"}, "test4")
		hrouterInsert.Add(DynPattern{Topic: "payment"}, "test2")
	}
}
