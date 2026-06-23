# Coder Report — Round 1

## What Was Done

在 v3 module 新增 `Sensitive[T]` 泛型安全型別，把遮罩從「開發者記得做」變成「洩漏原值要刻意呼叫 `Reveal()`」。所有自動輸出路徑（`fmt`、`%#v`、`json.Marshal`、`encoding.TextMarshaler`、`slog`）一律輸出遮罩值，原值只能透過唯一出口 `Reveal()` 取得。底層遮罩沿用 F001 既有 v3 package-level 函式（`Mobile/Email/...`），未自行重寫遮罩邏輯，masked 輸出逐字一致。

涵蓋 task-brief 全部 8 項任務：

1. `Sensitive[T]` 核心 struct（`raw`/`masked`/`mask` 全 unexported）+ value receiver `Reveal()`。
2. 通用建構子 `NewSensitive[T]`，建構時即執行 `maskFn(raw)` 算好 masked 並快取、保存 maskFn。
3. 五個安全輸出 interface 實作（皆 value receiver，只回 `s.masked`）：`String`、`GoString`、`MarshalJSON`（走 `json.Marshal` 正確跳脫）、`MarshalText`、`LogValue`；含五個編譯期介面斷言。
4. 九個內建建構子委派 `NewSensitive` 並綁定對應 v3 函式（`NewPhone→Mobile`、`NewTel→Tel` 不混用），各有繁中 GoDoc + 範例。
5. `Equal` 用 `reflect.DeepEqual` 只比 `raw`（支援非 comparable 的 T，不因 masked 相同誤判）。
6. `UnmarshalJSON`（pointer receiver）解碼進 `raw` 後用既有 `mask` 重算 masked。
7. `mask == nil` 退化：`NewSensitive(raw, nil)` 與未綁定 mask 的 `UnmarshalJSON` 皆使 masked 設為 `""`，不洩漏原值。
8. 測試 + 可執行 example（`ExampleNewPhone`/`ExampleNewEmail`/`ExampleSensitive_jSON`）。

對應 AC-1 ~ AC-17 全數覆蓋。

## Files Changed

- `v3/sensitive.go` — 新增：`Sensitive[T]` struct、`Reveal`、`NewSensitive`、五個安全輸出 interface、`Equal`、`UnmarshalJSON`、編譯期介面斷言。
- `v3/sensitive_constructors.go` — 新增：九個內建建構子（`NewPhone/NewEmail/NewPassword/NewID/NewCredit/NewName/NewAddress/NewTel/NewURL`）。
- `v3/sensitive_test.go` — 新增：核心型別測試（Reveal/String/GoString/MarshalJSON/MarshalText/LogValue/Equal/UnmarshalJSON/nil-mask 退化/自訂型別 `int`）+ JSON example。
- `v3/sensitive_constructors_test.go` — 新增：九建構子 parity 表格測試、Phone vs Tel、`ExampleNewPhone`/`ExampleNewEmail`。

## Verification

- `cd v3 && gofmt -l sensitive.go sensitive_constructors.go sensitive_test.go sensitive_constructors_test.go` → 無輸出（已格式化）
- `cd v3 && go vet ./...` → No issues found
- `cd v3 && go test -race -cover ./...` → 全綠，120 passed in 1 package（含 example output 驗證）

## Notes

- 編譯期斷言改用標準 `encoding` 套件（`var _ encoding.TextMarshaler = Sensitive[string]{}`），符合 task-brief 原意，較自訂 alias 乾淨。
- nil maskFn 測試以明確型別轉換 `(func(string) string)(nil)` 傳入，避免 untyped nil 推不出泛型 T。
- 僅在 `v3/` 內新增檔案，未動 v2 與 F001 既有 v3 檔。
