package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMobile(t *testing.T) {
	assert.Equal(t, "0978***978", Mobile("0978978978"))
}
