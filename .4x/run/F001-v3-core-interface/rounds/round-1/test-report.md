# Test Report — Round 1

## Summary
PASS — 16/16 criteria met

`4x verify` 三項指令全綠（build / vet / test -race），`go test -race` 32 個 top-level 測試（含子測試共 91 個）全數通過，race detector 乾淨，coverage 63.9%。

## Results
| # | Criterion | Status | Evidence |
|---|---|---|---|
| AC-1 | v3/go.mod module path + go 1.21+ | PASS | `v3/go.mod`：`module github.com/ggwhite/go-masker/v3`、`go 1.21` |
| AC-2 | v3 core 零外部依賴 | PASS | `grep -c require v3/go.mod` = 0，無 require 區塊 |
| AC-3 | Masker interface 簽名 Mask(value string) string | PASS | `grep -A3 "type Masker interface"` → `Mask(value string) string`（無 mask char 參數） |
| AC-4 | WithMaskChar + NewMaskerMarshaler 生效且不污染 default | PASS | `TestWithMaskChar` PASS：`NewMaskerMarshaler(WithMaskChar('#'))` → `0987###321`；default + package-level `Mobile` 仍 `0987***321` |
| AC-4b | WithMaskChar 例外：URL 固定 xxxxx、None 不適用 | PASS | `TestWithMaskChar_URLAndNoneException` PASS；GoDoc 亦明載例外 |
| AC-5 | 12 種 masker 實作新 interface，無殘留舊 Marshal(s,i) | PASS | `go build ./...` exit 0；`grep -rn "Marshal(s" v3/*.go` 無命中 |
| AC-6 | 遮罩輸出與 v2 逐字一致 | PASS | 各 masker test PASS，斷言含 `ggw****@gmail.com`、`A12345****`、`(02)2799-****` 等 v2 測資 |
| AC-7 | generic first-N/last-N 行為不變 | PASS | `TestParseGenericMask`、`TestMarshal_GenericTags`、`TestMarshal_GenericTags_Error` 全 PASS（含非數字 N 回 error、超長截斷） |
| AC-8 | MaskerType 精簡命名 TypeXxx，舊名不存在 | PASS | 舊 `MaskerTypeXxx` grep 無命中；`TypeMobile MaskerType = "mobile"` 存在於 masker.go:22 |
| AC-9 | package-level 便利函式回傳 string | PASS | 簽名 `func Mobile(value string) string` 等皆回 string；`TestConvenienceFunctions` PASS |
| AC-9b | Abuse 空字典回原值，載入後可遮罩 | PASS | `TestAbuseConvenience_EmptyDict`（`Abuse("hello")=="hello"`）、`TestAbuse_LoadedDictionary` PASS |
| AC-10 | 動態 Mask(t, value)，找不到 type 行為一致 | PASS | `TestMaskDynamic` PASS：`Mask("nonexistent", "value")=="value"`（回原值，符合文件定義） |
| AC-11 | Marshal 保留 + MustMarshal panic | PASS | `TestMarshalAndMustMarshal` PASS：`Marshal("bad")` 回 error、`MustMarshal("bad")` panic（recover 驗證）、合法型別兩者輸出相同 |
| AC-12 | Struct() 可用且與 v2 一致 | PASS | `TestStruct` PASS（含 mask tag 的 struct 欄位正確遮罩） |
| AC-13 | 全部測試通過且 race 乾淨 | PASS | `cd v3 && go test -race -cover ./...` exit 0；`ok ... coverage: 63.9%`；32 PASS / 0 FAIL |
| AC-14 | exported 繁中 GoDoc，名稱開頭 | PASS | `go vet ./...` No issues；抽查 `NewMaskerMarshaler`/`WithMaskChar`/`Mobile`/`Mask` GoDoc 首句皆繁中且以名稱開頭 |
| AC-15 | CI workflow 涵蓋 v3 | PASS | `.github/workflows/go.yml` 含 `build-v3` job、`working-directory: v3` 跑 build + test |

## Rules Check
- ✅ 既有 masker 遮罩邏輯不變，只改 interface 簽名（AC-6 逐字一致佐證）
- ✅ core package 零外部依賴（AC-2）
- ✅ exported 函式具繁中 GoDoc、名稱開頭（AC-14）
- ✅ 每個 masker 保留獨立檔案（mobile.go/email.go/... 皆存在）

## Verdict
PASS
