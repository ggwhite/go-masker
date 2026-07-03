# F013 tel masker 支援可配置區碼／號碼長度

來源：https://github.com/ggwhite/go-masker/issues/40

## 概述

現行 `TelephoneMasker`（`mask:"tel"`）寫死「8 碼免區碼／11 碼固定 2 碼區碼」，
輸出固定為 `(02)2799-****`。實際市話區碼長度不只 2 碼（中國深圳 `755`、美國
`212` 等），且部分場景需要含國際碼（E.164 風格但去掉不必要的括號格式）。

本 feature 新增一組動態 mask tag，讓使用者在 struct tag 上直接宣告區碼／
號碼／（可選）國際碼位數，不需要註冊新 masker、不影響既有 `tel` tag 行為。

## 背景

- 原始 issue 只要求「區碼＋號碼」可配置，且明講輸出要「去掉國家碼」
- 討論過程中確認：這次一併支援「國際碼＋區碼＋號碼」的組合（見下方 tag 語法），
  作為 issue 訴求的延伸，非 issue 原文要求
- `mobile`／`id`／`password` 等其他 masker type 的可配置需求本次不處理，
  待有實際 issue 再另開 feature（本 feature 建立的動態 tag 模式可直接複用）

## 目標

比照 `generic.go` 現有的 `first-N`/`last-N` 動態 tag 模式（不需註冊、tag 直接
帶參數、命中時不查 `maskers` map），新增 `tel-` 前綴的動態 tag。

## Tag 語法

```
tel-<regionLen>-<numberLen>                          // 無國際碼，預設 dash 分隔
tel-<regionLen>-<numberLen>-<sep>                     // 無國際碼，指定分隔符
tel-<intlLen>-<regionLen>-<numberLen>                 // 含國際碼，預設 dash 分隔
tel-<intlLen>-<regionLen>-<numberLen>-<sep>           // 含國際碼，指定分隔符
```

- `regionLen`／`intlLen`：正整數，位數
- `numberLen`：整數，`>= 4`（含被遮罩的末 4 碼）
- `sep`：關鍵字 `dash`（輸出 `-`，省略時預設值）或 `space`（輸出 ` `）

### 消歧義規則

`tel-` 後以 `-` split 出的 token 數量：

- 2 段 → `[regionLen, numberLen]`
- 3 段 → 檢查第三個 token：
  - 是 `dash` 或 `space` → `[regionLen, numberLen, sep]`
  - 能 parse 成正整數 → `[intlLen, regionLen, numberLen]`
  - 兩者都不是 → 非法 tag，回傳 error
- 4 段 → `[intlLen, regionLen, numberLen, sep]`
- 其他段數 → 非法 tag，回傳 error

因為 `dash`/`space` 關鍵字與正整數字串不可能同時成立，此規則無歧義。

## 拆分與輸出規則

1. 清理輸入：沿用現行 `TelephoneMasker` 的清理邏輯（去除空格、`(`、`)`、`-`），
   額外去除開頭的 `+`（相容已格式化的國際碼輸入，如 `+886-2-2799-3078`）
2. 令 `total = intlLen + regionLen + numberLen`（無國際碼時 `intlLen = 0`），
   清理後長度為 `L`：
   - `L == total` → 依序左切：`intlLen` 碼國際碼／`regionLen` 碼區碼／
     `numberLen` 碼本地號碼
   - 其他長度（含多 1 碼的情況）→ 視為不合法輸入，回傳清理後的原值（不
     遮罩），與現行 `tel` 對長度不符的行為一致（非 error）
3. 輸出格式：`[+<國際碼><sep>]<區碼><sep>[<本地號碼前(numberLen-4)碼><sep>]<mask*4>`
   - 無國際碼時省略 `+<國際碼><sep>` 段
   - `numberLen == 4`（無明碼前綴）時省略中間的前置碼片段，直接輸出
     `<區碼><sep><mask*4>`

**不做 trunk prefix（國內前導 0）猜測性剝除。** 是否要去掉開頭 0、國碼與
區碼怎麼組，是電話號碼正規化的問題（且規則因國家而異，例如美國 NANP 號碼
沒有前導 0），不是遮罩函式庫該猜的事。呼叫端必須自行把輸入正規化成剛好
`intlLen+regionLen+numberLen` 碼再傳入，見「Non-goal」。

### 範例

| tag | 輸入（呼叫端已正規化） | 輸出 |
|---|---|---|
| `tel-2-8` | `0227993078` | `02-2799-****` |
| `tel-3-8` | `75588888888`（呼叫端已去掉開頭 0） | `755-8888-****` |
| `tel-3-7` | `2125551234` | `212-555-****` |
| `tel-2-3-8` | `8675588888888`（`86`＝國碼、`755`＝區碼，呼叫端已去掉開頭 0） | `+86-755-8888-****` |
| `tel-2-3-8-space` | `8675588888888` | `+86 755 8888-****` |

