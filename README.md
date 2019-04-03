# Golang Masker

個資隱碼工具

## Example:

##### Code:

``` golang
package main

import (
	"log"

	masker "github.com/ggwhite/go-masker"
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

func main() {
	org := &Foo{
		Name:      "王大大餅",
		ID:        "A123456789",
		Mobile:    "0912345678",
		Email:     "qqabcd@taiwanmobile.com",
		Password:  "idontwantuknow",
		Address:   "台北市大安區敦化南路五段999號200樓",
		Telephone: "(02)0800-1234",
		Credit:    "1234567890123456",
	}

	masked := &Foo{}
	err := masker.Struct(org, masked)

	log.Println(err)
	log.Println(org)
	log.Println(masked)
}
```

##### Result:

```
2019/04/03 11:13:49 <nil>
2019/04/03 11:13:49 &{王大大餅 qqabcd@taiwanmobile.com idontwantuknow A123456789 台北市大安區敦化南路五段999號200樓 0912345678 (02)0800-1234 1234567890123456}
2019/04/03 11:13:49 &{王**餅 qqa****@taiwanmobile.com ************ A12345**** 台北市大安區****** 0912***678 (02)0800-**** 123456******3456}
```