package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Foo struct {
	Name      string `mask:"name"`
	Email     string `mask:"email"`
	Password  string `mask:"password"`
	ID        string `mask:"id"`
	Address   string `mask:"addr"`
	Mobile    string `mask:"mobile"`
	Telephone string `mask:"tel"`
	Credit    string `mask:"credit"`
}

func TestStruct(t *testing.T) {
	org := &Foo{
		Name:      "王大大餅",
		ID:        "A123456789",
		Mobile:    "0912345678",
		Email:     "qqabcd@taiwanmobile.com",
		Password:  "idontwantuknow",
		Address:   "台北市大安區敦化南路二段206號7樓",
		Telephone: "(02)0800-1234",
		Credit:    "4938170030000003",
	}

	masked := &Foo{}

	err := Struct(org, masked)
	assert.Nil(t, err)
	assert.Equal(t, "王**餅", masked.Name)
	assert.Equal(t, "A12345****", masked.ID)
	assert.Equal(t, "0912***678", masked.Mobile)
	assert.Equal(t, "qqa****@taiwanmobile.com", masked.Email)
	assert.Equal(t, "************", masked.Password)
	assert.Equal(t, "台北市大安區******", masked.Address)
	assert.Equal(t, "(02)0800-****", masked.Telephone)
	assert.Equal(t, "493817******0003", masked.Credit)
}
