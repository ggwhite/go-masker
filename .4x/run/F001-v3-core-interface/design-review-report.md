# Design Review Report — Round 0

## Summary

PASS

（本輪覆寫前一份 stale 的 FAIL 報告：前一份所列 4 點疑慮均已被 Designer 在當前 task-brief / acceptance-criteria 修正——見下方「Missing Requirements / 前輪疑慮追蹤」。）

設計完整、可實作，且把唯一紅線（遮罩輸出與 v2 逐字一致）標示明確。task-brief、acceptance-criteria、test-strategy 三者一致，15 條 AC 皆具可驗證方法。僅剩兩處次要的「待釘住」項目（未知 type 行為的 fail-safe 預設、abuse mask char 交叉參照文字），但 AC 機制已強制以測試保證一致性，coder 不會被卡，不阻擋開工。

## Architecture Risks

- **實例化策略正確、無單例污染風險**：task-brief Task 2（行 34）明確要求 `DefaultMaskerMarshaler` 與 `NewMaskerMarshaler(WithMaskChar('#'))` 各自持有獨立 masker 集合，AC-4 以「同程序內 default 仍輸出 `*`」把關。對照 v2 `masker.go` 之 `masker string` 欄位設計，v3 改為「建構時把 mask char 注入各 masker 實例內部欄位」是合理且低風險的轉換。✅
- **mask char 型別轉換**：design 採 `WithMaskChar(c rune)`，v2 內部 `strLoop` 走 `strings.Repeat(string, n)`、`masker` 欄位為 `string`。coder 需把 `rune` 轉 `string` 注入，屬 1 行轉換，多 byte rune 也由 `strings.Repeat` 正確處理，無風險。
- **generic first-N/last-N 的 mask char 來源**：`parseGenericMask`（generic.go）非註冊型 Masker，由 `MaskerMarshaler.Marshal` dispatch 時傳入 mask char。v3 仍由 marshaler 持有並於 dispatch 傳入即可，與 design 一致，無架構衝突。
- **Struct() 過渡相容**：AC-12 + task-brief 要求 `Struct()` 維持可用、內部 `Marshal` 呼叫配合新 interface 調整即可，正確地把 reflect cache 推遲到獨立 feature。對照 masker.go:238/298 的 `m.Marshal(...)` 呼叫點，遷移面收斂可控。✅

## Overengineering

- **無過度設計**。scope 紀律良好——`Sensitive[T]`、Struct reflect cache、`Format()`、zapfield/slog 都正確排在 out-of-scope，與 design 文件優先級表一致，未被提前拉進來。
- `Option func(*MaskerMarshaler)` functional option 為 Go 慣用模式，非為「未來需求」預先抽象。
- 維持「每個 masker 一個檔案」符合既有專案結構與紅線，未引入不必要的 registry／plugin 抽象。✅

## Missing Requirements

### 前輪疑慮追蹤（均已修正）

- ✅ **【前輪 Blocking】`WithMaskChar` 對 URL masker 契約矛盾**：當前 task-brief 行 29–32 已新增「`WithMaskChar` 例外」段，明示 `URLMasker` 使用 `url.Redacted()` 固定輸出 `xxxxx`、不受 `WithMaskChar` 影響、coder 不得改寫 URL 遮罩邏輯；AC-4b 新增對應斷言（`WithMaskChar('#')` 下 URL 仍為 `http://u:xxxxx@host`）。矛盾已消除。
- ✅ **【前輪 Should fix】Abuse 預設空字典行為**：AC-9b 已明確定義「default 空字典回原值；載入詞典後可正確遮罩」。
- ✅ **【前輪 Minor】AC-8 grep 誤報**：AC-8 已改為精確的 `grep -rnE "MaskerType(None|Password|...|MapStruct)"`，不再誤觸。
- ✅ **【前輪 確認項】instance-per-marshaler 不污染 default**：task-brief Task 2 行 34 已明文化。

### 本輪殘留（次要、不阻擋）

- **（次要，建議釘住 fail-safe）AC-10 未知 type 的 `Mask(t MaskerType, value string) string` 行為未拍板**：AC-10（行 16）與 task-brief Task 4（行 47）都寫「回原值或空字串，二擇一並一致」，但未指定哪個。對「遮罩／資安」函式庫而言，未知 type 時**回傳原值等同洩漏未遮罩明文**，屬安全相關預設值。建議採 **fail-safe：回傳空字串 `""`**——理由：(1) 與 v2 `Marshal` not-found 時回 `("", err)` 的空字串一致；(2) 寧可漏顯示也不漏明文。請 coder 採此選項、tester 於 AC-10 以 `""` 斷言。註：type 多為編譯期常數（`TypeMobile`），未知 type 僅在手寫字串時發生，風險面有限，故不阻擋開工。
- **（次要，文件文字缺漏）task-brief Task 2 對 AbuseMasker 的交叉參照不完整**：行 32 寫「AbuseMasker … 遮罩字元行為見下方 Task 4 的 Abuse 說明」，但 Task 4（行 49）只說明「空字典回原值」，**未說明命中詞時用哪個 mask char**。經查 abuse.go:140/165，`AbuseMasker.Marshal` 以 `maskChar` 遮罩命中詞（`strLoop(maskChar, j-i)`）。故正確行為為：AbuseMasker 比照他 masker **受 `WithMaskChar` 影響**（載入詞典後用設定字元遮罩）。AC-9b 只測 default `*`、未測 `WithMaskChar` 下 abuse，故不阻擋。建議 coder 讓 AbuseMasker 一致注入 mask char。

### 已抽查確認（無缺口）

- **AC-6 輸出對照**：對照 mobile.go（`overlay(i, strLoop(s,3), 4, 7)` → `0987***321`）、none.go、generic.go、masker.go GoDoc 範例，AC-6 期望值與 v2 實作逐字一致，紅線可被測試打臉，無誤植。
- **CI（AC-15）/ race（AC-13）/ GoDoc 繁中（AC-14）** 皆有對應 verify_commands 或驗證手段，無缺口。

## Verdict

PASS
