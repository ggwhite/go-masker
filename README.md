# Go Masker v2

[![build workflow](https://github.com/ggwhite/go-masker/actions/workflows/go.yml/badge.svg)](https://github.com/ggwhite/go-masker/actions)
[![GoDoc](https://godoc.org/github.com/ggwhite/go-masker?status.svg)](https://godoc.org/github.com/ggwhite/go-masker)
[![Go Report Card](https://goreportcard.com/badge/github.com/ggwhite/go-masker)](https://goreportcard.com/report/github.com/ggwhite/go-masker)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/ggwhite/go-masker/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/ggwhite/go-masker.svg?style=flat-square)](https://github.com/ggwhite/go-masker/releases/latest)

Go Masker v2 is a simple and extensible library for masking sensitive data in Go structs. Use struct tags to control how fields are masked — passwords, emails, IDs, credit cards, and more.

- [Go Masker v2](#go-masker-v2)
  - [Install](#install)
  - [Quick Start](#quick-start)
  - [Masker Types](#masker-types)
    - [Masking Slices](#masking-slices)
    - [Masking Nested Structs](#masking-nested-structs)
    - [Masking Maps](#masking-maps)
  - [Custom Mask Character](#custom-mask-character)
  - [Abuse Masker](#abuse-masker)
    - [Basic Usage](#basic-usage)
    - [Load Words from File](#load-words-from-file)
    - [Use with Struct Tags](#use-with-struct-tags)
  - [Custom Masker](#custom-masker)
  - [Contributors](#contributors)

## Install

```bash
go get -u github.com/ggwhite/go-masker/v2
```

## Quick Start

```go
package main

import (
    "log"
    masker "github.com/ggwhite/go-masker/v2"
)

type User struct {
    Name     string `mask:"name"`
    Email    string `mask:"email"`
    Password string `mask:"password"`
    Mobile   string `mask:"mobile"`
}

func main() {
    u := &User{
        Name:     "John Doe",
        Email:    "john@gmail.com",
        Password: "secret",
        Mobile:   "0987654321",
    }

    m := masker.NewMaskerMarshaler()
    masked, err := m.Struct(u)
    if err != nil {
        log.Fatal(err)
    }
    log.Println(masked) // &{J**n D**e joh****@gmail.com ************** 0987***321}
}
```

You can also use the package-level default instance:

```go
masked, err := masker.DefaultMaskerMarshaler.Struct(u)
```

## Masker Types

| Tag         | Description                                                                                                | Example Input                    | Example Output                  |
| ----------- | ---------------------------------------------------------------------------------------------------------- | -------------------------------- | ------------------------------- |
| `none`      | No masking, return as-is                                                                                   | `foo`                            | `foo`                           |
| `password`  | Always returns 14 asterisks                                                                                | `secret`                         | `**************`                |
| `name`      | Masks middle characters                                                                                    | `John Doe`                       | `J**n D**e`                     |
| `addr`      | Masks last 6 characters                                                                                    | `台北市內湖區內湖路一段737巷1號` | `台北市內湖區內湖路一段7******` |
| `email`     | Keeps first 3 chars and domain                                                                             | `john@gmail.com`                 | `joh****@gmail.com`             |
| `mobile`    | Masks 3 digits from 4th position                                                                           | `0987654321`                     | `0987***321`                    |
| `tel`       | Formats and masks last 4 digits                                                                            | `0227993078`                     | `(02)2799-****`                 |
| `id`        | Masks digits 7–10                                                                                          | `A123456789`                     | `A12345****`                    |
| `credit`    | Masks digits 7–12                                                                                          | `4111111111111111`               | `411111******1111`              |
| `url`       | Masks URL password                                                                                         | `http://user:pass@host`          | `http://user:xxxxx@host`        |
| `abuse`     | Masks abusive words via trie                                                                               | `bad word`                       | `*** word`                      |
| `struct`    | Recursively masks nested struct                                                                            | —                                | —                               |
| `mapstruct` | Recursively masks map values (supports nested maps, structs, pointers, slices, and pointer-to-slice forms) | —                                | —                               |

### Masking Slices

String slices are also supported:

```go
type Foo struct {
    Tags []string `mask:"name"`
}
```

### Masking Nested Structs

```go
type Address struct {
    Street string `mask:"addr"`
}

type User struct {
    Name    string  `mask:"name"`
    Address Address `mask:"struct"`
}
```

### Masking Maps

Use `mask:"mapstruct"` on map fields to recursively mask map values.

`mapstruct` tag limitations and rules:

- Works only on fields whose own tag is exactly `mask:"mapstruct"`.
- Recursion applies to map **values** only; map keys are never masked.
- Recurses through `map`, `struct`, `ptr`, and `slice` combinations (including nested map and pointer-to-slice forms).
- Leaf values that are not struct-like (for example `string`, `int`, `bool`, `time.Time` fields without mask tags) are kept as-is.
- `nil` values are preserved (`nil` map, `nil` pointer, `nil` slice).
- Struct masking still follows normal struct-tag rules: only exported fields are processed, and fields without `mask:"..."` stay unchanged.
- `interface{}` map values are not force-unwrapped by `mapstruct`; if you need masking, prefer concrete types.
- Cyclic object graphs are not specially handled; avoid self-referencing structures in recursive masking paths.

```go
type Item struct {
    ID string `mask:"id"`
}

type Payload struct {
    Items      map[int]Item                      `mask:"mapstruct"`
    Ptrs       map[int]*Item                     `mask:"mapstruct"`
    SliceItems map[int][]Item                    `mask:"mapstruct"`
    PtrSlices  map[int]*[]Item                   `mask:"mapstruct"`
    Nested     map[int]map[string]map[int][]Item `mask:"mapstruct"`
}

m := masker.NewMaskerMarshaler()
masked, _ := m.Struct(Payload{
    Items:      map[int]Item{1: {ID: "A123456789"}},
    Ptrs:       map[int]*Item{1: {ID: "A123456789"}, 2: nil},
    SliceItems: map[int][]Item{1: {{ID: "A123456789"}, {ID: "B223456789"}}},
    PtrSlices: func() map[int]*[]Item {
        x := []Item{{ID: "A123456789"}, {ID: "B223456789"}}
        return map[int]*[]Item{1: &x, 2: nil}
    }(),
    Nested: map[int]map[string]map[int][]Item{
        1: {
            "group": {
                10: {{ID: "C323456789"}},
            },
        },
    },
})
// Items[1].ID => A12345****
// Ptrs[1].ID  => A12345****
// Ptrs[2]     => nil
// SliceItems[1][0].ID => A12345****
// PtrSlices[1][1].ID  => B22345****
// PtrSlices[2]        => nil
// Nested[1]["group"][10][0].ID => C32345****
```

## Custom Mask Character

By default, `*` is used as the mask character. Use `SetMasker` to change it:

```go
m := masker.NewMaskerMarshaler()
m.SetMasker("#")

masked, _ := m.Struct(u)
// Name "John" -> "J##n"
```

## Abuse Masker

The abuse masker uses a trie for efficient word matching and replacement.

### Basic Usage

```go
abuseWords := []string{"bad", "terrible", "awful"}
abuseMasker := masker.NewAbuseMaskerWithWords(abuseWords)

text := "This is a bad and terrible situation"
masked := abuseMasker.Marshal("*", text)
// Output: "This is a *** and ******** situation"
```

### Load Words from File

```go
loader := masker.NewAbuseWordLoader()
words, err := loader.LoadFromFile("abuse_words.txt")
if err != nil {
    log.Fatal(err)
}
abuseMasker := masker.NewAbuseMaskerWithWords(words)
```

### Use with Struct Tags

```go
type Post struct {
    Title   string `mask:"name"`
    Content string `mask:"abuse"`
}

m := masker.NewMaskerMarshaler()
m.Register(masker.MaskerTypeAbuse, masker.NewAbuseMaskerWithWords([]string{"bad"}))

masked, _ := m.Struct(&Post{Title: "Hello", Content: "bad content"})
```

## Custom Masker

Implement the `Masker` interface to create your own masker:

```go
type Masker interface {
    Marshal(maskChar string, value string) string
}
```

Example:

```go
type SSNMasker struct{}

func (m *SSNMasker) Marshal(s, i string) string {
    if len(i) != 9 {
        return i
    }
    return "***-**-" + i[7:]
}

m := masker.NewMaskerMarshaler()
m.Register("ssn", &SSNMasker{})

type Person struct {
    SSN string `mask:"ssn"`
}
```

## Contributors

Thanks to all the people who already contributed!

<a href="https://github.com/ggwhite/go-masker/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=ggwhite/go-masker" />
</a>
