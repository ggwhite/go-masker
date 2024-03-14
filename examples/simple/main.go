package main

import (
	"log"

	"github.com/ggwhite/go-masker/v2"
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
	URL       string `mask:"url"`
	Foo       *Foo   `mask:"struct"`
}

func main() {
	m := masker.NewMaskerMarshaler()

	log.Println(m.List())

	foo1 := &Foo{
		Name:      "John Doe",
		Email:     "john@gmail.com",
		Password:  "password",
		ID:        "1234567890",
		Address:   "123 Main St",
		Mobile:    "1234567890",
		Telephone: "1234567890",
		Credit:    "4111111111111111",
		URL:       "http://john:password@localhost:3000",
		Foo: &Foo{
			Name:      "John Doe",
			Email:     "john@gmail.com",
			Password:  "password",
			ID:        "1234567890",
			Address:   "123 Main St",
			Mobile:    "1234567890",
			Telephone: "1234567890",
			Credit:    "4111111111111111",
			URL:       "http://john:password@localhost:3000",
		},
	}

	foo2, _ := m.Struct(foo1)

	log.Println(foo1)
	log.Println(foo1.Foo)
	log.Println(foo2)
	log.Println(foo2.(*Foo).Foo)
}
