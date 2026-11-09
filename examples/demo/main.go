// Command demo prints a struct before and after masking with go-masker.
//
// Run from the repo root:
//
//	go run ./examples/demo
package main

import (
	"fmt"
	"reflect"

	masker "github.com/ggwhite/go-masker/v3"
)

type User struct {
	Name     string `mask:"name"`
	Email    string `mask:"email"`
	Password string `mask:"password"`
	Mobile   string `mask:"mobile"`
	Tel      string `mask:"tel"`
	ID       string `mask:"id"`
	Credit   string `mask:"credit"`
	Website  string `mask:"url"`
}

const (
	red   = "\033[31m"
	green = "\033[32m"
	cyan  = "\033[36m"
	bold  = "\033[1m"
	reset = "\033[0m"
)

func main() {
	u := &User{
		Name:     "John Doe",
		Email:    "john@gmail.com",
		Password: "hunter2",
		Mobile:   "0987654321",
		Tel:      "0227993078",
		ID:       "A123456789",
		Credit:   "4111111111111111",
		Website:  "https://john:s3cret@acme.io",
	}

	m := masker.NewMaskerMarshaler()
	masked, err := m.Struct(u)
	if err != nil {
		panic(err)
	}

	before := reflect.ValueOf(u).Elem()
	after := reflect.ValueOf(masked).Elem()
	t := before.Type()

	fmt.Printf("%s%-10s %-10s %-28s %s%s\n", bold, "FIELD", "TAG", "BEFORE", "AFTER", reset)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Printf("%-10s %s%-10s%s %s%-28s%s %s%s%s\n",
			f.Name,
			cyan, f.Tag.Get("mask"), reset,
			red, before.Field(i).String(), reset,
			green, after.Field(i).String(), reset)
	}

}
