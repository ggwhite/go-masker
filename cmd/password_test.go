package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	assert.Equal(t, "************", Password("abc"))
	assert.Equal(t, "************", Password("台北市"))
	assert.Equal(t, "************", Password("09870987"))
}
