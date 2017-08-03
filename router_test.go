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
