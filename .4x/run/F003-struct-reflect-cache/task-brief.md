# Task Brief — F003 Struct() reflect cache 效能優化

## Goal

消除 `MaskerMarshaler.Struct()` 在每次呼叫時重複 reflect 解析 struct type metadata 的開銷。
對同一個 struct type 重複呼叫（典型情境：API response masking）時，tag 解析與欄位分類只需做一次，
之後直接查 cache 取已解析好的 metadata，只執行值的遮罩運算。

另新增 `Format()`：直接輸出遮罩後的字串表示，過程中不配置新的 masked struct，
適合「只要看遮罩後字串、不需要遮罩後物件」的記 log / debug 情境。

**硬性前提：`Struct()` 的對外行為（回傳值、error 條件、遮罩結果）必須與現況完全一致（向後相容）。**
本次只改內部解析路徑為「查 cache」，不改任何遮罩語意。

## Tasks (numbered)

1. **type-cache — typeInfo cache 實作**
   - 在 `masker.go`（或新增 `type_cache.go`）定義 per-field metadata 與 type metadata 結構：
     - `fieldMask`：描述單一可遮罩欄位，至少含：欄位 index、解析後的 `MaskerType`（tag 值）、
       欄位 `reflect.Kind` 的分類資訊（string / struct / ptr / map / slice + slice element kind），
       以及是否走 `struct` / `mapstruct` 遞迴的旗標。未匯出欄位、無 `mask` tag 的欄位之處理規則也要能由 metadata 表達。
     - `typeInfo`：持有 `[]fieldMask`（依欄位宣告順序，僅含需處理的欄位）。
   - 新增 package-level cache：`var typeCache sync.Map`（key 為 `reflect.Type`，value 為 `*typeInfo`）。
   - 新增 `func cachedTypeInfo(t reflect.Type) *typeInfo`：cache miss 時呼叫 `buildTypeInfo(t)` 解析並用
     `LoadOrStore` 寫回；cache hit 直接回傳。
   - 新增 `func buildTypeInfo(t reflect.Type) *typeInfo`：對該 type 的欄位**只迭代一次**，
     解析 tag、分類 kind、跳過未匯出欄位，產出 `[]fieldMask`。
   - 新增供測試/benchmark 在 **同 package（`package masker`）** 內重置 cache 的內部 helper
     （例如 `func resetTypeCache()`），不對外匯出，不污染 public API。

2. **struct-use-cache — Struct() 改用 cache**
   - 重構 `MaskerMarshaler.Struct()`：取得 input 的 `reflect.Type` 後改呼叫 `cachedTypeInfo()`，
     依 `typeInfo.fields` 迭代執行遮罩，**不再於每次呼叫時讀 `selem.Type().Field(i).Tag.Get(tagName)`**。
   - 巢狀遞迴（`mask:"struct"` / `mask:"mapstruct"`、slice/ptr/interface of struct）仍呼叫 `m.Struct()` /
     `m.maskMapStructValue()`，使其各自亦受惠於 cache。
   - 既有 `masker.go:195-366` 的所有分支語意（string / struct / ptr / map / slice 各 element kind /
     interface、nil pointer/slice/map 處理、未匯出欄位略過、無 tag 欄位直接複製）必須逐一保留，輸出不得漂移。

3. **format-func — Format() 直出字串**
   - 新增 `func (m *MaskerMarshaler) Format(s interface{}) string`（method on `*MaskerMarshaler`，與既有 API 一致）。
   - 走 `cachedTypeInfo()` 取 metadata，將遮罩後的值寫入 `strings.Builder`，**不配置新的 masked struct**
     （不呼叫 `reflect.New` 建一個完整的 masked 物件再印）。
   - 輸出為確定性字串（deterministic）：巢狀 struct / ptr-to-struct 以遞迴方式 inline 展開遮罩後內容，
     **不得輸出記憶體位址**（與 `fmt %v` 對 pointer 欄位印 `0xc...` 的非確定性行為不同）。
   - 行為約束（feature rules）：`Format()` 不得 panic；nil input 回空字串 `""`；nil pointer 欄位以 `<nil>` 表示。
   - **遮罩語意 parity**：對同一 input，`Format()` 對每個 string 欄位產生的遮罩字串，必須與 `Struct()` 結果中
     對應欄位的值 byte-identical（共用同一條遮罩路徑，不可自行重寫遮罩邏輯）。
   - 輸出格式採 Go `%v` struct 風格：`{v0 v1 ... vn}`，欄位依宣告順序、以單一空白分隔；
     ptr-to-struct 欄位以 `&{...}` 遞迴展開（遮罩後、無位址），nil pointer 為 `<nil>`。

4. **benchmark — 效能 benchmark**
   - 新增 `masker_bench_test.go`（`package masker`），涵蓋：
     - `BenchmarkStruct`：對同一 type 反覆呼叫 `Struct()`（cache 已 warm），代表穩態開銷。
     - `BenchmarkStructCold`：每次 iteration 先 `resetTypeCache()` 再 `Struct()`，代表含 tag 解析的冷路徑。
     - `BenchmarkFormat`：對同一 type 反覆呼叫 `Format()`。
   - 以 `-benchmem` 取得 ns/op 與 allocs/op，作為效能 artifact。
   - 另新增確定性測試（見 AC）：用 `testing.AllocsPerRun` 比較 `Format()` 與 `Struct()` 的 allocs，
     並驗證 metadata 對同一 type 只建立一次。

## Scope (files to modify)

- `masker.go` — 重構 `Struct()` 改用 cache；可在此或新檔加入 cache 結構與 `Format()`。
- 新增 `type_cache.go`（建議）— `fieldMask` / `typeInfo` / `typeCache` / `cachedTypeInfo` / `buildTypeInfo` / `resetTypeCache`。
- 新增 `masker_bench_test.go` — benchmarks。
- 新增 / 擴充 `*_test.go` — cache、Format、parity、concurrency 測試（同 `package masker` 以存取內部 cache）。

## Out of Scope

- 不改任何個別 masker 型別（password/name/email…）的遮罩演算法。
- 不改 `Struct()` 的對外簽章、回傳型別、error 條件。
- 不新增 package-level wrapper 函式（現有 API 全為 `*MaskerMarshaler` method，維持一致）。
- 不處理 cyclic object graph（沿用現況限制，文件已載明）。
- 不為 `Format()` 設計可插拔的格式選項 / 多種輸出格式（只做單一確定性 `%v` 風格）。
- 不改 `Marshal()` / `Register()` / `Get()` / `List()` 等既有方法。
