package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPasswordSha256(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	str := "12345678"
	assert.Equal("ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f", HashPasswordSha256(str))
}
