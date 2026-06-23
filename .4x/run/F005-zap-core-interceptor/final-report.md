# Final Report — F005-zap-core-interceptor: F005 zap Core interceptor: field name keyword 攔截

## Status
ready-for-review

## Summary
在 `v3/zapfield/core.go` 新增 `WrapCore` / `InterceptRules` / `maskingCore`，於 log 寫出前依 field name keyword（case-insensitive `strings.Contains`）或 regex 對頂層 string field 全遮罩，作為「最後一道防線」。12 條 AC 與三條 feature 紅線全部通過，`go test -race -cover` 全綠（coverage 97.1%），變更僅限新增 `core.go` 與 `core_test.go`。

## Open Issues
None — all issues resolved.

（Review、Test 兩份報告均判定 PASS，零 critical / 零 warning；本輪無 deep-review-report、無 escalation。）

### Feature 紅線檢核
- 只攔截 string 類型的 zap field — PASS（`interceptFields` 以 `Type == zapcore.StringType` 守門，core.go:104）
- keyword matching 用 `strings.Contains`、case-insensitive — PASS（`match` 雙向 `strings.ToLower` + `Contains`，core.go:35-37）
- 文件明確警告 keyword matching 的 false positive 風險 — PASS（`InterceptRules` / `WrapCore` GoDoc 以繁中明示 `phone_count` 誤遮、漏遮與「顯式遮罩存在時本層多餘」，core.go:16-21, 66-68）
