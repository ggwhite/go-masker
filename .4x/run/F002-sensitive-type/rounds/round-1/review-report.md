# Review Report — Round 1

## Summary

PASS

實作完整覆蓋 task-brief 全部 8 項任務與 AC-1～AC-17，四條 feature 紅線全數守住。`gofmt -l` 無輸出、`go vet` 無 issue、`go test -race -cover` 全綠（120 passed）。零 critical、零 warning。

## Checklist

| Item | Status | Notes |
|------|--------|-------|
| AC-1 三 field unexported | PASS | `sensitive.go:24-26` `raw`/`masked`/`mask` 皆小寫；套件外只能經 `Reveal()` 取值 |
| AC-2 Reveal 唯一原值出口 | PASS | value receiver 回傳 `s.raw`；`TestSensitiveReveal` + parity 測試各建構子皆驗 |
| AC-3 String 遮罩不洩漏 | PASS | 回 `s.masked`；`TestSensitiveString` 斷言 `0987***321` 且不含 `654` |
| AC-4 GoString 防 %#v | PASS | `TestSensitiveGoString` 斷言不含原帳號 `chang` |
| AC-5 MarshalJSON 正確跳脫 | PASS | `json.Marshal(s.masked)`；struct 序列化測試含中文訊息逐字一致 |
| AC-6 MarshalText | PASS | 回 `[]byte(s.masked)`、err nil；測試驗 `A12345****` |
| AC-7 LogValue | PASS | `slog.StringValue(s.masked)`；測試斷言 buf 含遮罩值、不含原值 |
| AC-8 九建構子 parity | PASS | `TestBuiltinConstructorsParity` 逐一比對 v3 package-level 函式輸出 |
| AC-9 Phone vs Tel 不混用 | PASS | `NewPhone→Mobile`、`NewTel→Tel`；`TestPhoneVsTel` 驗格式不同 |
| AC-10 NewSensitive 自訂型別 | PASS | `TestNewSensitiveCustomType` 以 `int` + 自訂 maskFn 驗 |
| AC-11 Equal 只比 raw | PASS | `reflect.DeepEqual(raw)`；測試含「masked 同但 raw 異」回 false |
| AC-12 UnmarshalJSON round-trip | PASS | pointer receiver，解碼後用既有 mask 重算；測試驗 Reveal + String |
| AC-13 nil mask 退化不洩漏 | PASS | 建構子與 UnmarshalJSON 皆設 masked `""`；兩測試覆蓋 |
| AC-14 編譯期斷言 | PASS | `sensitive.go:12-18` 五個 value-type 斷言齊全 |
| AC-15 可執行 example | PASS | `ExampleNewPhone`/`ExampleNewEmail`/`ExampleSensitive_jSON` 均有 Output 並通過 |
| AC-16 race + cover | PASS | `go test -race -cover ./...` 全綠 120 passed |
| AC-17 GoDoc + gofmt + vet | PASS | exported 皆有繁中 GoDoc（首句以名稱開頭）；gofmt clean、vet clean |
| 規則：raw field unexported | PASS | 同 AC-1 |
| 規則：自動輸出路徑不回傳原值 | PASS | String/GoString/MarshalJSON/MarshalText/LogValue 全部只讀 `s.masked`，不碰 `s.raw` |
| 規則：只有 Reveal() 取原值 | PASS | 僅 `Reveal()` 回傳 `s.raw`；其餘方法不暴露 |
| 規則：建構時算好並快取，後續零成本 | PASS | `NewSensitive` 建構時執行 `maskFn(raw)` 存入 `masked`；取值僅讀欄位 |

## Issues

無 critical / warning / info 等級問題。

補充觀察（不影響判定）：
- `UnmarshalJSON` 為 pointer receiver、`MarshalJSON` 為 value receiver，當 `Sensitive[T]` 作為 struct 欄位時 json 取址呼叫 pointer 方法正常運作，`TestSensitiveUnmarshalJSON` 已驗證 round-trip 成功。
- Scope 嚴格遵守：僅新增 `v3/sensitive.go`、`v3/sensitive_constructors.go` 及兩個測試檔，未動 v2 與 F001 既有 v3 檔；遮罩邏輯全程委派 v3 package-level 函式，無自行重寫。

## Verdict

PASS
