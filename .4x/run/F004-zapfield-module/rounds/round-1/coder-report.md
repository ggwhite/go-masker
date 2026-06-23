# Coder Report — Round 1

## What Was Done

實作 F004：建立獨立 sub-module `github.com/ggwhite/go-masker/v3/zapfield`，提供 zap 整合 helper，並把 zap 依賴隔離在此 sub-module、不污染 v3 core。

完成三項 task：

1. **module-setup**：建立 `v3/zapfield/go.mod`（module path、`go 1.21`、require v3 core + `go.uber.org/zap v1.28.0`、`replace ... => ../`），`go mod tidy` 產生 `go.sum`；新增 `doc.go` 說明 package 用途與依賴隔離設計意圖。
2. **field-helpers**：`field.go` 為 12 個 v3 core package-level 函式各提供 `func(key, value string) zap.Field` helper（Phone/Email/Password/Name/Address/ID/Credit/Tel/URL/Abuse/None/All），命名與 `Sensitive` 建構子家族一致（mobile→`Phone`）。
3. **sensitive-adapter**：`sensitive.go` 提供泛型 `Sensitive[T any](key string, s masker.Sensitive[T]) zap.Field`。

設計不變式全數遵守：
- 所有 helper 一律以 `zap.String(key, masked)` 建構，**未**使用 `zap.Any`／`zap.Reflect`（測試以 `field.Type == zapcore.StringType` 鎖住）。
- 遮罩值一律委派 v3 core 函式與 `Sensitive[T].String()`，未自行重寫遮罩邏輯（測試逐一比對 `zapfield.X(...) == zap.String(key, masker.X(...))` 證明逐字一致）。
- `Sensitive` adapter 只讀 `.String()` 快取值，**不**呼叫 `Reveal()`、不碰 raw；zero-value（mask 為 nil）安全退化為 `zap.String(key, "")`。

## Files Changed

- `v3/zapfield/go.mod` — 新增獨立 module 定義（隔離 zap 依賴）
- `v3/zapfield/go.sum` — `go mod tidy` 產生
- `v3/zapfield/doc.go` — package doc，說明用途與依賴隔離設計
- `v3/zapfield/field.go` — 12 個 masker type 的 zap.Field helper
- `v3/zapfield/sensitive.go` — `Sensitive[T]` adapter
- `v3/zapfield/field_test.go` — helper 委派/逐字一致/型別/空值測試
- `v3/zapfield/sensitive_test.go` — adapter 使用 masked、zero-value 退化、泛型測試

v3 core 未做任何修改（已驗證 `v3/go.mod` 不含 zap）。

## Verification

- `go build ./...`（v3/zapfield）: Success
- `go vet ./...`: No issues found
- `go test -race -cover ./...`: 17 passed in 1 package
- `grep -i zap v3/go.mod`: 無輸出 → core 未被 zap 污染

## Notes

- 依 Out of Scope，未實作 `Mask`（動態型別）對應 helper、未做 zap Core wrapper、未發 release tag。
- `go mod tidy` 額外帶入 `go.uber.org/multierr v1.10.0`（zap 的 indirect 依賴），屬正常。
