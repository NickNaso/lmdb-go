package lmdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	assert := assert.New(t)
	payload := []byte("hello")
	pval := toVal(payload)
	assert.EqualValues(5, pval.Size)
	assert.EqualValues(0, toVal(nil).Size)
	assert.Equal(payload, fromVal(pval))
	assert.Equal([]byte(nil), fromVal(nil))
}
