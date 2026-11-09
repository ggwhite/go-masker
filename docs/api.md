# Public API Reference (v3)

```go
import masker "github.com/ggwhite/go-masker/v3"
```

## MaskerMarshaler

```go
func NewMaskerMarshaler(opts ...Option) *MaskerMarshaler
var DefaultMaskerMarshaler = NewMaskerMarshaler()
```

Creates an instance with all built-in maskers registered. `DefaultMaskerMarshaler` is the package-level instance used by the convenience functions.

### Options

| Option | Description |
|--------|-------------|
| `WithMaskChar(c rune)` | Mask character for maskers created by this marshaler (default `'*'`). Does not affect `DefaultMaskerMarshaler`. `URLMasker` (fixed `xxxxx` password) and `NoneMasker` ignore it. |

### Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `Struct` | `(s interface{}) (interface{}, error)` | Masks a struct or struct pointer by its `mask` tags and returns a pointer to a masked copy. The input is not modified. |
| `Format` | `(s interface{}) string` | Deterministic masked string of `s` (Go `%v` style, `&{...}` for pointers, `<nil>` for nil pointers) without allocating a masked struct. Same masking as `Struct`. Never panics; nil input returns `""`. |
| `Marshal` | `(t MaskerType, value string) (string, error)` | Masks one value by masker type; error if the type is not registered. |
| `MustMarshal` | `(t MaskerType, value string) string` | Like `Marshal`, but panics if the type is not registered. |
| `Register` | `(t MaskerType, masker Masker)` | Adds or overrides a masker. |
| `Unregister` | `(t MaskerType)` | Removes a masker. |
| `Get` | `(t MaskerType) (Masker, error)` | Returns a registered masker. |
| `List` | `() []MaskerType` | Lists registered masker types. |

All exported methods are safe for concurrent use.

Limitation: `Struct` does not detect cycles; a self-referencing type under `mask:"struct"` or `mask:"mapstruct"` can recurse forever.

## Convenience functions

Single-value helpers backed by `DefaultMaskerMarshaler`:

```go
func Mask(t MaskerType, value string) string // unknown type returns value unchanged; supports first-N / last-N
func None(value string) string
func Password(value string) string
func Name(value string) string
func Address(value string) string
func Email(value string) string
func Mobile(value string) string
func Tel(value string) string
func ID(value string) string
func Credit(value string) string
func URL(value string) string
func Abuse(value string) string
func All(value string) string
```

## Masker interface

```go
type Masker interface {
    Mask(value string) string
}
```

The mask character is injected when a masker is built, so `Mask` takes only the value. Implement this interface for a custom masker and register it with `MaskerMarshaler.Register`.

Built-in implementations: `NoneMasker`, `PasswordMasker`, `NameMasker`, `AddressMasker`, `EmailMasker`, `MobileMasker`, `TelephoneMasker`, `IDMasker`, `CreditMasker`, `URLMasker`, `AbuseMasker`, `AllMasker`.

## MaskerType constants

```go
const (
    TypeNone      MaskerType = "none"
    TypePassword  MaskerType = "password"
    TypeName      MaskerType = "name"
    TypeAddress   MaskerType = "addr"
    TypeEmail     MaskerType = "email"
    TypeMobile    MaskerType = "mobile"
    TypeTel       MaskerType = "tel"
    TypeID        MaskerType = "id"
    TypeCredit    MaskerType = "credit"
    TypeURL       MaskerType = "url"
    TypeAbuse     MaskerType = "abuse"
    TypeStruct    MaskerType = "struct"
    TypeAll       MaskerType = "all"
    TypeMapStruct MaskerType = "mapstruct"
)
```

Dynamic tags (`first-N`, `last-N`, `mobile-`, `id-`, `mid-`, `tel-`) are described in [masker-types.md](masker-types.md).

## Sensitive[T]

A generic wrapper that masks on every automatic output path. The raw value is only available through `Reveal`.

```go
func NewSensitive[T any](raw T, maskFn func(T) string, opts ...SensitiveOption) Sensitive[T]

func NewPassword(v string, opts ...SensitiveOption) Sensitive[string]
func NewName(v string, opts ...SensitiveOption) Sensitive[string]
func NewAddress(v string, opts ...SensitiveOption) Sensitive[string]
func NewEmail(v string, opts ...SensitiveOption) Sensitive[string]
func NewPhone(v string, opts ...SensitiveOption) Sensitive[string]
func NewTel(v string, opts ...SensitiveOption) Sensitive[string]
func NewID(v string, opts ...SensitiveOption) Sensitive[string]
func NewCredit(v string, opts ...SensitiveOption) Sensitive[string]
func NewURL(v string, opts ...SensitiveOption) Sensitive[string]
```

