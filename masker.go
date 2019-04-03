// Package masker Provide mask format of Taiwan usually used(Name, Address, Email, ID ...etc.),
package masker

import (
	"fmt"
	"math"
	"reflect"
	"strings"
)

const tagName = "mask"

type mtype string

// Maske Types of format string
const (
	MPassword   mtype = "password"
	MName             = "name"
	MAddress          = "addr"
	MEmail            = "email"
	MMobile           = "mobile"
	MTelephone        = "tel"
	MId               = "id"
	MCreditCard       = "credit"
)

// Masker is a instance to marshal masked string
type Masker struct{}

// Struct must input two pointer struct, source(s) and target(t), add tag mask on struct fields, after Struct(), target(t)'s filed will be masked with the tag format type
//
// Example:
//   type Foo struct {
//   	Name      string `mask:"name"`
//   	Email     string `mask:"email"`
//   	Password  string `mask:"password"`
//   	ID        string `mask:"id"`
//   	Address   string `mask:"addr"`
//   	Mobile    string `mask:"mobile"`
//   	Telephone string `mask:"tel"`
//   	Credit    string `mask:"credit"`
//   }
//
//   func main() {
//   	s := &Foo{
//			Name: ...,
//			Email: ...,
//			Password: ...,
//   	}
//   	t := &Foo{}
//
//   	Struct(s, t)
//   	fmt.Println(t)
//   }
func (m *Masker) Struct(s, t interface{}) error {
	sv := reflect.ValueOf(s)
	if sv.Kind() != reflect.Ptr {
		return fmt.Errorf("source non-pointer %v", sv.Type())
	}
	tv := reflect.ValueOf(t)
	if tv.Kind() != reflect.Ptr {
		return fmt.Errorf("target non-pointer %v", tv.Type())
	}
	if sv.Type() != tv.Type() {
		return fmt.Errorf("source(%v) and target(%v) are not same type", sv.Type(), tv.Type())
	}
	sv = sv.Elem()
	tv = tv.Elem()

	for i := 0; i < sv.NumField(); i++ {
		if mtag, ok := sv.Type().Field(i).Tag.Lookup(tagName); ok {
			switch mtype(mtag) {
			case MPassword:
				tv.Field(i).SetString(m.Password(sv.Field(i).String()))
			case MName:
				tv.Field(i).SetString(m.Name(sv.Field(i).String()))
			case MAddress:
				tv.Field(i).SetString(m.Address(sv.Field(i).String()))
			case MEmail:
				tv.Field(i).SetString(m.Email(sv.Field(i).String()))
			case MMobile:
				tv.Field(i).SetString(m.Mobile(sv.Field(i).String()))
			case MId:
				tv.Field(i).SetString(m.ID(sv.Field(i).String()))
			case MTelephone:
				tv.Field(i).SetString(m.Telephone(sv.Field(i).String()))
			case MCreditCard:
				tv.Field(i).SetString(m.CreditCard(sv.Field(i).String()))
			default:
				tv.Field(i).Set(sv.Field(i))
			}
		} else {
			tv.Field(i).Set(sv.Field(i))
		}
	}

	return nil
}

// Name mask the second world and the third world
//
// Example:
//   input: ABCD
//   output: A**D
func (*Masker) Name(i string) string {
	l := len(i)

	if l == 2 || l == 3 {
		return overlay(i, "**", 1, 2)
	}

	if l > 3 {
		return overlay(i, "**", 1, 3)
	}

	return "**"
}

// ID mask last 4 worlds of ID number
//
// Example:
//   input: A123456789
//   output: A12345****
func (*Masker) ID(i string) string {
	return overlay(i, "****", 6, 10)
}

// Address keep first 6 worlds, mask the overs
//
// Example:
//   input: 台北市內湖區內湖路一段737巷1號1樓
//   output: 台北市內湖區******
func (*Masker) Address(i string) string {
	l := len(i)
	if l <= 6 {
		return "******"
	}
	return overlay(i, "******", 6, math.MaxInt64)
}

// CreditCard mask middle 6 worlds from 7'th world
//
// Example:
//   input1: 1234567890123456 (VISA, JCB, MasterCard)(len = 16)
//   output1: 123456******3456
//   input2: 123456789012345` (American Express)(len = 15)
//   output2: 123456******345`
func (*Masker) CreditCard(i string) string {
	return overlay(i, "******", 6, 12)
}

