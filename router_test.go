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

func AddPattern(t *testing.T) {
	assert := assert.New(t)

	hr := NewRouter()
	hr.Add(Pattern{"a": 1, "b": 2})

	assert.NotEqual(hr.Len(), 1, "Should contain one element")

}
