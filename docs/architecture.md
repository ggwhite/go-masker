# Architecture

## Module

```
github.com/ggwhite/go-masker/v3
```

Go 1.22+. Zero external dependencies.

## Design

go-masker is a single-package library. Each masker type lives in its own file, implementing a shared `Masker` interface. `MaskerMarshaler` is the central orchestrator that maps struct tags to maskers and uses reflection to produce masked copies.

```
User struct (with `mask` tags)
        │
        ▼
MaskerMarshaler.Struct()
        │
        ├─ reflect: iterate fields
        │   ├─ read `mask:"<type>"` tag
        │   ├─ lookup Masker by type
        │   └─ call Masker.Marshal(maskChar, value)
        │
        └─ return new struct (masked copy)
```

## Key Types

| Type | File | Role |
|------|------|------|
| `Masker` | `masker.go` | Interface — `Marshal(maskChar, value string) string` |
| `MaskerMarshaler` | `masker.go` | Registry + reflection engine |
| `MaskerType` | `masker.go` | String alias for tag values |
| `AbuseMasker` | `abuse.go` | Trie-based word masker (stateful, needs word list) |
| `AbuseWordLoader` | `abuse_loader.go` | Loads word lists from file or `io.Reader` |

## File Layout

```
├── masker.go              # Core: MaskerMarshaler, Struct(), Marshal(), interface
├── generic.go             # AllMasker, first-N/last-N dynamic patterns
├── password.go            # PasswordMasker
├── name.go                # NameMasker
├── address.go             # AddressMasker
├── email.go               # EmailMasker
├── mobile.go              # MobileMasker
├── telephone.go           # TelephoneMasker
├── id.go                  # IDMasker
├── credit.go              # CreditMasker
├── url.go                 # URLMasker
├── none.go                # NoneMasker
├── abuse.go               # AbuseMasker (trie)
├── abuse_loader.go        # AbuseWordLoader
├── examples/
│   ├── simple/            # Basic struct masking
│   └── customize/         # Custom masker + override
└── docs/
    └── design/            # Design proposals
```

## Concurrency

`MaskerMarshaler` is safe for concurrent use. All public methods acquire `sync.RWMutex`. Direct field access to `Maskers` map is NOT safe — use `Register`/`Get`.

## Extension Point

Implement `Masker` interface → `Register()` with a custom `MaskerType` string → use in struct tags.
