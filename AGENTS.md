@.4x/plugins/copilot-AGENTS.md

@.4x/plugins/codex-AGENTS.md

# go-masker

A Go library for masking sensitive data via struct tags.

## Module

```
github.com/ggwhite/go-masker/v2
```

Requires Go 1.17+.

## Project Structure

Each masker type has its own file:

| File | Masker Type |
|------|-------------|
| `masker.go` | Core: `MaskerMarshaler`, `Struct()`, `Marshal()`, `NewMaskerMarshaler()` |
| `password.go` | `password` — always 14 mask chars |
| `name.go` | `name` — masks middle chars |
| `address.go` | `addr` — masks last 6 chars |
| `email.go` | `email` — keeps first 3 chars + domain |
| `mobile.go` | `mobile` — masks positions 4–6 |
| `telephone.go` | `tel` — formats `(XX)XXXX-****` |
| `id.go` | `id` — masks positions 7–10 |
| `credit.go` | `credit` — masks positions 7–12 |
| `url.go` | `url` — masks URL password via `url.Redacted()` |
| `none.go` | `none` — returns value unchanged |
| `generic.go` | `all`, `first-N`, `last-N` — dynamic masking patterns |
| `abuse.go` | `abuse` — trie-based word masking |
| `abuse_loader.go` | Loads abuse word lists from file/reader |

## Adding a New Masker Type

1. Create `<type>.go` implementing the `Masker` interface:
   ```go
   type MyMasker struct{}
   func (m *MyMasker) Marshal(s, i string) string { ... }
   ```
2. Add a `MaskerType` constant in `masker.go`
3. Register it in `NewMaskerMarshaler()` and `DefaultMaskerMarshaler`
4. Add tests in `<type>_test.go`

## Testing

```bash
make test          # run all tests with race detector and coverage
make cover         # open coverage report in browser
```

## Key Types

- `Masker` interface — `Marshal(maskChar, value string) string`
- `MaskerMarshaler` — manages maskers, exposes `Struct()` and `Marshal()`
- `MaskerType` — string alias for masker tag values
- `AbuseMasker` — trie-based masker, must be initialized with words

## Code Style

- GoDoc comments in Traditional Chinese; first line starts with the identifier name
- Public functions must have GoDoc
- No unnecessary comments — good naming over comments
- No Pinyin naming — use standard English terms
- Follow gofmt/golint, explicit error handling
