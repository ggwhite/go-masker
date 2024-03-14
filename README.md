# Go Masker v2

[![build workflow](https://github.com/ggwhite/go-masker/actions/workflows/go.yml/badge.svg)](https://github.com/ggwhite/go-masker/actions)
[![GoDoc](https://godoc.org/github.com/ggwhite/go-masker/v2?status.svg)](https://godoc.org/github.com/ggwhite/go-masker/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/ggwhite/go-masker/v2)](https://goreportcard.com/report/github.com/ggwhite/go-masker/v2)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/ggwhite/go-masker/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/ggwhite/go-masker.svg?style=flat-square)](https://github.com/ggwhite/go-masker/releases/latest)

Go Masker v2 is a tool for masking sensitive data in Go code. It provides a simple and convenient way to replace sensitive information, such as passwords or API keys, with placeholder values.

* [Install](#install)
* [Usage](#usage)
* [Custom Masker](#custom-masker)
* [Contributors](#contributors)

## Install

To install Go Masker v2, you can use the following command:

```bash
go get -u github.com/ggwhite/go-masker
```

## Usage

To use Go Masker v2, you can create a new instance of the `Masker` type and then use its methods to mask sensitive data. For example:

```go
package main

import (
    "log"
    masker "github.com/ggwhite/go-masker/v2"
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

    m := masker.NewMaskerMarshaler()

    t, err := m.Struct(foo)
    log.Println(t)
    log.Println(err)
}
```

This will produce the following output:

```
t = &{g**hite 0987***987}
err = <nil>
```

For more information about how to use Go Masker v2, please refer to the [documentation](https://pkg.go.dev/github.com/ggwhite/go-masker/v2).

## Custom Masker

You can also create a custom masker by implementing the `Masker` interface. For example:

```go
package main

import (
    "log"
    masker "github.com/ggwhite/go-masker/v2"
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

    // Register custom masker and override default masker
	m.Register(masker.MaskerTypeEmail, &MyEmailMasker{})

	log.Println(m.Marshal(masker.MaskerTypeEmail, "email")) // myemailmasker <nil>

	// Register custom masker and use it
	m.Register("mymasker", &MyMasker{})

	log.Println(m.Marshal("mymasker", "1234567")) // mymasker <nil>
}
```


## Contributors

Thanks to all the people who already contributed!

<a href="https://github.com/ggwhite/go-masker/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=ggwhite/go-masker" />
</a>