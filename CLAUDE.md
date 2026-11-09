@.4x/plugins/CLAUDE.md

# go-masker

A Go library for masking sensitive data via struct tags.

## Module

```
github.com/ggwhite/go-masker/v3
```

Requires Go 1.22+.

## Project Structure

Core files:

| File | Contents |
|------|----------|
| `masker.go` | `Masker` interface, `MaskerType` constants (`TypeXxx`), `MaskerMarshaler` (`Struct`, `Marshal`, `MustMarshal`, `Register`, `Unregister`, `Get`, `List`), `NewMaskerMarshaler(opts ...Option)`, `DefaultMaskerMarshaler` |
| `format.go` | `MaskerMarshaler.Format` — masked string without building a masked struct |
| `type_cache.go` | Per-type field metadata cache used by `Struct` / `Format` |
| `generic.go` | `AllMasker` and the dynamic tags (`first-N`, `last-N`, `mobile-`, `id-`, `mid-`, `tel-`) via `parseGenericMask()` |
| `convenience.go` | Package-level one-value functions (`Mask`, `Mobile`, `Email`, ...) backed by `DefaultMaskerMarshaler` |
| `sensitive*.go` | `Sensitive[T]` safe type, its constructors (`NewPhone`, ...), options (`WithRedact`, `WithRedactText`) and `database/sql` support |
| `abuse_loader.go` | `AbuseWordLoader` — loads word lists from file, reader, string or slice |

Each built-in masker has its own file (full behavior and examples: `docs/masker-types.md`):

| File | Tag | Behavior |
|------|-----|----------|
| `none.go` | `none` | returns value unchanged |
| `password.go` | `password` | always 14 mask chars |
| `name.go` | `name` | masks middle chars per word |
| `address.go` | `addr` | keeps first 6 chars, rest → `******` |
| `email.go` | `email` | keeps first 3 chars + domain |
| `mobile.go` | `mobile` | masks 3 digits from position 4 (`0987***321`) |
| `telephone.go` | `tel` | formats `(XX)XXXX-****` |
| `id.go` | `id` | masks positions 7–10 |
| `credit.go` | `credit` | masks positions 7–12 |
| `url.go` | `url` | masks the URL password as fixed `xxxxx`; ignores the mask char |
| `abuse.go` | `abuse` | trie-based dictionary word masking; empty dictionary returns text unchanged |
| `generic.go` | `all` | replaces every char |

Integrations are separate Go sub-modules (own `go.mod` with `replace => ../`), so the core module keeps zero external dependencies:

| Dir | Module |
|-----|--------|
| `slogfield/` | `log/slog` attr helpers |
| `zapfield/` | `zap.Field` helpers + `WrapCore` key-based interceptor |
| `ginmasker/` | Gin access-log middleware |

Public API reference: `docs/api.md`.

## Adding a New Masker Type

1. Create `<type>.go` implementing the `Masker` interface. The mask char is injected at construction, so `Mask` takes only the value:
   ```go
   type MyMasker struct {
       mask string
   }

   func (m *MyMasker) Mask(value string) string { ... }
   ```
2. Add a `TypeMy MaskerType = "my"` constant in `masker.go`.
3. Register it in the map inside `NewMaskerMarshaler()` with `&MyMasker{mask: c}`. `DefaultMaskerMarshaler` is built by `NewMaskerMarshaler()`, so it needs no separate step.
4. Optionally add a package-level function in `convenience.go` and a `Sensitive` constructor in `sensitive_constructors.go`.
5. Add tests in `<type>_test.go`, and update `docs/masker-types.md`, `docs/api.md` and the README tag table with outputs produced by running the code.

Users can also add a masker at runtime without changing the library: `m.Register("my", &MyMasker{...})`.

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

## 4x Feature 命名規則

檔名、`id`、`name` 統一使用 `F00X-短名稱` 格式：

- 檔名：`.4x/features/F001-v3-core-interface.yaml`
- id：`F001-v3-core-interface`
- name：`"F001 v3 core: 新 Masker interface + package-level 短函式"`
- depends 引用完整 id：`F001-v3-core-interface`

編號按依賴順序遞增，不跳號。新增 feature 接續最大編號。

## Key Types

- `Masker` interface — `Mask(value string) string`; implement to create a custom masker
- `MaskerMarshaler` — manages maskers; `Struct()`, `Format()`, `Marshal()`, `MustMarshal()`; concurrent-safe
- `Option` / `WithMaskChar(rune)` — per-marshaler mask char; never changes `DefaultMaskerMarshaler`
- `MaskerType` — string alias for `mask:"..."` tag values (`TypeMobile`, ...)
- `Sensitive[T]` — every automatic output (fmt, JSON, text, slog) is masked; only `Reveal()` returns the raw value
- `AbuseMasker` — trie-based masker; masks nothing until words are added (`AddWords` / `NewAbuseMaskerWithWords`)
- `AbuseWordLoader` — loads word lists from file, `io.Reader`, string or slice

@.4x/learnings-context.md
