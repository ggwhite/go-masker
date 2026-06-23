# Test Report — Round 1

## Summary
PASS — 17/17 criteria met

verify.json `passed: true`。三項 verify 指令全數 exit 0：`gofmt -l`（無輸出）、`go vet ./...`（無錯）、`go test -race -cover ./...`（ok，coverage 65.9%，120 passed）。

## Results
| # | Criterion | Status | Evidence |
|---|---|---|---|
| AC-1 | `Sensitive[T]` 定義於 v3/sensitive.go，`raw`/`masked`/`mask` 三 field 皆 unexported | PASS | sensitive.go:23-27 三 field 首字小寫；`go vet ./...` exit 0；套件外只能透過 `Reveal()` 取值（無 exported field） |
| AC-2 | `Reveal()` 為取得原值唯一途徑，回傳原輸入 | PASS | sensitive.go:45-47 value receiver 回 `s.raw`；`TestSensitiveReveal` PASS（`NewPhone("0987654321").Reveal()=="0987654321"`）；parity 測試對九建構子皆驗 Reveal==原值 |
| AC-3 | `String()` 實作 fmt.Stringer，回傳遮罩值不含原值 | PASS | `TestSensitiveString` PASS：`fmt.Sprint(NewPhone("0987654321"))=="0987***321"` 且不含 `"654"` |
| AC-4 | `GoString()` 實作 fmt.GoStringer，`%#v` 不洩漏 raw/內部欄位 | PASS | `TestSensitiveGoString` PASS：`%#v` of NewEmail 等於 `Email(...)` 遮罩值且不含 `"chang"` |
| AC-5 | `MarshalJSON()` 實作 json.Marshaler，struct 序列化得遮罩結果 | PASS | `TestSensitiveMarshalJSON` PASS：`{"phone":"0987***321","message":"驗證碼 1234"}`（中文正確跳脫） |
| AC-6 | `MarshalText()` 實作 encoding.TextMarshaler，回傳遮罩 bytes | PASS | `TestSensitiveMarshalText` PASS：`NewID("A123456789").MarshalText()==[]byte("A12345****")`，err nil |
| AC-7 | `LogValue()` 實作 slog.LogValuer，slog 輸出遮罩值不含原值 | PASS | `TestSensitiveLogValue` PASS：slog TextHandler 輸出含 `0987***321`、不含 `0987654321` |
| AC-8 | 九建構子皆 `func(string) Sensitive[string]`，masked 與對應 v3 函式逐字一致 | PASS | `TestBuiltinConstructorsParity`（9 子測試全 PASS）：每個 `ctor(v).String()==Fn(v)` 且 `Reveal()==v` |
| AC-9 | NewPhone 綁 Mobile、NewTel 綁 Tel，格式不同不混用 | PASS | `TestPhoneVsTel` PASS：NewPhone=="0987***321"、NewTel=="(02)2799-****"；constructors.go:13,77 綁定來源不同 |
| AC-10 | `NewSensitive[T]` 支援自訂遮罩，建構時即執行並快取 | PASS | `TestNewSensitiveCustomType` PASS：對 `int` 自訂 maskFn（全 `#`）→ String()=="#####"、Reveal()==12345；sensitive.go:36-42 建構時算 masked |
| AC-11 | `Equal()` 用 reflect.DeepEqual 只比 raw，不因 masked 相同誤判 | PASS | `TestSensitiveEqual` PASS：同 raw→true、異 raw→false；兩個 Password masked 相同但 raw 不同→false |
| AC-12 | `UnmarshalJSON()` 還原原值並依 mask 重算 masked | PASS | `TestSensitiveUnmarshalJSON` PASS：預填 `NewPhone("")` 後 Unmarshal → Reveal()=="0987654321"、String()=="0987***321" |
| AC-13 | `mask==nil` 退化不洩漏原值（masked=""） | PASS | `TestSensitiveNilMaskConstructor`＋`TestSensitiveUnmarshalNilMask` PASS：`NewSensitive("secret",nil).String()==""`；裸 Sensitive UnmarshalJSON 後 String()=="" 且 Reveal()=="secret"；sensitive.go:37-40,90-92 安全預設 |
| AC-14 | 五安全輸出 interface 由 value type 滿足＋編譯期斷言 | PASS | sensitive.go:12-18 五個 `var _ ... = Sensitive[string]{}` 斷言（Stringer/GoStringer/json.Marshaler/encoding.TextMarshaler/slog.LogValuer）；編譯通過 |
| AC-15 | 提供可執行 example（ExampleNewPhone + JSON example），go test 驗證 | PASS | `ExampleNewPhone`、`ExampleNewEmail`、`ExampleSensitive_jSON` 皆有 `// Output:` 並 PASS |
| AC-16 | 全測試含 race detector 通過、合理覆蓋率 | PASS | `go test -race -cover ./...` ok，coverage 65.9%，120 passed；核心型別＋九建構子＋五 interface＋Equal＋UnmarshalJSON 皆有測試 |
| AC-17 | exported 識別字皆有繁中 GoDoc（首句以名稱開頭）；gofmt/go vet 通過 | PASS | `gofmt -l` 無輸出；`go vet ./...` exit 0；sensitive.go/constructors.go 各 exported 識別字皆有繁中 GoDoc 第一句以名稱開頭 |

## Verdict
PASS
