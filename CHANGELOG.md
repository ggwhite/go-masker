# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Features

- **v3 module** — 全新 `github.com/ggwhite/go-masker/v3`（Go 1.21+），包含：
  - 新 `Masker` interface：`Mask(value string) string`（移除 mask char 參數）
  - Functional option `WithMaskChar` 配置遮罩字元
  - Package-level 便利函式：`masker.Mobile()`、`masker.Email()` 等
  - `MaskerType` 精簡命名（`TypeMobile` 取代 `MaskerTypeMobile`）
- **Sensitive[T] 泛型安全型別** — 把遮罩從「開發者記得做」變成「洩漏要刻意寫 `.Reveal()`」，自動安全輸出至 fmt/json/slog
- **zapfield sub-module** — `github.com/ggwhite/go-masker/v3/zapfield`，提供 zap Field helpers 與 `Sensitive[T]` adapter
- **zap Core interceptor** — `WrapCore` 以 keyword/regex 攔截 log field 的最後一道防線
- **v2 Struct() reflect cache** — 首次解析 type metadata 後快取，重複呼叫效能提升 ≥ 50%
- **v2 Format()** — 直接輸出遮罩後字串，不配置新 struct，適合 log/debug 場景
- **mapstruct tag** — 遞迴遮罩 map values，支援巢狀 map/struct/ptr/slice 組合

### Fixes

- **v3 abuse trie data race** — AbuseMasker 加 sync.RWMutex，防止 AddWord 與 Mask 併發存取
- **v3 interface{} 欄位歸零** — Struct() 不再靜默歸零非 struct tag 的 interface 欄位，改為遮罩底層 string
- **v3 maskMapStructValue 漏處理 Interface** — 加入 reflect.Interface case 解包遞迴
- **v3 email Split 丟失多 '@' 後段** — 改用 SplitN 保留完整 domain
- **v3 Sensitive UnmarshalJSON nil mask** — mask 未綁定時回傳 error 而非靜默空值
- **v3 Maskers 欄位 data race** — unexport 為 maskers，強制走 Register/Get
- **v2 Format() 洩漏敏感資料** — Marshal 失敗時輸出 mask chars 而非原始值
- **v2 十項 code review bug 修正** — 含 Struct() non-struct panic、nil handling 等

### Internal

- **4x 開發流程** — 設定 features（F001–F006）、settings.json、profiles、roles
- **LLM wiki** — `docs/` 目錄結構化管理（AGENTS.md 治理 + architecture/api/masker-types 文件）
- **Claude Code 專案設定** — `.claude/settings.json` 權限、release skill

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
