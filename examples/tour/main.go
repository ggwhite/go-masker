// Command tour walks through go-masker features one at a time.
//
// Run from the repo root:
//
//	go run ./examples/tour tags|custom|nested|sensitive
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	masker "github.com/ggwhite/go-masker/v3"
)

const (
	dim   = "\033[2m"
	red   = "\033[31m"
	green = "\033[32m"
	cyan  = "\033[36m"
	reset = "\033[0m"
)

// Account uses dynamic tags: no registration needed.
type Account struct {
	APIKey string `mask:"mid-4-4"`
	Phone  string `mask:"tel-2-3-8"`
	Mobile string `mask:"mobile-3-4"`
	SSN    string `mask:"id-0-4"`
	Code   string `mask:"first-3"`
	Token  string `mask:"last-4"`
}

// SSNMasker is a custom Masker: implement Mask(string) string.
type SSNMasker struct{}

func (SSNMasker) Mask(v string) string {
	if len(v) != 9 {
		return v
	}
	return "***-**-" + v[5:]
}

type Person struct {
	Name string `mask:"name"`
	SSN  string `mask:"ssn"`
}

type Item struct {
	ID string `mask:"id"`
}

type Order struct {
	Owner  *Person         `mask:"struct"`
	Emails []string        `mask:"email"`
	Items  map[string]Item `mask:"mapstruct"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: tour tags|custom|nested|sensitive")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "tags":
		tags()
	case "custom":
		custom()
	case "nested":
		nested()
	case "sensitive":
		sensitive()
	default:
		fmt.Println("unknown section:", os.Args[1])
		os.Exit(2)
	}
}

// printDiff prints each tagged field of before next to its masked copy.
func printDiff(before, after interface{}) {
	b := reflect.ValueOf(before).Elem()
	a := reflect.ValueOf(after).Elem()
	for i := 0; i < b.NumField(); i++ {
		fmt.Printf("%s%-11s%s %s%-14s%s → %s%s%s\n",
			cyan, b.Type().Field(i).Tag.Get("mask"), reset,
			red, b.Field(i).String(), reset,
			green, a.Field(i).String(), reset)
	}
}

func tags() {
	a := &Account{
		APIKey: "sk-abc123xyz",
		Phone:  "8675588888888",
		Mobile: "09012345678",
		SSN:    "123456789",
		Code:   "ABC123",
		Token:  "secret99",
	}
	masked, err := masker.NewMaskerMarshaler().Struct(a)
	if err != nil {
		panic(err)
	}
	printDiff(a, masked)
}

func custom() {
	p := &Person{Name: "John Doe", SSN: "123456789"}

	m := masker.NewMaskerMarshaler()
	m.Register("ssn", SSNMasker{})
	fmt.Printf("%sm.Register(\"ssn\", SSNMasker{})%s\n", dim, reset)
	fmt.Println(m.Format(p))

	h := masker.NewMaskerMarshaler(masker.WithMaskChar('#'))
	h.Register("ssn", SSNMasker{})
	fmt.Printf("%sNewMaskerMarshaler(WithMaskChar('#'))%s\n", dim, reset)
	fmt.Println(h.Format(p))
}

func nested() {
	m := masker.NewMaskerMarshaler()
	m.Register("ssn", SSNMasker{})
	o := &Order{
		Owner:  &Person{Name: "John Doe", SSN: "123456789"},
		Emails: []string{"john@gmail.com", "jd@work.io"},
		Items:  map[string]Item{"a": {ID: "A123456789"}, "b": {ID: "B223456789"}},
	}
	out, err := m.Struct(o)
	if err != nil {
		panic(err)
	}
	masked := out.(*Order)
	fmt.Printf("%sOwner  struct   %s %s%v%s\n", cyan, reset, green, masked.Owner, reset)
	fmt.Printf("%sEmails []string %s %s%v%s\n", cyan, reset, green, masked.Emails, reset)
	fmt.Printf("%sItems  mapstruct%s %s%v%s\n", cyan, reset, green, masked.Items, reset)
}

func sensitive() {
	phone := masker.NewPhone("0987654321")
	b, _ := json.Marshal(phone)
	row := func(code, out string) { fmt.Printf("%s%-34s%s %s\n", cyan, code, reset, out) }
	row(`phone := masker.NewPhone("0987654321")`, "")
	row("fmt.Println(phone)", green+phone.String()+reset)
	row("json.Marshal(phone)", green+string(b)+reset)
	row("phone.Reveal()", red+phone.Reveal()+reset)
	row(`NewPhone(v, masker.WithRedact())`, green+masker.NewPhone("0987654321", masker.WithRedact()).String()+reset)
}
