package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	assert.Equal(t, "A12345****", ID("A123456789"))
	assert.Equal(t, "A12****", ID("A12"))
	assert.Equal(t, "A****", ID("A"))
}