| Method | Output |
|--------|--------|
| `String`, `GoString` | masked value (`fmt`, `%v`, `%#v`) |
| `MarshalJSON`, `MarshalText` | masked value |
| `LogValue` | masked value (`log/slog`) |
| `Reveal() T` | raw value — the only way out |
| `Equal(other) bool` | compares raw values |
| `UnmarshalJSON` | stores the raw value and recomputes the mask |
| `Value() (driver.Value, error)` | raw value for `database/sql` writes |
| `Scan(src any) error` | reads a DB value; `NULL` gives the zero value; an unbound zero-value receiver keeps the mask empty instead of leaking the raw value |

The masked value is computed once at construction and cached.

### SensitiveOption

| Option | Description |
|--------|-------------|
| `WithRedact()` | Every automatic output returns `DefaultRedactText` (`[REDACTED]`) instead of a partial mask. `Reveal` and `Value` still return the raw value. |
| `WithRedactText(text string)` | Same as `WithRedact`, with custom replacement text. |

## Abuse words

```go
func NewAbuseMasker(maskChar string) *AbuseMasker
func NewAbuseMaskerWithWords(maskChar string, words []string) *AbuseMasker
```

| Method | Description |
|--------|-------------|
| `AddWord(word string)` / `AddWords(words []string)` | Add dictionary words. With an empty dictionary, text is returned unchanged. |
| `Mask(text string) string` | Masks dictionary hits. |
| `ContainsAbuse(text string) bool` | Reports whether any word matches. |
| `GetAbuseWords(text string) []string` | Returns the matched words. |

`AbuseTrie` (`NewAbuseTrie`, `Insert`, `InsertAll`, `Contains`) is the underlying trie.

```go
func NewAbuseWordLoader() *AbuseWordLoader
func (l *AbuseWordLoader) LoadFromFile(filename string) ([]string, error)
func (l *AbuseWordLoader) LoadFromReader(reader io.Reader) ([]string, error)
func (l *AbuseWordLoader) LoadFromString(content string) ([]string, error)
func (l *AbuseWordLoader) LoadFromSlice(words []string) []string
```

## Integrations

### slogfield — `github.com/ggwhite/go-masker/v3/slogfield`

`slog.Attr` helpers that mask before logging:

```go
func Masked(key string, maskerType masker.MaskerType, value string) slog.Attr
func SensitiveAttr[T any](key string, s masker.Sensitive[T]) slog.Attr
// plus Password, Name, Address, Email, Phone, Tel, ID, Credit, URL, Abuse, All, None (key, value string)
```

### zapfield — `github.com/ggwhite/go-masker/v3/zapfield`

`zap.Field` helpers with the same names (`Phone`, `Email`, ..., `Sensitive[T]`), plus a core wrapper that masks fields by key:

```go
func WrapCore(core zapcore.Core, rules InterceptRules) zapcore.Core

type InterceptRules struct {
    Keywords []string          // case-insensitive substring match on the field key → TypeAll
    Patterns []*regexp.Regexp  // regex match on the field key → TypeAll
    Rules    []Rule            // checked after Keywords/Patterns; first match wins
}

type Rule struct {
    Keywords   []string
    Patterns   []*regexp.Regexp
    MaskerType masker.MaskerType // zero value means TypeAll
}
```

Key-name matching guesses which fields are sensitive, so it has false positives and false negatives. Use it as a fallback when business code cannot change; prefer `Sensitive[T]` or the explicit field helpers.

### ginmasker — `github.com/ggwhite/go-masker/v3/ginmasker`

```go
func Middleware(opts ...Option) gin.HandlerFunc
```

Masks sensitive fields in request/response bodies, query and headers before writing the access log. The real request and response are not changed. Bodies are read only when the logger's Info level is enabled.

| Option | Description |
|--------|-------------|
| `WithLogger(logger *zap.Logger)` | Logger for the access log. |
| `WithKeywords(keywords ...string)` | Body field keys to mask. |
| `WithQueryMask(keys ...string)` | Query parameters to mask. |
| `WithHeaderMask(keys ...string)` | Headers to mask. |
| `WithMaxBodySize(n int64)` | Maximum body size to read for logging. |
