# Design Review Report — Round 0

## Summary
PASS

設計完整且可直接進入 coding。task-brief、acceptance-criteria、test-strategy 三者一致，覆蓋全部 subtask（sensitive-struct / safe-output / constructors / unmarshal / equality），並逐條對應到 == Rules to Check ==。已實地驗證設計關鍵前提：F001 已提供 `Mobile/Email/Password/Name/Address/ID/Credit/Tel/URL` 九個 package-level 函式（`v3/convenience.go`，皆 `func(string) string` 走 `DefaultMaskerMarshaler`），以及 `TypeMobile/TypeTel/TypeURL` 等型別常數（`v3/masker.go`），複用基礎成立。

## Architecture Risks
- **複用既有遮罩邏輯（已驗證，無風險）**：建構子委派 v3 既有函式而非自行重寫，AC-8 以「與對應 v3 函式逐字一致」為斷言，避免行為分歧。方向正確。
- **method set 正確性（低風險，已被斷言涵蓋）**：五個安全輸出 interface 採 value receiver，`UnmarshalJSON` 採 pointer receiver。value receiver 確保 `Sensitive[string]{}` 與 struct 非指標欄位都能滿足 Marshaler/Stringer；pointer receiver 的 `UnmarshalJSON` 需欄位 addressable，AC-12 透過 `json.Unmarshal(&req)` 流程保證可寫回。設計已用編譯期斷言（AC-14）鎖死 value type 滿足介面，風險可控。
- **`MarshalJSON` 與 `MarshalText` 並存（無衝突）**：`encoding/json` 優先採用 `MarshalJSON`，`MarshalText` 僅在無 `MarshalJSON` 時生效。兩者輸出皆為同一 masked 字串，無語意分歧，僅屬刻意冗餘以同時滿足 json 與 text 兩種編碼器。
- **`UnmarshalJSON` 依賴 receiver 既有 `mask`（已正視）**：這是本設計最關鍵的耦合點——round-trip 正確性建立在「欄位先經建構子預填以綁定 maskFn」之上。task-brief Task 6 明確標為「設計前提」，AC-12（已綁定）與 AC-13（未綁定 → masked 為 `""`）兩條同時涵蓋正常與退化路徑，沒有遺漏。
- **`raw` 防 reflect 讀取的限度**：unexported field 阻擋一般 reflect，但 `unsafe` 或同 package 仍可讀。設計文件已誠實標註「除非用 unsafe」，未對此做過度承諾，符合預期。

## Overengineering
- 無過度設計。九個固定建構子 + 一個通用 `NewSensitive`，範圍精準；Out of Scope 明確排除 struct tag 整合、`Abuse/None/All/Mask` 動態型別建構子、並行可變狀態，避免擴張。
- `Equal` 採 `reflect.DeepEqual(raw)` 而非加 `T comparable` 約束，是合理取捨——換取支援 slice/struct 等型別，且只比 `raw` 不比 `masked`（避免 masked 碰撞誤判），AC-11 已專門構造「不同 raw 但 masked 相同」案例驗證。判斷正確，非過度設計。

## Missing Requirements
- 無阻斷性缺口。AC 對五個 interface、九建構子、`Equal`、`UnmarshalJSON`、nil edge case（AC-13）、example（AC-15）、race+cover（AC-16）、gofmt/vet/GoDoc（AC-17）皆有對應驗證方法。
- test-strategy 的 gofmt 指令僅列 4 個指定檔（`sensitive.go`/`sensitive_constructors.go` 及兩測試檔），與 Scope「僅新增這些檔案、僅在 v3/ 作業」一致，無遺漏；`go vet ./...` 與 `go test -race -cover ./...` 覆蓋全套件。
- 建議（非阻斷）：coding 時注意 AC-9 的 `NewPhone→Mobile`（手機）與 `NewTel→Tel`（市話）對應，task-brief 已特別標註不可混用；AC-12 的範例值 `NewPhone("")` 預填再 Unmarshal 為合法（空字串建構不洩漏），實作須確保空值不 panic。

## Verdict
PASS