（`tel-2-3-8` 的 `2` 是 `intlLen`，`3` 是 `regionLen`；沿用「消歧義規則」——
第三個 token `8` 能 parse 成正整數，所以整組被判讀為
`[intlLen, regionLen, numberLen]`。）

## struct 宣告範例

```go
type Contact struct {
    HomePhone   string `mask:"tel"`              // 原有固定格式，行為不變
    TWPhone     string `mask:"tel-2-8"`           // 台灣市話（區碼含前導 0，如 "02"）
    CNPhone     string `mask:"tel-3-8-space"`     // 深圳市話，空格分隔
    TWPhoneIntl string `mask:"tel-3-1-8"`         // 含國際碼（呼叫端須先去掉區碼前導 0）
}
```

## 已知相容性：已格式化但未遮罩的輸入

輸入若已是「格式化但未遮罩」的原始號碼（含 `-`／空格／`()`／開頭 `+`），
清理步驟會統一去除這些字元後再切，因此不論輸入是 `8860227993078` 或
`+886-02-2799-3078`，結果一致。（相對地，「已遮罩過的值」不在此保證範圍內，
見下方 Non-goal。）

## Non-goal

- **不做電話號碼正規化**，包含「有國碼時去掉國內前導 0」這類轉換。這是
  電話號碼格式的語意問題，且規則因國家而異（如美國 NANP 號碼沒有前導 0），
  不是遮罩函式庫該猜的事。呼叫端必須自行把輸入正規化成剛好
  `intlLen+regionLen+numberLen` 碼（trunk prefix 已去除）再傳入，長度不符
  一律視為不合法輸入、不遮罩，go-masker 不做任何猜測性剝除。
- **不保證對「已遮罩過的值」再次遮罩仍正確**（例如把 `02-2799-****` 再丟進
  `tel-2-8` 遮一次）。這套邏輯是純位置切割，不檢查字元是否為數字或遮罩符號；
  結果是否正確純屬長度巧合。與現有所有 masker type（`password`／`mobile`／
  `id`⋯）一致，均未提供「偵測已遮罩並跳過」機制，這次也不新增。
- **separator 不支援 `dash`/`space` 以外的自訂符號**（如逗號）。這次只解決
  issue 提出的空格／dash 兩種需求，不做開放式自訂分隔符。
- **不處理 `mobile`／`id`／`password` 等其他 masker type 的可配置需求**，
  範圍僅限 `tel`。

## 向後相容

- 既有 `TypeTel`（`mask:"tel"`）與 `TelephoneMasker` 完全不動，行為與輸出
  格式不變
- 新 tag 為純新增的動態解析路徑，不影響任何既有呼叫

## 實作落點

- `generic.go`：`parseGenericMask` 新增 `tel-` 前綴分派，新增
  `parseTelMask` 函式；抽出既有 `TelephoneMasker.Mask` 的清理邏輯為共用
  helper（如 `cleanTelValue`），避免重複程式碼
- `telephone.go`：不修改
- `generic_test.go`：新增對應測試（見「測試策略」）

## 錯誤處理

`parseTelMask` 在以下情況回傳 error（比照 `first-N`/`last-N` 對非法 N 的
處理）：

- tag 分段數不是 2、3、4 段
- `intlLen`／`regionLen` 不是正整數
- `numberLen` 不是整數或 `< 4`
- 第三段既非 `dash`/`space` 關鍵字、也非合法正整數（3 段時）／第四段不是
  合法分隔符關鍵字（4 段時）

以下情況不是 error，回傳清理後的原值（不遮罩）：

- 清理後輸入長度不等於 `total`（含「多 1 碼」這種疑似含 trunk prefix 的
  情況——go-masker 不猜測，一律視為長度不符）
- 空字串（沿用 `l == 0 → return ""` 的既有慣例）

## 測試策略

`generic_test.go` 新增測試涵蓋：

- issue #40 三個範例（`tel-2-8`、`tel-3-8`、`tel-3-7`，輸入為呼叫端已正規化
  的長度）
- 情況 1：已格式化未遮罩輸入（含 `+`/`-`/空格混雜）清理後仍正確
- `space` 分隔符版本
- 國際碼版本（`tel-2-3-8`、含 `space`）
- `numberLen == 4`（無明碼前綴，省略空 prefix 片段）邊界
- 非法 tag：`intlLen<=0`（3/4 段時）、`regionLen<=0`、`numberLen<4`、分段數
  錯誤、第三段既非數字也非關鍵字 → 都回傳 error
- 長度不符（含剛好多 1 碼的情況，驗證「不猜測、不剝除」）→ 回傳原值不遮罩
- 回歸測試：確認既有 `tel`（`TypeTel`）行為與輸出完全不受影響
