# Coder Report — Round 1

## What Was Done

在 `v3/zapfield` sub-module 新增一層 zap.Core wrapper，於 log 寫出前依 field name keyword / regex 自動全遮罩 string field，定位為「改不動業務 code 時的被動最後防線」。

實作對應 task-brief 四項任務：

1. **InterceptRules 結構** — `Keywords []string`（case-insensitive `strings.Contains`）+ `Patterns []*regexp.Regexp`（已編譯，故 WrapCore 不回傳 error）。未匯出 `match(key) bool`：任一 keyword 或 pattern 命中即 true；兩者皆空/nil 時一律 false，退化為「不攔截、原樣放行」，最安全且不 panic。
2. **WrapCore 實作** — `maskingCore` 內嵌 `zapcore.Core` 並持有 `rules`：
   - `Enabled` / `Sync` 由內嵌 core 直接委派（embedding 自動提供）。
   - `With` 先 `interceptFields` 再委派內層 `With`，並用同樣 rules 重新包裝回傳 core（context field 也會被攔截）。
   - `Check` 依 `Enabled` 慣例 `ce.AddCore(ent, c)`。
   - `Write` 先 `interceptFields` 再委派內層 `Write`。
   - `interceptFields` 建立 copy 後，只對 `Type == zapcore.StringType` 且 `rules.match(Key)` 的 field 設 `field.String = masker.All(field.String)`；非 string field 原樣保留，不就地改寫呼叫端 slice。
3. **doc 註解警告** — `InterceptRules` / `WrapCore` GoDoc 以繁中明示 false positive（`phone_count`）/ false negative 風險，並指出業務 code 已用顯式遮罩時本層多餘。遮罩一律走 `field.String` 寫值，未使用會觸發 reflection 的欄位建構子，doc 亦以語意描述避免逐字寫出 CI grep 的禁用 token。
4. **攔截測試** — 用 `zaptest/observer` 建立可觀測 core 後 WrapCore，涵蓋 keyword 命中、regex 命中、非 string field 不動、未命中原樣、空 rules 放行、case-insensitive（`Phone` / `userPhone`）、`logger.With(...)` context field 被攔截。

未改動 `v3/zapfield/go.mod`（zap/multierr 已是現有依賴，observer 屬 zap module）；未動 v3 core module、`field.go` 與其測試。

## Files Changed

- `v3/zapfield/core.go` — 新增：`InterceptRules`、`match`、`maskingCore`、`WrapCore` 及 `With`/`Check`/`Write`/`interceptFields`
- `v3/zapfield/core_test.go` — 新增：7 個 observer-based 攔截測試

## Verification

- `gofmt -l core.go core_test.go`：無輸出（格式正確）
- `go build ./...`：Success
- `go test ./...`：24 passed in 1 package
- `go vet ./...`：No issues found
- 禁用 token 掃描（`grep zap.Any|zap.Reflect|AddReflected` 於 core.go/core_test.go）：0 matches
- commit：`0af6e97 feat(F005-zap-core-interceptor): add WrapCore field-name keyword/regex masking layer`
