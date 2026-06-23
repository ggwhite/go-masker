# Review Report — Round 1

## Summary
PASS

v3 core 基礎建設完整實作：新 `Masker` interface（`Mask(value string) string`）、functional option `WithMaskChar`、精簡 `MaskerType` 常數、package-level 便利函式、`MustMarshal`。全部 12 種 masker 已遷移，遮罩演算法與 v2 逐字一致（僅將 `Marshal(s, i string)` 的 mask char 參數改為建構時注入的內部欄位）。build / vet / test 全綠（95 測試）。零外部依賴。

## Checklist
| Item | Status | Notes |
|------|--------|-------|
| R1: 既有 masker 遮罩邏輯不變，只改 interface 簽名 | PASS | mobile/email/password/name/address/id/credit/tel/url/none/all/abuse 逐一與 v2 比對，演算法逐字一致；唯一差異為 `s string` 參數改為注入欄位 `m.mask`。`parseGenericMask`(first-N/last-N) 亦完全一致 |
| R2: core package 零外部依賴 | PASS | `v3/go.mod` 無 require 區塊，無 go.sum；僅用 stdlib（fmt/reflect/strings/sync/net/url/strconv/unicode） |
| R3: 所有 exported 函式有 GoDoc（繁中、名稱開頭） | PASS | heuristic 掃描全部非 test .go 檔，無任何 exported func/type 缺 GoDoc；逐句確認皆以識別字名稱開頭、內容繁中、多附 Example |
| R4: 每個 masker 保留獨立檔案 | PASS | mobile/email/password/name/address/id/credit/telephone/url/none/generic/abuse/abuse_loader 各自獨立檔案 |
| R5: 新 interface Mask(value string) string | PASS | `v3/masker.go:36` 定義；`grep "Marshal(s" v3/*.go` 無殘留舊簽名 |
| R6: WithMaskChar 不污染 default、URL/None 例外 | PASS | `NewMaskerMarshaler` 每次建構各自獨立 masker 實例；`TypeURL` 用無 mask char 的 `URLMasker`、`TypeNone` 用 `NoneMasker`。`TestWithMaskChar` / `TestWithMaskChar_URLAndNoneException` 已覆蓋隔離與例外 |
| R7: convenience 函式 + 動態 Mask | PASS | `convenience.go` 提供 Mobile/Email/Password/Name/Address/ID/Credit/Tel/URL/Abuse/None/All；動態 `Mask(t, value)` 未知型別回原值，與契約一致 |
| R8: MustMarshal panic、Marshal 保留 | PASS | `Marshal` 保留回 error；`MustMarshal` 未知型別 panic；`Struct()` 維持可用，內部改呼叫 `masker.Mask(value)` |
| R9: 編譯/測試/vet 通過 | PASS | `go build`/`go vet` 乾淨；`go test ./...` 95 測試全綠 |

## Issues

無 critical / warning / info 問題。

備註（非阻擋）：coder-report 提及 coverage 63.9%，未覆蓋的主要是 `Struct()` 的 slice/map/ptr reflect 容器分支，該重構屬獨立 feature（out of scope），核心 interface、各 masker、convenience、Option 路徑皆已覆蓋。本 feature 不要求覆蓋率門檻，不構成問題。

## Verdict
PASS
