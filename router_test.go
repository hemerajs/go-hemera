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

	assert.Equal(hr.Tree.Len(), 3, "Should contain one element")

}

func TestMatchedLookup(t *testing.T) {
	assert := assert.New(t)

	ch := make(chan interface{})

	hr := NewRouter()
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "payment"}, "test2")
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "test3")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")

	go hr.Lookup(ch, DynPattern{Topic: "math", Cmd: "add"})

	p := <-ch
	assert.Equal(p.(PatternSet).Callback, "test3", "Should be `test3`")

	ch = make(chan interface{})

	go hr.Lookup(ch, DynPattern{Topic: "math", Cmd: "add", A: "1"})

	p = <-ch
	assert.Equal(p.(PatternSet).Callback, "test4", "Should be `test4`")

}

func TestUnMatchedLookupWhenSubset(t *testing.T) {
	assert := assert.New(t)

	ch := make(chan interface{})

	hr := NewRouter()
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "payment"}, "test2")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")

	go hr.Lookup(ch, DynPattern{Topic: "math", Cmd: "add"})

	p := <-ch
	assert.Equal(p.(error).Error(), "Pattern not found", "Should pattern not found")

}

func TestUnMatchedLookup(t *testing.T) {
	assert := assert.New(t)

	ch := make(chan interface{})

	hr := NewRouter()
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "payment"}, "test2")
	hr.Add(DynPattern{Topic: "math", Cmd: "add22"}, "test3")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")

	go hr.Lookup(ch, DynPattern{Topic: "math", Cmd: "add222"})

	select {
	case <-ch:
		panic("Incorrect Pattern matched")
	default:
		assert.Equal(true, true, "Should not match")
	}

}

func BenchmarkLookup(b *testing.B) {

	hr := NewRouter()
	hr.Add(DynPattern{Topic: "math"}, "test")
	hr.Add(DynPattern{Topic: "payment"}, "test2")
	hr.Add(DynPattern{Topic: "math", Cmd: "add"}, "test3")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1"}, "test4")
	hr.Add(DynPattern{Topic: "math", Cmd: "add", A: "1", B: "1"}, "test5")

	for n := 0; n < b.N; n++ {
		ch := make(chan interface{})
		go hr.Lookup(ch, DynPattern{Topic: "math", Cmd: "add"})
		<-ch
	}

}
