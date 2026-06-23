# Task Brief — F005 zap Core interceptor: field name keyword 攔截

## Goal

在 `v3/zapfield` sub-module 新增一層 **zap.Core wrapper**，於 log 寫出前根據 **field name keyword / regex** 自動遮罩 string field。定位是「最後一道防線」——適用於改不動業務 code、但想對既有 logger 加一層被動防護的場景。

設計上必須誠實標示其侷限：keyword matching 本質是猜測，會誤遮（`phone_count`）也會漏掉（field name 不在清單）。若業務 code 已改用 `Sensitive[T]` 或 `zapfield.Phone()` 等顯式 helper，這層攔截即為多餘，不應視為主力方案。

對外 API 範例（與 feature 描述一致，WrapCore 不回傳 error）：

```go
core := zapfield.WrapCore(originalCore, zapfield.InterceptRules{
    Keywords: []string{"phone", "password", "token", "secret"},
})
logger := zap.New(core)
```

## Tasks (numbered, specific)

1. **InterceptRules 結構**（subtask: intercept-rules）— 在新檔 `v3/zapfield/core.go` 定義：
   ```go
   type InterceptRules struct {
       Keywords []string         // 對 field key 做 case-insensitive 子字串比對（strings.Contains）
       Patterns []*regexp.Regexp // 對 field key 做 regex 比對（已編譯，避免 WrapCore 回傳 error）
   }
   ```
   - 提供未匯出的判定方法 `func (r InterceptRules) match(key string) bool`：任一 keyword（小寫化後 `strings.Contains`）或任一 pattern `MatchString` 命中即回傳 true。
   - Keywords 與 Patterns 皆為空（或 nil）時，`match` 一律回傳 false（退化為「不攔截、原樣放行」，最安全且不 panic）。

2. **WrapCore 實作**（subtask: wrap-core）— 在 `v3/zapfield/core.go`：
   - `func WrapCore(core zapcore.Core, rules InterceptRules) zapcore.Core`，回傳一個內嵌 / 持有 `zapcore.Core` 與 `rules` 的 wrapper 型別（如 `maskingCore`）。
   - 實作 `zapcore.Core` 介面全部方法：
     - `Enabled(zapcore.Level) bool`、`Sync() error` → 直接委派內層 core。
     - `With(fields []zapcore.Field) zapcore.Core` → 先以 `interceptFields` 處理 fields，再委派內層 `With`，並用同樣的 `rules` 重新包裝回傳的 core（確保 `logger.With(...)` 累積的 context field 也被攔截）。
     - `Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry` → 依 `Enabled` 慣例 `ce.AddCore(ent, c)` 把自己加入（標準 wrapper 寫法）。
     - `Write(ent zapcore.Entry, fields []zapcore.Field) error` → 先 `interceptFields(fields)`，再委派內層 `Write`。
   - 私有 helper `func (c *maskingCore) interceptFields(fields []zapcore.Field) []zapcore.Field`：
     - 逐一檢查 field；**只有當 `field.Type == zapcore.StringType` 且 `rules.match(field.Key)` 為 true** 時，才把 `field.String = masker.All(field.String)`（等價 `masker.Mask("all", ...)`，每字元換成遮罩字元）。
     - 非 string field（Int/Bool/Reflect/Object…）一律原樣保留，避免誤殺數字/結構。
     - 不就地改寫呼叫者傳入的 slice 元素時若有疑慮，可建立 copy；以不影響呼叫端為準。

3. **doc 註解警告**（與 wrap-core 同檔）— `WrapCore` / `InterceptRules` 的 GoDoc 必須以繁體中文明確警告 false positive / false negative 風險，並指出「業務 code 已用顯式遮罩時本層多餘」。GoDoc 第一句以識別字名稱開頭。
   - ⚠️ 撰寫 doc 時**避免逐字寫出** CI grep 會誤判的禁用 token（如 `zap.Any` / `zap.Reflect`）；用語意描述。masking 一律走 `field.String` 欄位寫值，不用會觸發 reflection 的欄位建構子。

4. **攔截測試**（subtask: tests）— 在 `v3/zapfield/core_test.go`，使用 `go.uber.org/zap/zaptest/observer` 建立可觀測 core 後 `WrapCore` 包裝，斷言：
   - keyword 命中 → string field 值被全遮罩。
   - regex pattern 命中 → string field 值被全遮罩。
   - 非 string field（如 `zap.Int("phone_count", 3)`）即使 key 命中 keyword 也**不**被改動。
   - 未命中的 string field 原樣保留。
   - 空 `InterceptRules`（nil keywords + nil patterns）→ 所有 field 原樣放行。
   - case-insensitive：keyword `"phone"` 命中 key `"Phone"` / `"userPhone"`。
   - `logger.With(...)` 加入的 context field 也會被攔截（驗證 With 路徑）。

## Scope (which files/dirs to modify)

- 新增 `v3/zapfield/core.go`
- 新增 `v3/zapfield/core_test.go`
- 不改動 `v3/zapfield/go.mod`（zap 與 multierr 已是現有依賴；observer 屬 zap module，無需新增 require）

## Out of Scope

- 不修改 v3 core module 任何檔案（`masker.All` / `masker.Mask` 已存在，直接使用）。
- 不新增 first-N / last-N 等動態遮罩型別選擇——本層固定用 `all`（全遮罩）以求最安全。
- 不處理巢狀 object / array field 內部欄位的遞迴攔截（zap 的 ObjectMarshaler 內容在此層不可見）；僅處理頂層 StringType field。
- 不提供 keyword → 不同 masker type 的對應表（非本 feature 目標，屬未來擴充）。
- 不改動 `field.go`（顯式 helper）與其測試。
