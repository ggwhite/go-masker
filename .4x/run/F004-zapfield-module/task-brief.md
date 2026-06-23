# Task Brief — F004 zapfield sub-module: zap 整合

## Goal

建立獨立的 Go sub-module `github.com/ggwhite/go-masker/v3/zapfield`，提供：

1. 每個 masker type 對應的 zap `Field` helper，讓 zap 使用者直接寫
   `zapfield.Phone("phone", rawValue)` 取得已遮罩的 `zap.Field`，免去手動 `masker.Mobile(...)` + `zap.String(...)` 的轉接。
2. `Sensitive[T]` adapter：`zapfield.Sensitive("phone", user.Phone)` 把 `masker.Sensitive[T]` 直接轉成已遮罩的 `zap.Field`，使用者不需手動 `.String()`。
3. 獨立的 `go.mod`，使 `go.uber.org/zap` 依賴**只**進 zapfield，不污染 core（`v3/go.mod` 不得出現 zap）。

設計不變式（跨所有 helper）：

- 所有 helper 一律回傳 `zap.Field`，且一律以 `zap.String(key, masked)` 建構 — **不得**用 `zap.Any`／`zap.Reflect`（避免 reflection、避免洩漏原值）。
- 遮罩值一律委派 v3 core 既有的 package-level 函式（`masker.Mobile`、`masker.Email`…）與 `Sensitive[T].String()`，**不得**自行重寫遮罩邏輯，確保輸出與 v3 core 逐字一致。
- `Sensitive` adapter 只讀 `Sensitive[T]` 快取的 masked 值（透過 `.String()`），**絕不**呼叫 `Reveal()`、絕不碰 raw。

## Tasks

### 1. module-setup — zapfield module 結構

1.1 在 `v3/zapfield/` 建立 `go.mod`：
- module path：`github.com/ggwhite/go-masker/v3/zapfield`
- `go 1.21`（與 v3 core 一致）
- `require github.com/ggwhite/go-masker/v3`
- `require go.uber.org/zap v1.28.0`（module cache 已有此版本）
- 因 v3 尚未發 tag，加 `replace github.com/ggwhite/go-masker/v3 => ../` 指回 core 目錄
- 執行 `go mod tidy` 補齊 `go.sum`

1.2 建立 package doc（建議 `v3/zapfield/doc.go`），說明 package 用途與「依賴隔離」設計意圖。package 名稱：`zapfield`。

### 2. field-helpers — masker type Field helpers（`v3/zapfield/field.go`）

為每個 v3 core package-level 遮罩函式提供對應的 `func(key, value string) zap.Field`，命名與 `Sensitive` 建構子家族一致（mobile 用 `Phone`，與 `masker.NewPhone` 對齊）：

| zapfield 函式 | 委派 core 函式 | 對應 masker type |
|---|---|---|
| `Phone(key, value string) zap.Field` | `masker.Mobile` | mobile |
| `Email(key, value string) zap.Field` | `masker.Email` | email |
| `Password(key, value string) zap.Field` | `masker.Password` | password |
| `Name(key, value string) zap.Field` | `masker.Name` | name |
| `Address(key, value string) zap.Field` | `masker.Address` | addr |
| `ID(key, value string) zap.Field` | `masker.ID` | id |
| `Credit(key, value string) zap.Field` | `masker.Credit` | credit |
| `Tel(key, value string) zap.Field` | `masker.Tel` | tel |
| `URL(key, value string) zap.Field` | `masker.URL` | url |
| `Abuse(key, value string) zap.Field` | `masker.Abuse` | abuse |
| `None(key, value string) zap.Field` | `masker.None` | none |
| `All(key, value string) zap.Field` | `masker.All` | all |

每個實作形如：

```go
func Phone(key, value string) zap.Field {
    return zap.String(key, masker.Mobile(value))
}
```

所有 exported helper 須加繁體中文 GoDoc（第一句以函式名開頭），對外 helper 附簡短使用範例。

### 3. sensitive-adapter — Sensitive[T] adapter（`v3/zapfield/sensitive.go`）

提供泛型 adapter，把任意 `masker.Sensitive[T]` 轉成已遮罩的 `zap.Field`：

```go
func Sensitive[T any](key string, s masker.Sensitive[T]) zap.Field {
    return zap.String(key, s.String())
}
```

- 透過 `s.String()` 取快取 masked 值（與 `fmt`/`json`/`slog` 路徑同源）。
- zero-value `Sensitive[T]`（未經建構子、mask 為 nil）的 `.String()` 回空字串 → adapter 輸出 `zap.String(key, "")`，不洩漏 raw（安全退化）。

對外加繁體中文 GoDoc + 範例（示範搭配 `masker.NewPhone(...)`）。

## Scope（要新增／修改的檔案）

- 新增 `v3/zapfield/go.mod`
- 新增 `v3/zapfield/go.sum`（`go mod tidy` 產生）
- 新增 `v3/zapfield/doc.go`（package doc）
- 新增 `v3/zapfield/field.go`（task 2）
- 新增 `v3/zapfield/sensitive.go`（task 3）
- 新增 `v3/zapfield/field_test.go`、`v3/zapfield/sensitive_test.go`（測試）

## Out of Scope

- zap `Core` wrapper／keyword 攔截（屬 F005-zap-core-interceptor，本 feature 不做）。
- 修改 v3 core 任何檔案（core 不得新增 zap 依賴）。
- 新增 `Mask`（動態 type）對應的 zapfield helper — 動態型別非本 feature 重點，需要時另案。
- slog / logr 等其他 logging 框架的 adapter。
- 為 zapfield 發布獨立 release tag。
