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
| `tel-<regionLen>-<numberLen>[-<sep>]` | Configurable-length phone masking, dash/space separator | `mask:"tel-2-8"` on `0227993078` → `02-2799-****` |
| `tel-<intlLen>-<regionLen>-<numberLen>[-<sep>]` | Same, with leading international code | `mask:"tel-2-3-8"` on `8675588888888` → `+86-755-8888-****` |
| `mobile-<keepFront>-<keepEnd>` | Keep first F and last E chars, mask middle | `mask:"mobile-3-4"` on `09012345678` → `090****5678` |
| `id-<keepFront>-<keepEnd>` | Keep first F and last E chars, mask middle | `mask:"id-0-4"` on `123456789` → `*****6789` |
| `mid-<keepFront>-<keepEnd>` | Generic version of mobile-/id- | `mask:"mid-4-4"` on `sk-abc123xyz` → `sk-a***3xyz` |

### tel- dynamic tag

Grammar: `tel-` followed by 2–4 `-`-separated tokens.

- `regionLen`, `intlLen`: positive integers (digit counts)
- `numberLen`: integer `>= 4` (includes the masked last 4 digits)
- `sep`: keyword `dash` (default, outputs `-`) or `space` (outputs ` `)

Disambiguation for 3-token tags: if the 3rd token is `dash`/`space`, it's
`[regionLen, numberLen, sep]`; if it parses as a positive integer, it's
`[intlLen, regionLen, numberLen]`. A keyword and a valid integer can never
be the same string, so this is unambiguous.

Splitting: the cleaned input (formatting chars and leading `+` stripped,
same cleanup as `tel`) must have length exactly
`intlLen + regionLen + numberLen`. Any other length — including one digit
too many — is treated as invalid and the cleaned value is returned
unmasked. **No trunk-prefix (leading domestic `0`) guessing is performed.**
Phone number normalization (e.g. converting `0928xxxxxx` to the E.164-style
`928xxxxxx` before adding a country code) is the caller's responsibility —
it's country-specific and out of scope for a masking library.

| Tag | Input (caller-normalized) | Output |
|-----|---------------------------|--------|
| `tel-2-8` | `0227993078` | `02-2799-****` |
| `tel-3-8-space` | `75588888888` | `755 8888-****` |
| `tel-2-3-8` | `8675588888888` | `+86-755-8888-****` |

See [`docs/design/F013-tel-configurable-region-spec.md`](design/F013-tel-configurable-region-spec.md)
for the full design rationale.

### mobile- dynamic tag

Grammar: `mobile-<keepFront>-<keepEnd>`

- `keepFront`: non-negative integer — number of leading characters to keep
- `keepEnd`: non-negative integer — number of trailing characters to keep
- `keepFront` and `keepEnd` cannot both be 0

Everything between the kept portions is replaced with the mask character.
If `keepFront + keepEnd >= len(value)`, the value is returned unchanged
(nothing to mask).

| Tag | Input | Country | Output |
|-----|-------|---------|--------|
| `mobile-3-4` | `09012345678` | Japan (11 digits) | `090****5678` |
| `mobile-3-4` | `2025551234` | USA (10 digits) | `202***1234` |
| `mobile-0-4` | `447911123456` | UK (12 digits) | `********3456` |
| `mobile-4-0` | `0987654321` | Custom | `0987******` |

### id- dynamic tag

Grammar: `id-<keepFront>-<keepEnd>` — same semantics as `mobile-`.

| Tag | Input | Country | Output |
|-----|-------|---------|--------|
| `id-0-4` | `123456789` | USA SSN (9 digits) | `*****6789` |
| `id-4-0` | `123456789012` | Japan My Number (12 digits) | `1234********` |
| `id-3-3` | `S1234567D` | Singapore NRIC (9 chars) | `S12***67D` |

### mid- dynamic tag

Grammar: `mid-<keepFront>-<keepEnd>` — generic version of `mobile-` and
`id-`, same semantics. Use when the field doesn't fit the mobile/id
category (API keys, tokens, account numbers, etc.).

| Tag | Input | Output |
|-----|-------|--------|
| `mid-4-4` | `sk-abc123xyz` | `sk-a***3xyz` |
| `mid-2-2` | `ABCDEF` | `AB**EF` |
| `mid-1-0` | `secret` | `s*****` |

See [`docs/design/F014-mobile-id-configurable-format-spec.md`](design/F014-mobile-id-configurable-format-spec.md)
for the full design rationale.

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