// Email keep domain and first 3 worlds
//
// Example:
//   input: ggw.chang@gmail.com
//   output: ggw****@gmail.com
func (*Masker) Email(i string) string {
	tmp := strings.Split(i, "@")
	addr := tmp[0]
	domain := tmp[1]

	addr = overlay(addr, "****", 3, 7)

	return addr + "@" + domain
}

// Mobile mask mobile 3 worlds from 4'th world
//
// Example:
//   input: 0987654321
//   output: 0987***321
func (*Masker) Mobile(i string) string {
	return overlay(i, "***", 4, 7)
}

// Telephone remove `(`, `)`, ` `, `-` chart, and mask last 4 worlds of telephone number, format to `(??)????-????`
//
// Example:
//   input: 0227993078
//   output: (02)2799-****
func (*Masker) Telephone(i string) string {
	i = strings.Replace(i, " ", "", -1)
	i = strings.Replace(i, "(", "", -1)
	i = strings.Replace(i, ")", "", -1)
	i = strings.Replace(i, "-", "", -1)

	l := len(i)

	if l != 10 && l != 8 {
		return i
	}

	ans := ""

	if l == 10 {
		ans += "("
		ans += i[:2]
		ans += ")"
		i = i[2:]
	}

	ans += i[:4]
	ans += "-"
	ans += "****"

	return ans
}

// Password always return `************`
func (*Masker) Password(i string) string {
	return "************"
}

// New create Masker
func New() *Masker {
	return &Masker{}
}

var instance *Masker

func init() {
	instance = New()
}

// Struct must input two pointer struct, source(s) and target(t), add tag mask on struct fields, after Struct(), target(t)'s filed will be masked with the tag format type
//
// Example:
//   type Foo struct {
//   	Name      string `mask:"name"`
//   	Email     string `mask:"email"`
//   	Password  string `mask:"password"`
//   	ID        string `mask:"id"`
//   	Address   string `mask:"addr"`
//   	Mobile    string `mask:"mobile"`
//   	Telephone string `mask:"tel"`
//   	Credit    string `mask:"credit"`
//   }
//
//   func main() {
//   	s := &Foo{
//			Name: ...,
//			Email: ...,
//			Password: ...,
//   	}
//   	t := &Foo{}
//
//   	Struct(s, t)
//   	fmt.Println(t)
//   }
func Struct(s, t interface{}) error {
	return instance.Struct(s, t)
}

// Name mask the second world and the third world
//
// Example:
//   input: ABCD
//   output: A**D
func Name(i string) string {
	return instance.Name(i)
}

// ID mask last 4 worlds of ID number
//
// Example:
//   input: A123456789
//   output: A12345****
func ID(i string) string {
	return instance.ID(i)
}

// Address keep first 6 worlds, mask the overs
//
// Example:
//   input: 台北市內湖區內湖路一段737巷1號1樓
//   output: 台北市內湖區******
func Address(i string) string {
	return instance.Address(i)
}

// CreditCard mask middle 6 worlds from 7'th world
//
// Example:
//   input1: 1234567890123456 (VISA, JCB, MasterCard)(len = 16)
//   output1: 123456******3456
//   input2: 123456789012345 (American Express)(len = 15)
//   output2: 123456******345
func CreditCard(i string) string {
	return instance.CreditCard(i)
}

// Email keep domain and first 3 worlds
//
// Example:
//   input: ggw.chang@gmail.com
//   output: ggw****@gmail.com
func Email(i string) string {
	return instance.Email(i)
}

// Mobile mask mobile 3 worlds from 4'th world
//
// Example:
//   input: 0987654321
//   output: 0987***321
func Mobile(i string) string {
	return instance.Mobile(i)
}

// Telephone remove `(`, `)`, ` `, `-` chart, and mask last 4 worlds of telephone number, format to `(??)????-????`
//
// Example:
//   input: 0227993078
//   output: (02)2799-****
func Telephone(i string) string {
	return instance.Telephone(i)
}

// Password always return `************`
func Password(i string) string {
	return instance.Password(i)
}
