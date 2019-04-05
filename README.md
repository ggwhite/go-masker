# Golang Masker

Utility for mask personal information.

You can mask the string by the function of the package instance, or new a masker to mask your string.

#### Use package
``` golang
package main

import (
	masker "github.com/ggwhite/go-masker"
)

func main() {
	masker.Name("ggwhite")
	masker.ID("A123456789")
	masker.Mobile("0978978978")
}
```

#### Use new masker
``` golang
package main

import (
	masker "github.com/ggwhite/go-masker"
)

func main() {
	m := masker.New()
	m.Name("ggwhite")
	m.ID("A123456789")
	m.Mobile("0978978978")
}
```

## Mask Types

|Type        |Const        |Tag        |Description                                                                                            |
|:----------:|:-----------:|:---------:|:------------------------------------------------------------------------------------------------------|
|Name        |MName        |name       |mask the second world and the third world                                                              |
|Password    |MPassword    |password   |always return `************`                                                                           |
|Address     |MAddress     |addr       |keep first 6 worlds, mask the overs                                                                    |
|Email       |MEmail       |email      |keep domain and first 3 worlds                                                                         |
|Mobile      |MMobile      |mobile     |mask mobile 3 worlds from 4'th world                                                                   |
|Telephone   |MTelephone   |tel        |remove `(`, `)`, ` `, `-` chart, and mask last 4 worlds of telephone number, format to `(??)????-????` |
|ID          |MID          |id         |mask last 4 worlds of ID number                                                                        |
|CreditCard  |MCreditCard  |credit     |mask middle 6 worlds from 7'th world                                                                   |
|Struct      |MStruct      |struct     |mask the struct                                                                                        |

## Mask the `String`

Give mask type and input string

``` golang
package main

import (
	masker "github.com/ggwhite/go-masker"
)

func main() {
	masker.String(masker.MName, "ggwhite")
	masker.String(masker.MID, "A123456789")
	masker.String(masker.MMobile, "0987987987")
}
```
Result:
```
g**hite
A12345****
0987***987
```

## Mask the `Struct`

You can define your struct and add tag `mask` to let masker know what kind of the format to mask.

> Field must be **public** in the struct.

``` golang
package main

import (
	"log"
	masker "github.com/ggwhite/go-masker"
)

type Foo struct {
	Name   string `mask:"name"`
	Mobile string `mask:"mobile"`
}

func main() {
	foo := &Foo{
		Name:   "ggwhite",
		Mobile: "0987987987",
	}
	t, err := masker.Struct(foo)
	log.Println(t)
	log.Println(t.(*Foo))
	log.Println(err)
}
```

Result:
```
t = &{g**hite 0987***987} 
err = <nil>
```

### Struct contain struct

``` golang
package main

import (
	masker "github.com/ggwhite/go-masker"
)

type Foo struct {
	Name   string `mask:"name"`
	Mobile string `mask:"mobile"`
	Qoo    *Qoo   `mask:"struct"`
}

type Qoo struct {
	Name      string `mask:"name"`
	Telephone string `mask:"tel"`
}

func main() {
	foo := &Foo{
		Name:   "ggwhite",
		Mobile: "0987987987",
		Qoo: &Qoo{
			Name:      "gino",
			Telephone: "0287658765",
		},
	}
	t, err := masker.Struct(foo)
	log.Println(t)
	log.Println(t.(*Foo).Qoo)
	log.Println(err)
}
```

Result:
```
t = &{g**hite 0987***987 0xc00000a080}
t.Qoo = &{g**o (02)8765-****}
err = <nil>
```

### Struct contain string slice
``` golang
package main

import (
	masker "github.com/ggwhite/go-masker"
)

type Foo struct {
	Name   string `mask:"name"`
	Mobile string `mask:"mobile"`
	IDs    []string   `mask:"id"`
}

func main() {
	foo := &Foo{
		Name:   "ggwhite",
		Mobile: "0987987987",
		IDs: []string{
			"A123456789",
			"A987654321",
		},
	}
	t, err := masker.Struct(foo)
	log.Println(t)
	log.Println(err)
}
```

Result:
```
t = &{g**hite 0987***987 [A12345**** A98765****]}
err = <nil>
```