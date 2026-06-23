# Acceptance Criteria

| # | Criterion | Verification Method |
|---|---|---|
| AC-1 | `zapfield.WrapCore(core zapcore.Core, rules InterceptRules) zapcore.Core` 存在，簽名不回傳 error，回傳值滿足 `zapcore.Core`，可直接傳入 `zap.New(core)`。 | `go build ./...`；測試中 `zap.New(WrapCore(...))` 編譯通過並能寫 log。 |
| AC-2 | `InterceptRules` 結構含 `Keywords []string` 與 `Patterns []*regexp.Regexp` 兩個 exported 欄位。 | code review + 測試直接以 struct literal 建構兩欄位。 |
| AC-3 | keyword 命中時，對應 string field 的值被全遮罩（等同 `masker.All(value)`，每字元換成遮罩字元）。 | observer 測試：`zap.String("password", "secret")` 在 rules `Keywords:["password"]` 下，輸出值 == `masker.All("secret")`。 |
| AC-4 | keyword 比對為 **case-insensitive 子字串**（`strings.Contains`）：`"phone"` 命中 `"Phone"`、`"userPhone"`。 | observer 測試多組 key 大小寫 / 前後綴皆遮罩。 |
| AC-5 | regex pattern 命中時，對應 string field 的值被全遮罩。 | observer 測試：`Patterns:[regexp.MustCompile("(?i)token")]` 命中 key `"access_token"`，值被遮罩。 |
| AC-6 | **非 string field 不被改動**：`zap.Int("phone_count", 3)` 在 keyword `"phone"` 命中 key 的情況下，整數值維持 `3` 不變。 | observer 測試斷言 Int field 的 Integer/值不變。 |
| AC-7 | 未命中的 string field 原樣保留；空 `InterceptRules`（nil Keywords + nil Patterns）下所有 field 原樣放行、不 panic。 | observer 測試：未命中 key 與空 rules 兩情境，輸出與輸入逐字一致。 |
| AC-8 | `logger.With(fields...)` 加入的 context field 同樣被攔截遮罩（With 路徑生效）。 | observer 測試：`logger.With(zap.String("secret","x")).Info("m")`，輸出該 field 被遮罩。 |
| AC-9 | `WrapCore` 與 `InterceptRules` 具繁體中文 GoDoc，明確警告 keyword matching 的 false positive（如 `phone_count`）與 false negative（不在清單即漏遮）風險，並說明業務 code 已用顯式遮罩時本層多餘。 | code review 檢視 GoDoc 內容。 |
| AC-10 | 不修改 v3 core module；變更僅限新增 `v3/zapfield/core.go` 與 `v3/zapfield/core_test.go`。 | `git status` 僅顯示這兩個新檔（及狀態檔）；core module 檔案無 diff。 |
| AC-11 | 全測試通過且帶 race detector；既有 `field_test.go` / `sensitive_test.go` 不受影響。 | `cd v3/zapfield && go test -race ./...` 全綠。 |
| AC-12 | doc/註解不逐字出現會觸發 CI grep 誤判的禁用 token（reflection 相關欄位建構子名稱）；遮罩一律經 `field.String` 寫值。 | grep 檢查 `core.go` 不含被禁 token 字面；code review。 |
