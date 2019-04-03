package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	assert.Equal(t, "台北市敦化南******", Address("台北市敦化南路2段206號7樓"))
	assert.Equal(t, "******", Address("台北市"))
}
