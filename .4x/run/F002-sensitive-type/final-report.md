# Final Report — F002-sensitive-type: F002 Sensitive[T] 泛型安全型別

## Status
ready-for-review

## Summary
在 v3 module 新增 `Sensitive[T]` 泛型安全型別（`v3/sensitive.go`、`v3/sensitive_constructors.go`），把遮罩從「開發者記得做」反轉為「洩漏原值要刻意呼叫 `Reveal()`」。所有自動輸出路徑（`String`/`GoString`/`MarshalJSON`/`MarshalText`/`LogValue`）只回傳建構時快取的 masked，原值唯一出口為 `Reveal()`。AC-1～AC-17 全數 PASS，四條 feature 紅線全數守住。

## Open Issues
None — all issues resolved.

四項 feature 規則驗證結果：
- raw field 必須 unexported — PASS（`raw`/`masked`/`mask` 三 field 皆小寫，套件外只能經 `Reveal()` 取值）
- 任何 String()/GoString()/MarshalJSON()/MarshalText()/LogValue() 不得回傳原始值 — PASS（五個方法皆 value receiver 只讀 `s.masked`，不碰 `s.raw`；含 `mask==nil` 退化情況 masked 為空字串）
- 只有 Reveal() 能取得原始值 — PASS（僅 `Reveal()` 回傳 `s.raw`）
- 建構子執行遮罩運算並快取，後續取值零額外成本 — PASS（`NewSensitive` 建構時執行 `maskFn(raw)` 存入 `masked`，取值僅讀欄位）

驗證證據：`gofmt -l v3/` 無輸出、`go vet ./...` exit 0、`go test -race -cover ./...` 全綠（120 passed，coverage 65.9%）。Scope 嚴格遵守：僅新增四個檔案，未動 v2 與 F001 既有 v3 檔，遮罩邏輯全程委派 v3 package-level 函式。無 deep-review-report、無 escalation。
