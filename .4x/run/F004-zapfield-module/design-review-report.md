# Design Review Report — Round 0

## Summary

PASS

設計目標明確、範圍切割乾淨，且所有對 v3 core 的 API 假設皆經實際驗證屬實。zapfield 為薄轉接層（thin delegating wrapper），無多餘抽象，依賴隔離與「不洩漏原值」兩條紅線都有對應 AC 與 verify command 守住。可進入 coding。

## Architecture Risks

- **依賴隔離（核心紅線）成立**：`replace github.com/ggwhite/go-masker/v3 => ../` 由 `v3/zapfield/` 指回 `v3/` core，路徑正確；zap 依賴只寫進 `v3/zapfield/go.mod`，不會回寫 core `v3/go.mod`。AC-2 用 `! grep -q zap v3/go.mod` 驗證，且 core module path `.../v3` 不含 "zap" 子字串，檢查有效。
- **巢狀 module 不污染 core 測試**：`v3/zapfield/` 有獨立 `go.mod`，故 `cd v3 && go test ./...`（AC-11）不會 descend 進 zapfield，core 維持綠燈、未被修改的前提成立。
- **輸出逐字一致無重寫風險**：所有 helper 委派 core 既有 package-level 函式（`masker.Mobile`…）與 `Sensitive[T].String()`，不自行重寫遮罩邏輯，避免兩套邏輯 drift。已驗證 12 個函式簽名皆為 `func(string) string`，與 `func(key, value string) zap.Field` 包裝完全相容。
- **offline 可建**：`go.uber.org/zap v1.28.0` 與其傳遞依賴（multierr、atomic）已在 module cache，coder 跑 `go mod tidy` / `go build` 無需連網。

## Overengineering

- 無。設計刻意維持最薄轉接：每個 helper 一行 `zap.String(key, masker.Xxx(value))`，Sensitive adapter 一行 `zap.String(key, s.String())`。
- 範圍紀律良好：Core wrapper／keyword 攔截（F005）、`Mask` 動態型別、slog/logr adapter、獨立 release tag 皆明確列入 Out of Scope，未提前實作。
- 未引入非必要的新介面或可擴充點。

## Missing Requirements

- AC 覆蓋完整：go.mod 結構（AC-1）、依賴隔離（AC-2）、可編譯（AC-3）、12 個 helper 齊全（AC-4）、值逐字一致＋key 正確（AC-5）、強制 StringType／禁 reflection（AC-6）、Sensitive adapter（AC-7）、不洩漏原值含 zero-value 與非 string T（AC-8）、race 測試（AC-9）、繁中 GoDoc（AC-10）、core 未動（AC-11）。兩條 feature rule 各有對應 AC（AC-2、AC-6）。
- test-strategy 的 verify_commands 與 AC 對齊，且 grep / build / vet / race 皆可機檢，無模糊真相源。
- 非阻擋性小瑕疵（coder 可順手修正，不需退回重設計）：
  - AC-8 與 task-brief 1.16 描述 zero-value 時用「mask 為 nil」措辭，但 core 實作的是 `masked string` 欄位，zero-value 為 `""`（非 nil）。行為斷言（`f.String == ""`）正確，僅文字描述需精準化。
  - verify_commands 未包含 `go mod tidy`，但 build 需 `go.sum`；task-brief 1.1 已要求 coder 執行 tidy，且依賴已在 cache，無實質風險，僅提醒 coder 勿遺漏。

## Verdict

PASS
