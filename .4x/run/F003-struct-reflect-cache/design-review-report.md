# Design Review Report — Round 0

## Summary

PASS

設計紮實、向後相容意圖明確、cache 架構（`sync.Map` + `LoadOrStore`，key 為 `reflect.Type`）正確且 concurrency-safe，scope 與 out-of-scope 界定清楚，12 條 AC 皆可驗證。以下記錄的是 Coder 需留意的風險與須補強的測試覆蓋，皆非阻擋性問題，故判定 PASS。

## Architecture Risks

- **Cache 設計正確**：`reflect.Type` 作 key、`*typeInfo` 作 value、`LoadOrStore` 處理 cache miss 競態，是 Go 標準做法。不同 type 必有不同 `reflect.Type`，無 key 碰撞風險。embedded／同名 type 也安全。
- **巢狀遞迴受惠 cache**：`mask:"struct"` / `mapstruct` / slice-of-struct 仍走 `m.Struct()` / `maskMapStructValue()`，各層各自查 cache，設計一致，無漏網。
- **AC-3 觀測計數的 race 風險（實作注意）**：AC-3 要求「metadata 只建立一次」需「同 package 可觀測的內部計數」。若 `buildTypeInfo` 用普通 `int++` 計數，會與 AC-4 的並行 `-race` 測試衝突（data race 告警）。**Coder 必須用 `sync/atomic`（如 `atomic.AddInt64`）或將計數限定在單執行緒 AC-3 測試中**。此為實作層級提醒，非設計缺陷。
- **語意保真是最大實作風險**：`masker.go:224-362` 的所有分支（string / struct / ptr / map / 5 種 slice element kind / interface、nil pointer/slice/map 略過、未匯出欄位 skip、無 tag 直接複製、有 tag 但 kind 落 `default` 仍複製）必須由 `fieldMask` metadata 完整表達且輸出 byte-identical。brief 已逐項列出，AC-1 以既有測試把關，方向正確。

## Overengineering

- 無過度設計。`fieldMask` 攜帶的分類資訊（欄位 index、`MaskerType`、kind 分類、slice element kind、遞迴旗標）都是 `Struct()` 既有分支實際需要的，非臆測性擴充。
- `Format()` 明確排除可插拔格式選項／多輸出格式（out-of-scope #6），只做單一確定性 `%v` 風格，克制得宜。
- `resetTypeCache()` 限同 package 內部、不匯出，不污染 public API，符合「不新增 package-level wrapper」的既有 API 一致原則。

## Missing Requirements

- **`Format()` 對非 string／非 ptr-to-struct 欄位的渲染未完整規格化（須補強）**：Task #3 與 AC-7/AC-8 精確定義了 string 欄位 parity、ptr-to-struct 以 `&{...}` 遞迴展開、nil pointer 顯示 `<nil>`，但**未明定 slice（of string / struct / ptr）、map / mapstruct、nested struct value、interface 欄位在 Format 輸出中如何呈現**。鑑於 `Format()` 的用途是「輸出可安全記 log 的遮罩字串」，若這些欄位未一併遮罩會造成**敏感資料外洩**。設計的核心原則（「共用同一條遮罩路徑」「確定性 %v」）已足以推導正確行為——Coder 應對所有欄位 kind 套用與 `Struct()` 相同的遮罩語意、再以確定性 %v 渲染，**不得對任何欄位直印原值**。建議在測試策略補一個含「遮罩 slice + mapstruct + 巢狀 struct」的 Format fixture，明確擋住外洩。此為測試覆蓋與行為釐清，brief 內部無矛盾，不需 Designer 重寫，故不阻擋。
- **AC-10「≥50% CPU 降低」是經驗性目標，須以實測佐證**：`BenchmarkStruct`(warm) vs `BenchmarkStructCold` 兩者唯一差異是 `buildTypeInfo`（tag 解析 + kind 分類 + 欄位迭代）。要達 ≥50%，`buildTypeInfo` 須占 `Struct()` 總成本 ≥50%。Go 的 `StructTag.Get` 每次重新 parse 整段 tag 字串、確實昂貴，對欄位數中等的 struct 達標可行；但實際遮罩運算（每欄 `m.Marshal` 含 `RLock` 與 regexp）也很重，可能稀釋比例。**Coder 應以實測數據驗證；若 benchmark fixture 欄位過少導致未達 50%，需用更具代表性的 struct（多欄位、含巢狀）而非調參數造數據。** 測試策略已含 `-benchmem`、cold/warm 雙 benchmark，工具齊備。
- **測試策略覆蓋充分**：AC-1（既有測試不改斷言）、AC-4（≥50 goroutine `-race`）、AC-6（nil/`(*T)(nil)`/空 struct/含 nil ptr 欄位、`recover` 斷言無 panic）、AC-9（`AllocsPerRun` 比較）皆已對應到 `verify_commands`，無缺口。

## Verdict

PASS
