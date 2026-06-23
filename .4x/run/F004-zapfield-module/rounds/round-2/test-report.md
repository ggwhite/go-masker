# Test Report — Round 2

## Summary
PASS — 11/11 criteria met

## Results
| # | Criterion | Status | Evidence |
|---|---|---|---|
| AC-1 | `v3/zapfield/go.mod` 內容正確（module path、go 1.21、require zap + core、replace） | PASS | `cat v3/zapfield/go.mod`：`module github.com/ggwhite/go-masker/v3/zapfield`、`go 1.21`、`require github.com/ggwhite/go-masker/v3 v3.0.0` 與 `go.uber.org/zap v1.28.0`、`replace github.com/ggwhite/go-masker/v3 => ../` 全部齊備。 |
| AC-2 | core `v3/go.mod` 不含 zap 依賴 | PASS | `grep -i zap v3/go.mod` exit=1（無 match）；verify group `test -z "$(grep -i zap v3/go.mod)"` exit 0 → `AC-2 core clean of zap`。 |
| AC-3 | zapfield 可編譯 | PASS | verify `cd v3/zapfield && go build ./...` exit 0（172ms）。 |
| AC-4 | 12 個 `func(key, value string) zap.Field` helper | PASS | `grep -cE "func [A-Z].*\(key, value string\) zap\.Field" field.go` = 12；helper：Phone、Email、Password、Name、Address、ID、Credit、Tel、URL、Abuse、None、All。 |
| AC-5 | 每個 Field helper 的 `Key`==key 且字串值與 core 函式逐字一致 | PASS | `TestFieldHelpers_DelegateToCore` table-driven 逐一斷言 `want := zap.String(key, tt.core(tt.value))` 並比對 `got.Key`/`got.String`，對照 `masker.Mobile/Email/Password/Name/Address/ID/Credit/Tel/URL`；test 全 PASS（17 passed）。 |
| AC-6 | 全以 `zap.String` 建構，`f.Type == zapcore.StringType`，無 `zap.Any`/`zap.Reflect` | PASS | `grep -REn 'zap.Any|zap.Reflect' v3/zapfield/*.go` exit=1（無命中，round-1 誤判已解除）；`zap.String` 出現 13 次；test 斷言 `got.Type != zapcore.StringType` 為失敗條件。 |
| AC-7 | `Sensitive[T]` 回傳 `zap.String(key, s.String())`，值==`s.String()`，型別 StringType | PASS | `TestSensitive_UsesMaskedValue`：`want := zap.String("phone", s.String())`，斷言 `got.Type == zapcore.StringType` 且 key 正確；test PASS。 |
| AC-8 | `Sensitive` 不洩漏原值：非 string T 與 zero-value 皆只輸出 masked，zero-value 輸出空字串 | PASS | `TestSensitive_ZeroValueSafeDegradation` 斷言 zero-value `got.String == ""`；`TestSensitive_GenericType`（int）斷言輸出 `"*****"` == `s.String()`；`TestSensitive_UsesMaskedValue` 另斷言 `got.String != s.Reveal()` 防洩漏。 |
| AC-9 | 全部 zapfield 測試通過（含 race） | PASS | verify `cd v3/zapfield && go test -race ./...` exit 0；`ok github.com/ggwhite/go-masker/v3/zapfield`，17 passed。 |
| AC-10 | exported 識別字皆有繁中 GoDoc，第一句以名稱開頭 | PASS | `go vet ./...` exit 0；GoDoc 抽查：`// Package zapfield 提供...`、`// Sensitive 把...`、`// Phone 回傳...`、`// Email 回傳...` 等 13 條 helper/adapter 註解皆繁中且以識別字名稱開頭。 |
| AC-11 | core 未被修改、core 測試全綠 | PASS | `git status --porcelain v3`（排除 zapfield）無輸出 → core 檔無改動；verify `cd v3 && go test ./...` exit 0（`ok github.com/ggwhite/go-masker/v3`）。 |

## Rules Check
- zapfield 的 go.mod 獨立於 core，避免 zap 依賴污染 core → ✅（AC-1 + AC-2）
- Field helper 回傳值一律是 zap.Field 且以 zap.String 建構（非 zap.Any，避免 reflection）→ ✅（AC-4 + AC-6）

## Verdict
PASS
