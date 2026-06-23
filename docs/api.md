# Public API Reference

## MaskerMarshaler

```go
func NewMaskerMarshaler() *MaskerMarshaler
```

Creates a new instance with all built-in maskers registered and `*` as default mask char.

```go
var DefaultMaskerMarshaler = NewMaskerMarshaler()
```

Package-level default instance.

### Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `Struct` | `(s interface{}) (interface{}, error)` | Masks struct fields by `mask` tags, returns new masked copy |
| `Marshal` | `(t MaskerType, value string) (string, error)` | Masks a single value by masker type |
| `Register` | `(t MaskerType, masker Masker)` | Adds or overrides a masker |
| `Unregister` | `(t MaskerType)` | Removes a masker |
| `Get` | `(t MaskerType) (Masker, error)` | Retrieves a registered masker |
| `List` | `() []MaskerType` | Lists all registered masker types |
| `SetMasker` | `(masker string)` | Changes mask character |

All methods are concurrent-safe (`sync.RWMutex`).

## Masker Interface

```go
type Masker interface {
    Marshal(maskChar string, value string) string
}
```

Implement this to create a custom masker. Register with `MaskerMarshaler.Register()`.

## MaskerType Constants

```go
MaskerTypeNone      = "none"
MaskerTypePassword  = "password"
MaskerTypeName      = "name"
MaskerTypeAddress   = "addr"
MaskerTypeEmail     = "email"
MaskerTypeMobile    = "mobile"
MaskerTypeTel       = "tel"
MaskerTypeID        = "id"
MaskerTypeCredit    = "credit"
MaskerTypeURL       = "url"
MaskerTypeAbuse     = "abuse"
MaskerTypeStruct    = "struct"
MaskerTypeAll       = "all"
MaskerTypeMapStruct = "mapstruct"
```

## AbuseMasker

```go
func NewAbuseMasker() *AbuseMasker                       // empty trie
func NewAbuseMaskerWithWords(words []string) *AbuseMasker // pre-loaded trie
```

## AbuseWordLoader

```go
func NewAbuseWordLoader() *AbuseWordLoader
func (l *AbuseWordLoader) LoadFromFile(path string) ([]string, error)
func (l *AbuseWordLoader) LoadFromReader(r io.Reader) ([]string, error)
```
