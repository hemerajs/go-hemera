package hemera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateRouter(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter()

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

func TestAddPattern(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter()
	hr.Add(DynPattern{Topic: "math", Cmd: "add"})
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"})
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"})

	assert.Equal(hr.Len(), 3, "Should contain 3 elements")

}

func TestMatchedLookup(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter()
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "payment"}, "test2")
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "test3")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})

	assert.Equal(p.Callback, "test3", "Should be `test3`")

	p = hr.Lookup(DynPattern{Topic: "math", Cmd: "add", A: "1"})

	assert.Equal(p.Callback, "test4", "Should be `test4`")

}

func TestMatchedLookupWhenEqualWeight(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter()
	hr.Add(DynPattern{Topic: "math"}, "test")

	p := hr.Lookup(DynPattern{Topic: "math"})

	assert.Equal(p.Callback, "test", "Should be `test`")

}

func TestMatchedLookupSubset(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter()
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test1")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})

	assert.Equal(p.Callback, "test", "Should be `test`")

}

func TestMatchedLookupLast(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter()
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test1")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add", A: "1"})

	assert.Equal(p.Callback, "test1", "Should be `test`")

}

func TestMatchedLookupWhenSubset(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter()
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "payment"}, "test2")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})

	assert.Equal(p.Callback, "test", "Should be `test`")

}

func TestUnMatchedLookupNoPartialMatchSupport(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter()
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})

	assert.Empty(p, "Pattern not found", "Should pattern not found")

}

func TestUnMatchedLookupWhenTreeEmpty(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter()

	p := hr.Lookup(DynPattern{Topic: "math", Cmd: "add222"})

	assert.Empty(p, "Pattern not found", "Should pattern not found")

}

func BenchmarkLookup(b *testing.B) {

	hr := NewRouter()
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "payment"}, "test2")
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "test3")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")

	for n := 0; n < b.N; n++ {
		hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})
	}

}

func BenchmarkLookup10000(b *testing.B) {

	hr := NewRouter()

	for n := 0; n < 10000; n++ {
		hr.Add(DynPattern{Topic: "payment"}, "test2")
		hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "test3")
		hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
		hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")
	}

	for n := 0; n < b.N; n++ {
		hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})
	}

}

func BenchmarkLookup100000(b *testing.B) {

	hr := NewRouter()

	for n := 0; n < 100000; n++ {
		hr.Add(DynPattern{Topic: "payment"}, "test2")
		hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "test3")
		hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
		hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")
	}

	for n := 0; n < b.N; n++ {
		hr.Lookup(DynPattern{Topic: "math", Cmd: "add"})
	}

}

func BenchmarkAdd(b *testing.B) {

	hr := NewRouter()

	for n := 0; n < b.N; n++ {
		hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test")
	}

}
