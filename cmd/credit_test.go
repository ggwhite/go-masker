package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCredit(t *testing.T) {
	assert.Equal(t, "123456******3456", Credit("1234567890123456"))
	assert.Equal(t, "123456******345", Credit("123456789012345"))
}
