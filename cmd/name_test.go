package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {

	assert.Equal(t, "王*", Name("王蛋"))
	assert.Equal(t, "王*蛋", Name("王八蛋"))
	assert.Equal(t, "王**蛋", Name("王七八蛋"))
	assert.Equal(t, "王**九蛋", Name("王七八九蛋"))
	assert.Equal(t, "王**九十蛋", Name("王七八九十蛋"))

	assert.Equal(t, "A**n", Name("Alen"))
	assert.Equal(t, "A**n Lin", Name("Alen Lin"))
	assert.Equal(t, "J**ge Marry", Name("Jorge Marry"))

}
