package main

import (
	"log"

	"github.com/ggwhite/go-masker/v3"
)

type MyEmailMasker struct{}

func (m *MyEmailMasker) Mask(value string) string {
	return "myemailmasker"
}

type MyMasker struct{}

func (m *MyMasker) Mask(value string) string {
	return "mymasker"
}

func main() {
	m := masker.NewMaskerMarshaler()

	log.Println(m.Marshal(masker.TypeNone, "none"))                               // none <nil>
	log.Println(m.Marshal(masker.TypePassword, "password"))                       // ************** <nil>
	log.Println(m.Marshal(masker.TypeName, "name"))                               // n**e <nil>
	log.Println(m.Marshal(masker.TypeAddress, "address"))                         // addres****** <nil>
	log.Println(m.Marshal(masker.TypeEmail, "email"))                             // ema**** <nil>
	log.Println(m.Marshal(masker.TypeMobile, "mobile"))                           // mobi*** <nil>
	log.Println(m.Marshal(masker.TypeTel, "tel"))                                 // tel <nil>
	log.Println(m.Marshal(masker.TypeID, "id"))                                   // id**** <nil>
	log.Println(m.Marshal(masker.TypeCredit, "4111111111111111"))                 // 411111******1111 <nil>
	log.Println(m.Marshal(masker.TypeURL, "http://john:password@localhost:3000")) // http://john:xxxxx@localhost:3000 <nil>

	// Register custom masker and override default masker
	m.Register(masker.TypeEmail, &MyEmailMasker{})

	log.Println(m.Marshal(masker.TypeEmail, "email"))

	// Register custom masker and use it
	m.Register("mymasker", &MyMasker{})

	log.Println(m.Marshal("mymasker", "1234567"))
}
