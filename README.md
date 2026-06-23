# Go Masker

[![build workflow](https://github.com/ggwhite/go-masker/actions/workflows/go.yml/badge.svg)](https://github.com/ggwhite/go-masker/actions)
[![GoDoc](https://godoc.org/github.com/ggwhite/go-masker?status.svg)](https://godoc.org/github.com/ggwhite/go-masker)
[![Go Report Card](https://goreportcard.com/badge/github.com/ggwhite/go-masker)](https://goreportcard.com/report/github.com/ggwhite/go-masker)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/ggwhite/go-masker/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/ggwhite/go-masker.svg?style=flat-square)](https://github.com/ggwhite/go-masker/releases/latest)

Go Masker is a simple and extensible library for masking sensitive data in Go structs. Use struct tags to control how fields are masked — passwords, emails, IDs, credit cards, and more.

> Looking for v2? See the [`release/v2`](https://github.com/ggwhite/go-masker/tree/release/v2) branch.

## Install

```bash
go get -u github.com/ggwhite/go-masker/v3
```

Requires Go 1.22+.

## Quick Start

```go
package main

import (
    "fmt"

    masker "github.com/ggwhite/go-masker/v3"
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

    // Option 1: Get a masked copy of the struct
    m := masker.NewMaskerMarshaler()
    masked, err := m.Struct(u)
    if err != nil {
        panic(err)
    }
    fmt.Println(masked) // &{J**n D**e joh****@gmail.com ************** 0987***321}

    // Option 2: Get a masked string directly (great for logging)
    fmt.Println(m.Format(u)) // &{J**n D**e joh****@gmail.com ************** 0987***321}
}
```

You can also use the package-level default instance:

```go
masked, err := masker.DefaultMaskerMarshaler.Struct(u)
```

## Sensitive[T] — Type-Safe Masking

`Sensitive[T]` wraps a value so that all automatic output (fmt, JSON, slog) is masked. You must call `.Reveal()` to get the original value — leaking data requires an explicit opt-in.

```go
phone := masker.NewPhone("0987654321")

fmt.Println(phone)          // 0987***321
fmt.Println(phone.Reveal()) // 0987654321

b, _ := json.Marshal(phone) // "0987***321"
```

Built-in constructors: `NewPhone`, `NewEmail`, `NewPassword`, `NewID`, `NewCredit`, `NewName`, `NewAddress`, `NewTel`, `NewURL`.

### Redact Mode

For compliance scenarios where partial masking is not enough:

```go
phone := masker.NewPhone("0987654321", masker.WithRedact())
fmt.Println(phone)          // [REDACTED]
fmt.Println(phone.Reveal()) // 0987654321

// Custom redact text
phone2 := masker.NewPhone("0987654321", masker.WithRedactText("***"))
fmt.Println(phone2)         // ***
```

### Database Integration

`Sensitive[T]` implements `sql.Scanner` and `driver.Valuer`, so it works directly with GORM and other ORMs:

```go
type Player struct {
    Phone masker.Sensitive[string] `gorm:"column:phone"`
}

// db.Find(&player) → auto scan + bind mask
// fmt.Println(player.Phone) → 0987***321 (masked)
// player.Phone.Reveal() → original value
```

- `Value()` returns the original value (DB stores plaintext)
- `Scan(nil)` sets to zero value without panic

## Format — Log-Friendly Output

`Format()` returns a deterministic masked string without allocating a new struct — ideal for logging and debugging:

```go
m := masker.NewMaskerMarshaler()

type Foo struct {
    Name  string `mask:"name"`
    Email string `mask:"email"`
    Self  *Foo   `mask:"struct"`
}

foo := &Foo{
    Name:  "John Doe",
    Email: "john@gmail.com",
    Self:  &Foo{Name: "Jane Doe", Email: "jane@gmail.com"},
}

fmt.Println(m.Format(foo))
// &{J**n D**e joh****@gmail.com &{J**e D**e jan****@gmail.com <nil>}}
```

- Pointer-to-struct fields are expanded as `&{...}` (no memory address)
- `nil` pointers display as `<nil>`
- Output is deterministic — same input always produces the same string

## Masker Types

| Tag         | Description                      | Example Input                    | Example Output                  |
| ----------- | -------------------------------- | -------------------------------- | ------------------------------- |
| `none`      | No masking, return as-is         | `foo`                            | `foo`                           |
| `password`  | Always returns 14 asterisks      | `secret`                         | `**************`                |
| `name`      | Masks middle characters          | `John Doe`                       | `J**n D**e`                     |
| `addr`      | Masks last 6 characters          | `台北市內湖區內湖路一段737巷1號` | `台北市內湖區內湖路一段7******` |
| `email`     | Keeps first 3 chars and domain   | `john@gmail.com`                 | `joh****@gmail.com`             |
| `mobile`    | Masks 3 digits from 4th position | `0987654321`                     | `0987***321`                    |
| `tel`       | Formats and masks last 4 digits  | `0227993078`                     | `(02)2799-****`                 |
| `id`        | Masks digits 7–10                | `A123456789`                     | `A12345****`                    |
| `credit`    | Masks digits 7–12                | `4111111111111111`               | `411111******1111`              |
| `url`       | Masks URL password               | `http://user:pass@host`          | `http://user:xxxxx@host`        |
| `all`       | Replaces every character         | `secret`                         | `******`                        |
| `abuse`     | Masks abusive words via trie     | `bad word`                       | `*** word`                      |
| `struct`    | Recursively masks nested struct  | —                                | —                               |
| `mapstruct` | Recursively masks map values     | —                                | —                               |

### Dynamic Mask Tags

Use `first-N` and `last-N` to mask a specific number of characters from the start or end:

```go
type Token struct {
    Code   string `mask:"first-3"`  // masks first 3 chars
    Suffix string `mask:"last-4"`   // masks last 4 chars
}

m := masker.NewMaskerMarshaler()
masked, _ := m.Struct(&Token{Code: "ABC123", Suffix: "secret99"})
// Code:   "***123"
// Suffix: "secr****"
```

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

Use `mask:"mapstruct"` on map fields to recursively mask map values:

- Recurses through `map`, `struct`, `ptr`, and `slice` combinations (including nested map and pointer-to-slice forms)
- Map keys are never masked; only values are processed
- `nil` values are preserved (`nil` map, `nil` pointer, `nil` slice)
- Leaf values without mask tags are kept as-is

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

By default, `*` is used as the mask character. Use `WithMaskChar` to change it:

```go
m := masker.NewMaskerMarshaler(masker.WithMaskChar('#'))

masked, _ := m.Struct(u)
// Name "John" -> "J##n"
```

## Abuse Masker

The abuse masker uses a trie for efficient word matching and replacement.

### Basic Usage

```go
abuseMasker := masker.NewAbuseMaskerWithWords("*", []string{"bad", "terrible", "awful"})

masked := abuseMasker.Mask("This is a bad and terrible situation")
// "This is a *** and ******** situation"
```

### Load Words from File

```go
loader := masker.NewAbuseWordLoader()
words, err := loader.LoadFromFile("abuse_words.txt")
if err != nil {
    log.Fatal(err)
}
abuseMasker := masker.NewAbuseMaskerWithWords("*", words)
```

### Use with Struct Tags

```go
type Post struct {
    Title   string `mask:"name"`
    Content string `mask:"abuse"`
}

m := masker.NewMaskerMarshaler()
m.Register(masker.TypeAbuse, masker.NewAbuseMaskerWithWords("*", []string{"bad"}))

masked, _ := m.Struct(&Post{Title: "Hello", Content: "bad content"})
```

## Custom Masker

Implement the `Masker` interface to create your own masker:

```go
type Masker interface {
    Mask(value string) string
}
```

Example:

```go
type SSNMasker struct{}

func (m *SSNMasker) Mask(value string) string {
    if len(value) != 9 {
        return value
    }
    return "***-**-" + value[7:]
}

m := masker.NewMaskerMarshaler()
m.Register("ssn", &SSNMasker{})

type Person struct {
    SSN string `mask:"ssn"`
}
```

## Sub-Modules

Optional sub-modules for integrating masking into your logging and HTTP stack. Each is a separate `go get` — they don't add dependencies to the core module.

### zapfield — Zap Logger Integration

```bash
go get github.com/ggwhite/go-masker/v3/zapfield
```

```go
import "github.com/ggwhite/go-masker/v3/zapfield"

logger.Info("user login",
    zapfield.Phone("phone", "0987654321"),  // 0987***321
    zapfield.Email("email", "john@g.com"),  // joh****@g.com
)
```

WrapCore intercepts log fields by keyword, with per-rule masker type support:

```go
core = zapfield.WrapCore(core, zapfield.InterceptRules{
    Rules: []zapfield.Rule{
        {Keywords: []string{"phone", "mobile"}, MaskerType: masker.TypeMobile},
        {Keywords: []string{"password", "secret"}, MaskerType: masker.TypePassword},
    },
})
```

### slogfield — slog Integration

```bash
go get github.com/ggwhite/go-masker/v3/slogfield
```

```go
import "github.com/ggwhite/go-masker/v3/slogfield"

slog.Info("user login",
    slogfield.Phone("phone", "0987654321"),  // 0987***321
    slogfield.Email("email", "john@g.com"),  // joh****@g.com
)
```

### ginmasker — Gin Middleware

```bash
go get github.com/ggwhite/go-masker/v3/ginmasker
```

Access log middleware that automatically masks sensitive fields in request/response JSON body. Only affects log output — never modifies the actual request or response.

```go
import "github.com/ggwhite/go-masker/v3/ginmasker"

r := gin.New()
r.Use(ginmasker.Middleware(
    ginmasker.WithLogger(logger),
    ginmasker.WithKeywords("password", "token", "secret"),
    ginmasker.WithQueryMask("api_key"),
    ginmasker.WithHeaderMask("Authorization"),
))
```

## Performance

`Struct()` caches type metadata on first call via `sync.Map`, so repeated calls on the same struct type skip reflection overhead. Benchmarks show ≥ 50% improvement on subsequent calls.

`Format()` avoids allocating a new struct entirely — it writes masked output directly to a `strings.Builder`, making it the better choice for log-only use cases.

## Migration from v2

The module path changes from `github.com/ggwhite/go-masker/v2` to `github.com/ggwhite/go-masker/v3`. Key API changes:

| v2 | v3 |
|---|---|
| `Masker.Marshal(maskChar, value)` | `Masker.Mask(value)` |
| `MaskerTypeMobile` | `TypeMobile` |
| `m.SetMasker("#")` | `NewMaskerMarshaler(WithMaskChar('#'))` |
| `NewAbuseMaskerWithWords(words)` | `NewAbuseMaskerWithWords("*", words)` |

## Contributors

Thanks to all the people who already contributed!

<a href="https://github.com/ggwhite/go-masker/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=ggwhite/go-masker" />
</a>
