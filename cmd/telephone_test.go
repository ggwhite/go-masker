package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTelephone(t *testing.T) {
	assert.Equal(t, "(02)2799-****", Telephone("(02-)27   99-3--078"))
	assert.Equal(t, "(02)2799-****", Telephone("0227993078"))
	assert.Equal(t, "(07)8807-****", Telephone("0788079966"))
}
