# F014 mobile/id masker 支援可配置格式（多國支援）

## 概述

現行 `MobileMasker`（`mask:"mobile"`）和 `IDMasker`（`mask:"id"`）的遮罩
規則寫死台灣格式：手機固定遮罩第 4~7 位、身分證固定遮罩第 7~10 位。
其他國家的手機號碼和身分證件格式與台灣不同，無法直接使用。

本 feature 沿用 F013（`tel-` 動態 tag）建立的模式，新增 `mobile-` 和 `id-`
前綴的動態 tag，讓使用者在 struct tag 上直接宣告保留前幾碼、保留後幾碼，
中間全部遮罩。不需要註冊新 masker、不影響既有 `mobile` / `id` tag 行為。

## 背景

- F013 建立了動態 tag 的模式（`tel-`），本次直接複用
- `mobile` 和 `id` 的遮罩本質相同：「保留頭尾，遮罩中間」
- 既有的 `first-N` / `last-N` 動態 tag 只能保留頭或尾，缺少「同時保留頭尾」能力

## 目標

比照 `generic.go` 現有的動態 tag 模式，新增 `mobile-` 和 `id-` 前綴。
兩者共用相同的參數語意（`keepFront-keepEnd`），但各自有獨立的前綴以保留
語意清晰度——看到 tag 就知道這個欄位放的是手機號還是身分證。

## Tag 語法

```
mobile-<keepFront>-<keepEnd>
id-<keepFront>-<keepEnd>
```

- `keepFront`：非負整數，保留的前幾碼數量
- `keepEnd`：非負整數，保留的後幾碼數量
- `keepFront` 和 `keepEnd` 不可同時為 0（否則跟 `all` 等效，沒意義）

## 遮罩規則

1. 令 `L = len(value)`（原始字串長度，以 rune 計）
2. `keepFront + keepEnd >= L` → 原樣返回（無法遮罩）
3. `L == 0` → 原樣返回
4. 否則 → `value[:keepFront] + mask*(L-keepFront-keepEnd) + value[L-keepEnd:]`
   - `keepEnd == 0` 時省略尾部，即 `value[:keepFront] + mask*(L-keepFront)`

## 範例

### mobile

| tag | 輸入 | 國家 | 輸出 |
|---|---|---|---|
| `mobile` | `0987654321` | 台灣（預設，不變） | `0987***321` |
| `mobile-3-4` | `09012345678` | 日本（11位） | `090****5678` |
| `mobile-3-4` | `2025551234` | 美國（10位） | `202***1234` |
| `mobile-0-4` | `447911123456` | 英國（12位） | `********3456` |
| `mobile-4-2` | `0987654321` | 自訂 | `0987****21` |

### id

| tag | 輸入 | 國家 | 輸出 |
|---|---|---|---|
| `id` | `A123456789` | 台灣（預設，不變） | `A12345****` |
| `id-0-4` | `123456789` | 美國 SSN（9位） | `*****6789` |
| `id-4-0` | `123456789012` | 日本 My Number（12位） | `1234********` |
| `id-3-3` | `S1234567D` | 新加坡 NRIC（9位） | `S12***67D` |
| `id-2-1` | `AB1234567` | 英國 NI Number（9位） | `AB******7` |

### 邊界案例

| tag | 輸入 | 輸出 | 原因 |
|---|---|---|---|
| `mobile-5-5` | `1234567890`（10位） | `1234567890` | keepFront+keepEnd=10 >= L，無法遮罩 |
| `mobile-3-4` | `12345`（5位） | `12345` | keepFront+keepEnd=7 > L=5 |
| `id-3-0` | `` | `` | 空字串 |

## struct 宣告範例

```go
type User struct {
    // 台灣手機（預設，行為不變）
    Phone string `mask:"mobile"`

    // 日本手機 090-XXXX-XXXX（11位，保留前3後4）
    PhoneJP string `mask:"mobile-3-4"`

    // 美國 SSN XXX-XX-XXXX（9位，只露後4碼）
    SSN string `mask:"id-0-4"`

    // 台灣身分證（預設，行為不變）
    TaiwanID string `mask:"id"`

    // 日本 My Number（12位，保留前4）
    MyNumber string `mask:"id-4-0"`
}
```

## 向後相容

- 既有 `TypeMobile`（`mask:"mobile"`）與 `MobileMasker` 完全不動，行為不變
- 既有 `TypeID`（`mask:"id"`）與 `IDMasker` 完全不動，行為不變
- 新 tag 為純新增的動態解析路徑，不影響任何既有呼叫

## 實作落點

- `generic.go`：`parseGenericMask` 新增 `mobile-` 和 `id-` 前綴分派，
  新增 `parseMobileMask` 和 `parseIDMask` 函式（兩者內部共用相同的
  keepFront/keepEnd 解析與遮罩邏輯，可抽為共用 helper）
- `mobile.go` / `id.go`：不修改
- `generic_test.go`：新增對應測試

## 錯誤處理

回傳 error 的情況（比照 `first-N` / `last-N` / `tel-` 的處理）：

- tag 分段數不是 2 段（`mobile-` 或 `id-` 後面不是恰好兩個 `-` 分隔的數字）
- `keepFront` 或 `keepEnd` 無法 parse 為非負整數
- `keepFront` 和 `keepEnd` 同時為 0

不是 error，回傳原值（不遮罩）的情況：

- `keepFront + keepEnd >= len(value)`
- 空字串

## 測試策略

`generic_test.go` 新增測試涵蓋：

### mobile

- 台灣格式：`mobile-4-3`，輸入 `0987654321`，預期 `0987***321`（等同預設行為）
- 日本格式：`mobile-3-4`，輸入 `09012345678`，預期 `090****5678`
- 美國格式：`mobile-3-4`，輸入 `2025551234`，預期 `202***1234`
- 全遮罩前段：`mobile-0-4`，輸入 `447911123456`，預期 `********3456`
- keepEnd=0：`mobile-4-0`，輸入 `0987654321`，預期 `0987******`

### id

- 台灣格式：`id-6-0`，輸入 `A123456789`，預期 `A12345****`（等同預設行為）
- 美國 SSN：`id-0-4`，輸入 `123456789`，預期 `*****6789`
- 日本 My Number：`id-4-0`，輸入 `123456789012`，預期 `1234********`
- 新加坡 NRIC：`id-3-3`，輸入 `S1234567D`，預期 `S12***67D`

### 邊界與錯誤

- keepFront+keepEnd >= L → 原值不遮罩
- 空字串 → 原值
- 非法 tag（缺參數、非數字、同時為 0） → error
- 回歸測試：確認既有 `mobile` / `id`（無參數）行為完全不受影響

## Non-goal

- **不做 region preset**（`mobile:jp`、`id:us` 等預設組合）——本次只做
  動態參數，未來有需求再加
- **不做輸入正規化**——跟 F013 一致，呼叫端自行處理格式清理
- **不改 Masker interface 或 Option 機制**
