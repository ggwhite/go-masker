package main

import (
	"log"

	"github.com/ggwhite/go-masker/v2"
)

type MyEmailMasker struct{}

func (m *MyEmailMasker) Marshal(s, i string) string {
	return "myemailmasker"
}

type MyMasker struct{}

func (m *MyMasker) Marshal(s, i string) string {
	return "mymasker"
}

func main() {

	m := masker.NewMaskerMarshaler()

	log.Println(m.Marshal(masker.MaskerTypeNone, "none"))
	log.Println(m.Marshal(masker.MaskerTypePassword, "password"))
	log.Println(m.Marshal(masker.MaskerTypeName, "name"))
	log.Println(m.Marshal(masker.MaskerTypeAddress, "address"))
	log.Println(m.Marshal(masker.MaskerTypeEmail, "email"))
	log.Println(m.Marshal(masker.MaskerTypeMobile, "mobile"))
	log.Println(m.Marshal(masker.MaskerTypeTel, "tel"))
	log.Println(m.Marshal(masker.MaskerTypeID, "id"))
	log.Println(m.Marshal(masker.MaskerTypeCredit, "4111111111111111"))
	log.Println(m.Marshal(masker.MaskerTypeURL, "http://john:password@localhost:3000"))

	// Register custom masker and override default masker
	m.Register(masker.MaskerTypeEmail, &MyEmailMasker{})

	log.Println(m.Marshal(masker.MaskerTypeEmail, "email"))

	// Register custom masker and use it
	m.Register("mymasker", &MyMasker{})

	log.Println(m.Marshal("mymasker", "1234567"))
}
