package cmd

import (
	"fmt"
	"reflect"
)

const tagName = "mask"

type mtype string

// mask types
const (
	MPassword  mtype = "password"
	MName            = "name"
	MAddress         = "addr"
	MEmail           = "email"
	MMobile          = "mobile"
	MTelephone       = "tel"
	MId              = "id"
	MCredit          = "credit"
)

// Struct mask struct with tag `mask`
func Struct(source interface{}, target interface{}) error {
	sv := reflect.ValueOf(source)
	if sv.Kind() != reflect.Ptr {
		return fmt.Errorf("source non-pointer %v", sv.Type())
	}
	tv := reflect.ValueOf(target)
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
				tv.Field(i).SetString(Password(sv.Field(i).String()))
			case MName:
				tv.Field(i).SetString(Name(sv.Field(i).String()))
			case MAddress:
				tv.Field(i).SetString(Address(sv.Field(i).String()))
			case MEmail:
				tv.Field(i).SetString(Email(sv.Field(i).String()))
			case MMobile:
				tv.Field(i).SetString(Mobile(sv.Field(i).String()))
			case MId:
				tv.Field(i).SetString(ID(sv.Field(i).String()))
			case MTelephone:
				tv.Field(i).SetString(Telephone(sv.Field(i).String()))
			case MCredit:
				tv.Field(i).SetString(Credit(sv.Field(i).String()))
			default:
				tv.Field(i).Set(sv.Field(i))
			}
		} else {
			tv.Field(i).Set(sv.Field(i))
		}
	}

	return nil
}
