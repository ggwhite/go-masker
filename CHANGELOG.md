# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [3.3.0] - 2026-07-07

### Features

- **`mobile-` 動態遮罩 tag** — `mobile-<keepFront>-<keepEnd>` 支援各國手機號碼格式，保留前幾碼與後幾碼、中間全遮罩（如日本 `mobile-3-4`、英國 `mobile-0-4`）；不影響既有 `mobile` tag 行為
- **`id-` 動態遮罩 tag** — `id-<keepFront>-<keepEnd>` 支援各國身分證件格式（如美國 SSN `id-0-4`、日本 My Number `id-4-0`），與 `mobile-` 共用相同的 keepFront/keepEnd 語意；不影響既有 `id` tag 行為
- **`mid-` 動態遮罩 tag** — `mid-<keepFront>-<keepEnd>` 通用版本，適用於 API key、token 等不屬於 mobile/id 類別的欄位

## [3.2.0] - 2026-07-07

### Features

- **`tel-` 動態遮罩 tag** — 支援自訂區碼／號碼／（可選）國際碼位數與分隔符（`dash`/`space`），解決固定 `(XX)XXXX-****` 格式無法涵蓋 3 碼區碼（如中國 `755`）或含國際碼場景的問題（[#40](https://github.com/ggwhite/go-masker/issues/40)）；不影響既有 `tel` tag 行為

## [3.1.0] - 2026-06-23

### Features

- **slogfield sub-module** — `github.com/ggwhite/go-masker/slogfield`，提供 `slog.Attr` helper 函式（Phone、Email、Name 等），與 zapfield 對等的 slog 版本
- **WrapCore 自訂 masker type** — `zapfield.Rule` 新增 `MaskerType` 欄位，攔截時可指定遮罩型別（如 `TypeMobile`），不再只能全遮罩；zero value 維持向後相容
- **Sensitive[T] SQL Scanner/Valuer** — 實作 `sql.Scanner` + `driver.Valuer`，ORM model 敏感欄位可直接宣告為 `Sensitive[string]`，DB 讀寫自動處理
- **ginmasker sub-module** — `github.com/ggwhite/go-masker/ginmasker`，Gin middleware 自動遮罩 access log 中的 request/response body、query、header 敏感欄位
- **Sensitive[T] Redact 模式** — `WithRedact()` 建構選項，輸出 `[REDACTED]` 取代部分遮罩，適用於合規場景與對外 API response

### Internal

- **CI 更新** — GitHub Actions 升級至 Go 1.22、移除舊 v3 子目錄 job、新增 zapfield job

## [3.0.0] - 2026-06-23

### Breaking Changes

- **Module path 升級** — `github.com/ggwhite/go-masker/v2` → `github.com/ggwhite/go-masker/v3`，需更新 import path
- **Masker interface 簡化** — `Marshal(maskChar, value string) string` → `Mask(value string) string`，mask char 改由 `WithMaskChar` option 設定
- **MaskerType 精簡命名** — `MaskerTypeMobile` → `TypeMobile`、`MaskerTypeEmail` → `TypeEmail` 等
- **Go 最低版本** — 1.17 → 1.22

### Features

- **Functional options** — `NewMaskerMarshaler(WithMaskChar("#"))` 取代 `SetMasker()`
- **Package-level 便利函式** — `masker.Mobile()`、`masker.Email()` 等，免建 marshaler 直接遮罩
- **Sensitive[T] 泛型安全型別** — 把遮罩從「開發者記得做」變成「洩漏要刻意寫 `.Reveal()`」，自動安全輸出至 fmt/json/slog
- **zapfield sub-module** — `github.com/ggwhite/go-masker/zapfield`，提供 zap Field helpers 與 `Sensitive[T]` adapter
- **zap Core interceptor** — `WrapCore` 以 keyword/regex 攔截 log field 的最後一道防線

### Fixes

- **abuse trie data race** — AbuseMasker 加 sync.RWMutex，防止 AddWord 與 Mask 併發存取
- **interface{} 欄位歸零** — Struct() 不再靜默歸零非 struct tag 的 interface 欄位，改為遮罩底層 string
- **maskMapStructValue 漏處理 Interface** — 加入 reflect.Interface case 解包遞迴
- **email Split 丟失多 '@' 後段** — 改用 SplitN 保留完整 domain
- **Sensitive UnmarshalJSON nil mask** — mask 未綁定時回傳 error 而非靜默空值
- **Maskers 欄位 data race** — unexport 為 maskers，強制走 Register/Get

## [2.4.2] - 2026-06-23

### Features

- **mapstruct tag** — 遞迴遮罩 map values，支援巢狀 map/struct/ptr/slice/ptr-to-slice 組合 (#38)
- **Struct() reflect cache** — 首次解析 type metadata 後快取（sync.Map），重複呼叫同型別效能提升 ≥ 50%
- **Format()** — 直接輸出遮罩後確定性字串，不配置新 struct，適合 log/debug 場景

### Fixes

- **Format() 敏感資料洩漏** — Marshal 失敗時改輸出 mask chars，不再回傳原始值
- **Format() dead code** — 移除 writeMaskedField 中永遠不會到達的 reflect.String 分支
- **10 項 code review 修正** — 含 Struct() non-struct input panic、nil handling、telephone 格式等

## [2.3.1] - 2026-04-20

### Docs

- **GoDoc 註解改善** — 修正並改進所有 masker type 的 godoc 註解
- **Release 規則** — CLAUDE.md 加入 semantic versioning 規則

## [2.3.0] - 2026-04-20

### Features

- **Generic mask types** — 新增 `all`（全字元遮罩）、`first-N`、`last-N` 動態遮罩 tag (#36)

### Fixes

- **Struct() 不再 panic** — 非 struct 輸入改為回傳 error (#35)

## [2.2.1] - 2026-04-20

### Docs

- **README 更新** — 新增 CLAUDE.md，修正 address.go 的 math.MaxInt 問題

## [2.2.0] - 2025-11-26

### Features

- **Abuse Masker** — 新增 trie-based 敏感詞遮罩功能 (#32)

### Fixes

- **Email 遮罩修正** — 正確套用 email mask (#30)

## [2.1.0] - 2024-09-02

### Fixes

- **Email mask 修正** — 修正 email 遮罩邏輯 (#30)

## [2.0.0] - 2024-03-14

### Features

- **v2 重構** — module path 升級為 `github.com/ggwhite/go-masker/v2`，改為 `MaskerMarshaler` 架構，支援自訂 masker 註冊
