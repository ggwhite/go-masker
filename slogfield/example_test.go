package slogfield_test

import (
	"fmt"

	masker "github.com/ggwhite/go-masker/v3"
	"github.com/ggwhite/go-masker/v3/slogfield"
)

// 為求 Output 可決定性，example 不接真實 logger（時間戳不固定），改印 Key 與 Value。

func ExamplePhone() {
	attr := slogfield.Phone("phone", "0987654321")
	fmt.Println(attr.Key, attr.Value.String())
	// Output: phone 0987***321
}

func ExampleMasked() {
	attr := slogfield.Masked("phone", masker.TypeMobile, "0987654321")
	fmt.Println(attr.Key, attr.Value.String())
	// Output: phone 0987***321
}

func ExampleSensitiveAttr() {
	phone := masker.NewPhone("0987654321")
	attr := slogfield.SensitiveAttr("phone", phone)
	fmt.Println(attr.Key, attr.Value.String())
	// Output: phone 0987***321
}
