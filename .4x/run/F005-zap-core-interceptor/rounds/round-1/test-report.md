# Test Report — Round 1

## Summary
PASS — 12/12 criteria met

`4x verify` 全綠（go build / go vet / go test -race -cover 皆 exit 0，coverage 97.1%）。
7 個 observer-based 攔截測試全數通過，code review 證據齊全。

## Results
| # | Criterion | Status | Evidence |
|---|---|---|---|
| AC-1 | WrapCore 簽名不回傳 error，回傳 zapcore.Core 可傳入 zap.New | PASS | `core.go:72` `func WrapCore(core zapcore.Core, rules InterceptRules) zapcore.Core`，無 error 回傳值；`core_test.go:15` `zap.New(WrapCore(obsCore, rules))` 編譯並寫 log 成功；`go build ./...` exit 0 |
| AC-2 | InterceptRules 含 Keywords []string 與 Patterns []*regexp.Regexp | PASS | `core.go:22-28`：`Keywords []string`、`Patterns []*regexp.Regexp` 兩 exported 欄位；測試以 struct literal 建構兩欄位（`core_test.go:31,41`） |
| AC-3 | keyword 命中時 string field 全遮罩（== masker.All） | PASS | `TestWrapCore_KeywordMatchMasksString`：`zap.String("phone","0987654321")` → `"**********"`（10 字元全遮）；`interceptFields` 走 `masker.All`（`core.go:105`），`All` 為每字元換遮罩字元（`v3/convenience.go:99`） |
| AC-4 | keyword 為 case-insensitive 子字串（strings.Contains） | PASS | `TestWrapCore_KeywordCaseInsensitive`：keyword `"phone"` 命中 `"Phone"` 與 `"userPhone"`，兩者皆 `"**********"`；`match` 用 `strings.ToLower` + `strings.Contains`（`core.go:35-37`） |
| AC-5 | regex pattern 命中時 string field 全遮罩 | PASS | `TestWrapCore_PatternMatchMasksString`：`Patterns:[(?i)secret]` 命中 key `"api_secret"`，值 `"abcdef"` → `"******"` |
| AC-6 | 非 string field 不被改動 | PASS | `TestWrapCore_NonStringFieldUntouched`：keyword `"phone"` 命中 `"phone_count"` key，但 `zap.Int` 的 Integer 維持 `3`；`interceptFields` 僅處理 `Type == zapcore.StringType`（`core.go:104`） |
| AC-7 | 未命中保留原樣；空 rules 放行不 panic | PASS | `TestWrapCore_UnmatchedStringUntouched`（`nickname` 維持 `"johnny"`）+ `TestWrapCore_EmptyRulesPassThrough`（空 rules 下 phone/token 原樣放行、無 panic）；`match` 於 Keywords/Patterns 皆空時回 false（`core.go:33-47`） |
| AC-8 | logger.With 的 context field 同樣被攔截 | PASS | `TestWrapCore_WithContextFieldIntercepted`：`logger.With(zap.String("password","hunter2")).Info("ctx")` → `"*******"`；`With` 先 `interceptFields` 再重新包裝（`core.go:78-83`） |
| AC-9 | 繁中 GoDoc 警告 false positive/negative 風險 | PASS | `core.go:16-21` 明示 false positive（`phone_count` 誤遮）、false negative（不在清單漏遮），並說明業務 code 已用顯式遮罩（Sensitive[T]/zapfield.Phone）時本層多餘；`WrapCore` GoDoc（`core.go:56-71`）同樣警告 |
| AC-10 | 僅新增 core.go 與 core_test.go，不改 v3 core module | PASS | commit `0af6e97 --stat`：僅 `v3/zapfield/core.go`、`v3/zapfield/core_test.go` 兩檔 217 行新增；無其他 module 檔案 diff |
| AC-11 | 全測試通過帶 race；既有測試不受影響 | PASS | `cd v3/zapfield && go test -race -cover ./...` exit 0，`coverage: 97.1%`；整包測試（含既有 field_test.go/sensitive_test.go）全綠 |
| AC-12 | doc/註解不含禁用 token；遮罩經 field.String | PASS | `grep -E 'zap.Any\|zap.Reflect\|AddReflected\|Reflect\('` 於 core.go/core_test.go → 0 matches；遮罩一律 `out[i].String = masker.All(...)`（`core.go:105`） |

## Verdict
PASS
