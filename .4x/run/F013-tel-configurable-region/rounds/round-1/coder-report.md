# Coder Report — Round 1

## What Was Done

在 `generic.go` 新增 `tel-` 前綴動態 tag 解析，支援可配置區碼／號碼／國際碼長度的市話遮罩。

1. **`cleanTelValue` helper**：清理電話號碼格式字元（`()- ` 及開頭 `+`）
2. **`parseTelMask` 函式**：解析 `tel-` tag（2/3/4 段），依消歧義規則判讀 intlLen/regionLen/numberLen/sep，組出遮罩輸出
3. **`parseGenericMask` 新增分派**：在 `first-`/`last-` 判斷之後加入 `tel-` 前綴檢查
4. **測試**：`TestParseGenericMask_Tel`（27 個 table case，涵蓋 AC-1 至 AC-27）、`TestMarshal_GenericTags_Tel`（端對端整合 + 既有 `TypeTel` 回歸）

既有 `telephone.go` 未修改（AC-23），`mask:"tel"` 行為不變（AC-21）。

## Files Changed

- `generic.go` — 新增 `cleanTelValue`、`parseTelMask`、`parseGenericMask` 加 `tel-` 分派（+128 行）
- `generic_test.go` — 新增 `TestParseGenericMask_Tel`、`TestMarshal_GenericTags_Tel`（+96 行）

## Verification

- `go build ./...`: pass
- `cd zapfield && go build ./...`: pass
- `go vet ./...`: pass
- `cd zapfield && go vet ./...`: pass
- `go test -race -count=1 ./...`: 241 tests pass
- `cd zapfield && go test -race -count=1 ./...`: 30 tests pass
- `gofmt -l generic.go generic_test.go`: 零輸出
- `git diff HEAD -- telephone.go`: 零輸出（未修改）
- `4x check F013-tel-configurable-region`: ✅ All checks passed
