# Coder Report — Round 1

## What Was Done

F006 是**驗收型（parity）feature**，目標是證明 `v3/` 完整承繼 v2（root package）的遮罩行為。v3 核心實作已由 F001/F003 完成，本輪**未修改任何 v3 原始碼**——移植的 v2 測試案例全部一次通過，未揭露任何 v3 與 v2 的行為不一致，因此不需要動 `v3/*.go`。

本輪產出為 4 個新測試檔，逐字移植 v2 既有測試的「輸入 → 期望輸出」，僅將 API 名稱適配 v3（常數改名、`WithMaskChar`、`NewAbuseMasker(maskChar)` 簽章）。期望輸出值與 v2 完全相同。

## Files Changed

- `v3/struct_parity_test.go`（新增）— Task 1 struct-migration：string 各 type、巢狀 struct（value/pointer）、`[]string`、`[]struct`、`[]*struct`、`[]interface{}`、`interface{}` 欄位、mapstruct 全表、nil 保留、非 struct 輸入 error、unexported 跳過、退化案例、generic tag。
- `v3/abuse_parity_test.go`（新增）— Task 2 abuse-migration：`AbuseTrie`（Insert/InsertAll/Contains/findAbuseWords）、`AbuseMasker`（各建構子/AddWords/AddWord/ContainsAbuse/GetAbuseWords/多次出現/空字典/no-substring/保留空白/不同 mask 字元）、`cleanWord`、abuse + struct tag 整合。
- `v3/abuse_loader_parity_test.go`（新增）— Task 2：補齊 v3 尚缺的 loader 案例（空輸入、純註解、混合內容）。
- `v3/api_surface_test.go`（新增）— Task 3 api-surface：預設 masker 集合 + `List()` parity、`DefaultMaskerMarshaler` 可用、未知型別 `Marshal`/`Get` 回傳 error、`Register` 覆蓋、並行 `Marshal`/`Struct` 安全。

## Verification

| Command | Result |
|---|---|
| `cd v3 && go build ./...` | PASS |
| `cd v3 && go vet ./...` | PASS（No issues found） |
| `cd v3 && go test -race -cover ./...` | PASS — 173 tests，total coverage **93.4%** |
| v2 root `*.go` 未修改檢查 | PASS（`git diff --name-only` 不含 root `*.go`） |

## v2 → v3 測試對照表（Task 4）

| v2 測試函式 | v3 對應測試 | 備註 |
|---|---|---|
| `TestMaskerMarshaler_Marshal` | `TestMarshalAndMustMarshal`（既有）、`TestMarshalAndGet_UnknownType` | 已知/未知型別 |
| `TestMaskerMarshaler_Register` | `TestRegisterGetListUnregister`（既有）、`TestRegisterOverride` | |
| `TestMaskerMarshaler_Unregister` | `TestRegisterGetListUnregister`（既有） | |
| `TestMaskerMarshaler_Get` | `TestRegisterGetListUnregister`（既有）、`TestMarshalAndGet_UnknownType` | |
| `TestMaskerMarshaler_List` | `TestNewMaskerMarshaler_DefaultMaskers` | 集合比對 |
| `TestMaskerMarshaler_SetMasker` | `TestWithMaskChar`（既有） | **parity 例外**：`SetMasker` → `WithMaskChar` |
| `TestMaskerMarshaler_Struct` | `TestStruct`（既有，string 欄位）、`TestStruct_StringFieldsAllTypes`、`TestStruct_NestedStructValue`、`TestStruct_NestedStructPointer`、`TestStruct_SliceString`、`TestStruct_NonStructInput` | 拆分為多個聚焦測試 |
| `TestNewMaskerMarshaler` | `TestNewMaskerMarshaler_DefaultMaskers` | 12 個靜態 masker |
| `Test_strLoop` / `Test_overlay` | 不移植 | 內部 helper，非 parity 範圍（已由上層行為間接涵蓋） |
| `TestMaskerMarshaler_MapStructs` | `TestStruct_MapStructs` | 逐字移植全表 |
| `TestStruct_TypedNilPointer` | `TestStruct_TypedNilPointer` | |
| `TestStruct_SlicePtrNilElement` | `TestStruct_SlicePtrNilElement` | |
| `TestStruct_StructFieldWithNonStructTag` | `TestStruct_StructFieldWithNonStructTag` | |
| `TestStruct_PtrFieldWithNonStructTag` | `TestStruct_PtrFieldWithNonStructTag` | |
| `TestStruct_SliceOfIntPreserved` | `TestStruct_SliceOfIntPreserved` | |
| `TestMarshal_AllMaskerOverridable` | `TestRegisterOverride` | |
| `TestMarshal_AbuseMaskerDefault` | `TestAbuseConvenience_EmptyDict`（既有）、`TestAbuseMasker_EmptyTrieBehavior` | |
| `TestMaskerMarshaler_ConcurrentMarshal` | `TestConcurrent_MarshalAndStruct` | 並行涵蓋 Marshal + Struct |
| `TestAbuseTrie` | `TestAbuseTrie_Full`、`TestAbuseTrie_Contains`（既有） | |
| `TestAbuseMasker` | `TestAbuseMasker_WithWords` | 含不同 mask 字元 |
| `TestAbuseMaskerAddWords` | `TestAbuseMasker_AddWordsAndAddWord` | |
| `TestAbuseMaskerContainsAbuse` | `TestAbuseMasker_ContainsAbuse` | |
| `TestAbuseMaskerGetAbuseWords` | `TestAbuseMasker_GetAbuseWords` | |
| `TestAbuseMaskerMultipleOccurrences` | `TestAbuseMasker_MultipleOccurrences` | |
| `TestAbuseMaskerEmptyTrie` | `TestAbuseMasker_EmptyTrieBehavior` | |
| `TestCleanWord` | `TestCleanWordParity` | |
| `TestAbuseMasker_NoSubstringMatch` | `TestAbuseMasker_NoSubstringMatchParity` | |
| `TestAbuseMasker_PreservesWhitespace` | `TestAbuseMasker_PreservesWhitespaceParity` | |
| `TestAbuseWordLoaderLoadFromString` | `TestAbuseWordLoader_LoadFromString`（既有） | |
| `TestAbuseWordLoaderLoadFromSlice` | `TestAbuseWordLoader_LoadFromSlice`（既有） | |
| `TestAbuseWordLoaderLoadFromReader` | `TestAbuseWordLoader_LoadFromReader`（既有） | |
| `TestAbuseWordLoaderEmptyInput` | `TestAbuseWordLoader_EmptyInput` | |
| `TestAbuseWordLoaderOnlyComments` | `TestAbuseWordLoader_OnlyComments` | |
| `TestAbuseWordLoaderMixedContent` | `TestAbuseWordLoader_MixedContent` | |
| `TestParseGenericMask` | `TestParseGenericMask`（既有 v3） | |
| `TestMarshal_GenericTags` | `TestMarshal_GenericTags`（既有 v3） | |
| `TestStruct_GenericTags` | `TestStruct_GenericTags` | v2 在 generic_test.go，v3 補進 struct_parity |
| （v2 無）`LoadFromFile` | `TestAbuseWordLoader_LoadFromFile`（既有 v3） | v3 額外涵蓋 |
| （v2 無）abuse+struct 整合 | `TestStruct_AbuseField` | AC-14 |

