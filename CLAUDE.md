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

Or directly:

```bash
go test ./...
```

## Release

使用 [Semantic Versioning](https://semver.org/)：

| 版本類型 | 條件 | 範例 |
|----------|------|------|
| patch `x.y.Z` | bug fix、文件更新、不影響 API 的內部改動 | `2.2.1` |
| minor `x.Y.0` | 新增功能、新增 masker type、向後相容的 API 擴充 | `2.3.0` |
| major `X.0.0` | 破壞性變更（移除 API、改變行為、升級 module path） | `3.0.0` |

### 發布流程

```bash
gh release create vX.Y.Z --title "Release vX.Y.Z" --notes "..."
```

release notes 格式：

```
## What's New
- 新功能描述

## Bug Fixes
- 修正描述

## Breaking Changes（如有）
- 描述
```

## Key Types

- `Masker` interface — implement to create a custom masker
- `MaskerMarshaler` — manages maskers, exposes `Struct()` and `Marshal()`
- `MaskerType` — string alias for masker tag values
- `AbuseMasker` — trie-based masker, must be initialized with words
- `AbuseWordLoader` — loads word lists from file or `io.Reader`
