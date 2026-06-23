# Masker Types

Complete reference for all built-in masker types.

## Static Types

| Tag | File | Behavior | Example |
|-----|------|----------|---------|
| `none` | `none.go` | Returns value unchanged | `foo` → `foo` |
| `password` | `password.go` | Always 14 mask chars | `secret` → `**************` |
| `name` | `name.go` | Masks middle chars per word | `John Doe` → `J**n D**e` |
| `addr` | `address.go` | Masks last 6 chars | `台北市內湖區內湖路一段737巷1號` → `台北市內湖區內湖路一段7******` |
| `email` | `email.go` | Keeps first 3 + domain | `john@gmail.com` → `joh****@gmail.com` |
| `mobile` | `mobile.go` | Masks 3 digits from pos 4 | `0987654321` → `0987***321` |
| `tel` | `telephone.go` | Formats `(XX)XXXX-****` | `0227993078` → `(02)2799-****` |
| `id` | `id.go` | Masks positions 7–10 | `A123456789` → `A12345****` |
| `credit` | `credit.go` | Masks positions 7–12 | `4111111111111111` → `411111******1111` |
| `url` | `url.go` | Masks URL password | `http://u:p@host` → `http://u:xxxxx@host` |
| `all` | `generic.go` | Replaces every char | `secret` → `******` |
| `abuse` | `abuse.go` | Trie-based word match | `bad word` → `*** word` |

## Dynamic Tags

Handled by `parseGenericMask()` in `generic.go`. No registration needed.

| Pattern | Behavior | Example |
|---------|----------|---------|
| `first-N` | Masks first N chars | `mask:"first-3"` on `ABCDEF` → `***DEF` |
| `last-N` | Masks last N chars | `mask:"last-4"` on `ABCDEFGH` → `ABCD****` |

## Structural Tags

These control recursion, not masking logic:

| Tag | Behavior |
|-----|----------|
| `struct` | Recursively masks nested struct or pointer-to-struct |
| `mapstruct` | Recursively masks map values (supports nested maps, structs, ptrs, slices) |

### mapstruct rules

- Only on fields tagged `mask:"mapstruct"`
- Map keys are never masked
- Recurses through map → struct → ptr → slice combinations
- Leaf scalars without mask tags kept as-is
- nil values preserved
- No cycle detection

## name masker edge cases

| Input | Output | Rule |
|-------|--------|------|
| `""` | `""` | Empty string |
| `"A"` | `"*"` | Single char → all masked |
| `"AB"` | `"A*"` | Two chars → mask second |
| `"John"` | `"J**n"` | Standard → keep first and last |
| `"John Doe"` | `"J**n D**e"` | Space-separated → per-word |

## abuse masker

Stateful — must be initialized with words before use:

```go
// Standalone
abuseMasker := masker.NewAbuseMaskerWithWords([]string{"bad", "terrible"})
result := abuseMasker.Marshal("*", "a bad day")  // "a *** day"

// With struct tags
m := masker.NewMaskerMarshaler()
m.Register(masker.MaskerTypeAbuse, masker.NewAbuseMaskerWithWords(words))
```

Word list loading:

```go
loader := masker.NewAbuseWordLoader()
words, _ := loader.LoadFromFile("words.txt")       // one word per line
words, _ := loader.LoadFromReader(someReader)       // io.Reader
```