## docs/design/v2-parity.md 清單勾選

**Masker Types（14）**：none/password/name/addr/email/mobile/tel/id/credit/url/abuse/all/first-N/last-N — 全部於 `TestStruct_StringFieldsAllTypes`、`TestNewMaskerMarshaler_DefaultMaskers`、generic 測試、abuse 測試涵蓋 ✅

**Structural Tags（2）**：`struct`（巢狀 value/pointer/slice/interface）✅、`mapstruct`（map/ptr/slice/ptr-to-slice/巢狀 map）✅

**Struct() 處理能力**：string 欄位 ✅、巢狀 struct（value/ptr）✅、`[]string` ✅、`[]struct`/`[]*struct`/`[]interface{}` ✅、map（mapstruct 各組合）✅、`interface{}` 欄位 ✅、nil 保留 ✅、非 struct 輸入 error ✅、unexported 跳過 ✅

**API 表面**：`NewMaskerMarshaler` ✅、`DefaultMaskerMarshaler` ✅、`Marshal` ✅、`Struct` ✅、`Register` ✅、`Unregister` ✅、`Get` ✅、`List` ✅、~~`SetMasker`~~ → `WithMaskChar`（**parity 例外**，見下）✅

**Abuse 子系統**：`AbuseMasker`（`NewAbuseMasker`/`NewAbuseMaskerWithWords`）✅、`AbuseTrie` ✅、`AbuseWordLoader`（`LoadFromFile`/`LoadFromReader`，另含 `LoadFromString`/`LoadFromSlice`）✅

**非功能性**：Concurrency-safe（`-race` + `TestConcurrent_MarshalAndStruct`）✅、零外部依賴（`go vet` + build 無第三方 import）✅、Go 版本 ✅

## Parity 例外（行為一致、API 形狀刻意不同，非缺漏）

1. **`SetMasker` → `WithMaskChar`**：v3 改用 functional option，未新增 `SetMasker`（task-brief 明令禁止）。`TestWithMaskChar` 驗證自訂字元生效且不污染 `DefaultMaskerMarshaler`（L002）。
2. **`AbuseMasker.Mask(text)` vs v2 `Marshal(maskChar, text)`**：v3 將 mask 字元於建構時注入（`NewAbuseMasker(maskChar)`），遮罩入口改為 `Mask(text)`，輸出值與 v2 相同。
3. **`URLMasker` 固定輸出 `xxxxx`、`NoneMasker` 原樣回傳**：不受 mask 字元設定影響，期望值與 v2 一致（`TestWithMaskChar_URLAndNoneException` 驗證）。
4. v3 增量 `MustMarshal`、`convenience.go` package-level 短函式不在 v2 parity 範圍，但未破壞既有行為。

## Findings

- **無** v3 與 v2 行為不一致——所有移植測試一次通過，未修改 `v3/*.go`。
- **無** `docs/design/v2-parity.md` 缺漏需標記。
- `Test_strLoop` / `Test_overlay` 為 v2 內部 helper 的白盒測試，不屬對外 parity 範圍，未移植；其行為已由上層遮罩測試間接涵蓋。
