package hemera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateRouter(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(true)

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

var hrouterDepth = NewRouter(true)

/**
* Pattern weight order
 */

func TestAddPatternDepth(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(true)
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "dede")
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "deded")
	hr.Add(DynPattern{Topic: "payment", Cmd: "add"}, "deded")
	hr.Add(DynPattern{Topic: "payment", Cmd: "add", A: "2"}, "deded")

	assert.Equal(len(hr.List()), 4, "Should contain 4 pattern")

}

func TestMatchedLookupDepth(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(true)
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

	hr := NewRouter(true)
	hr.Add(DynPattern{Topic: "math"}, "test")

	p := hr.Lookup(DynPattern{Topic: "math"})

	assert.Equal(p.Callback, "test", "Should be `test`")

}

func TestDepthSupport(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(true)
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "test1")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})

	assert.Equal(p.Callback, "test1", "Should be `test1`")

}

func TestOrderSupport(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(false)
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "test1")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})

	assert.Equal(p.Callback, "test", "Should be `test`")

}

func TestMatchedLookupNotExistKeyDepth(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(true)
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "test1")

	p := hr.Lookup(TestIntPattern{Topic: "math", Cmd: "add", A: 1, B: 1})

	assert.Equal(p.Callback, "test1", "Should be `test1`")

}

func TestMatchedLookupLastDepth(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(true)
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test1")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add", A: "1"})

	assert.Equal(p.Callback, "test1", "Should be `test`")

}

func TestMatchedLookupWhenSubsetDepth(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(true)
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "payment"}, "test2")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})

	assert.Equal(p.Callback, "test", "Should be `test`")

}

func TestUnMatchedLookupNoPartialMatchSupportDepth(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(true)
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})

	assert.Empty(p, "Pattern not found", "Should pattern not found")

}

func TestUnMatchedLookupWhenTreeEmptyDepth(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter(true)

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

func BenchmarkLookupWeightDepth4(b *testing.B) {

	for n := 0; n < b.N; n++ {
		hrouterDepth.Lookup(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2"})
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

func BenchmarkListDepth100000(b *testing.B) {

	for n := 0; n < b.N; n++ {
		hrouterDepth.List()
	}

}

func BenchmarkAddDepth(b *testing.B) {

	hr := NewRouter(true)

	for n := 0; n < b.N; n++ {
		hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test")
	}

}

func init() {

	for n := 0; n < 1000; n++ {
		hrouterDepth.Add(DynPattern{Topic: "payment"}, "test1")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add"}, "test2")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test3")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test3")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo"}, "test5")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test4")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo"}, "test5")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo", D: "dedede"}, "test6")
		hrouterDepth.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "2", C: "foo", D: "11", E: "d23"}, "test7")
		hrouterDepth.Add(DynPattern{Topic: "order"}, "test1")
	}
}
