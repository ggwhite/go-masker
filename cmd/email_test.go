package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmail(t *testing.T) {
	assert.Equal(t, "ggw****ng@gmail.com", Email("ggw.chang@gmail.com"))
	assert.Equal(t, "qq****@gmail.com", Email("qq@gmail.com"))
	assert.Equal(t, "qqa****@taiwanmobile.com", Email("qqabcd@taiwanmobile.com"))
}
