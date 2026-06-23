# Acceptance Criteria

| # | Criterion | Verification Method |
|---|---|---|
| AC-1 | `Sensitive[T]` 定義於 `v3/sensitive.go`，為泛型 value struct，`raw`、`masked`、`mask` 三個 field 皆 unexported（首字小寫）。 | 程式碼審查 + `go vet ./...`；測試確認套件外無法直接存取欄位（只能透過 `Reveal()`）。 |
| AC-2 | `Reveal()` 是取得原值的唯一途徑，回傳建構時傳入的原始值；對任一建構子建立的值，`s.Reveal()` 等於原輸入。 | 單元測試：`NewPhone("0987654321").Reveal() == "0987654321"`，各建構子皆驗。 |
| AC-3 | `String()` 實作 `fmt.Stringer`，回傳遮罩值且不含原值。`fmt.Sprint`/`fmt.Println` 輸出為遮罩字串。 | 單元測試：`fmt.Sprint(NewPhone("0987654321"))` == `"0987***321"`，且不含 `"654"`。 |
| AC-4 | `GoString()` 實作 `fmt.GoStringer`，`%#v` 輸出為遮罩值，不洩漏 raw 或 struct 內部欄位。 | 單元測試：`fmt.Sprintf("%#v", NewEmail("ggw.chang@gmail.com"))` 等於該 email 的遮罩值，且不含原帳號 `"chang"`。 |
| AC-5 | `MarshalJSON()` 實作 `json.Marshaler`，`json.Marshal` 輸出遮罩值（正確 JSON 跳脫）；含 `Sensitive` 欄位的 struct 序列化得到遮罩結果。 | 單元測試：`json.Marshal(SMSRequest{Phone:NewPhone("0987654321"),Message:"驗證碼 1234"})` == `{"phone":"0987***321","message":"驗證碼 1234"}`。 |
| AC-6 | `MarshalText()` 實作 `encoding.TextMarshaler`，回傳遮罩值的 bytes。 | 單元測試：`NewID("A123456789").MarshalText()` 回傳 `[]byte("A12345****")`，err 為 nil。 |
| AC-7 | `LogValue()` 實作 `slog.LogValuer`，回傳遮罩值；經 `slog` 輸出時欄位為遮罩字串，不含原值。 | 單元測試：用 `slog.New(slog.NewTextHandler(buf))` 記錄含 `Sensitive` 的 attr，斷言 buf 含遮罩值、不含原值。 |
| AC-8 | 九個內建建構子（`NewPhone/NewEmail/NewPassword/NewID/NewCredit/NewName/NewAddress/NewTel/NewURL`）皆 `func(string) Sensitive[string]`，masked 輸出與對應 v3 package-level 函式逐字一致。 | 單元測試：對每個建構子，`c(v).String()` == 對應 v3 函式 `Fn(v)`（沿用既有 v3 測資作為斷言）。 |
| AC-9 | `NewPhone` 綁定 `Mobile`（手機）、`NewTel` 綁定 `Tel`（市話），兩者輸出格式不同且不混用。 | 單元測試：`NewPhone("0987654321").String()`=="0987***321"；`NewTel("0227993078").String()`=="(02)2799-****"。 |
| AC-10 | 通用建構子 `NewSensitive[T](raw T, maskFn func(T) string)` 支援自訂遮罩邏輯；建構時即執行 maskFn 並快取 masked。 | 單元測試：以自訂 `maskFn`（如全部變 `#`）建構非 string 型別（如 `int`）的 `Sensitive`，驗 `String()` 為自訂遮罩結果、`Reveal()` 為原值。 |
| AC-11 | `Equal(other)` 比較兩值原始值是否相等（用 `reflect.DeepEqual(raw)`），相同 raw → true，不同 raw → false；不因 masked 相同而誤判相等。 | 單元測試：相同 raw 兩建構子值 Equal 為 true；不同 raw 為 false；構造兩個不同 raw 但 masked 相同的情況驗證仍回 false。 |
| AC-12 | `UnmarshalJSON()` 能從原始 JSON 值還原 `Sensitive`：對先以建構子預填欄位的 struct 做 `json.Unmarshal`，`Reveal()` 得到 JSON 中的原值，且 masked 依綁定 maskFn 重算正確。 | 單元測試：`req := SMSRequest{Phone:NewPhone("")}`；`json.Unmarshal([]byte(`{"phone":"0987654321"}`), &req)`；斷言 `req.Phone.Reveal()=="0987654321"` 且 `req.Phone.String()=="0987***321"`。 |
| AC-13 | `mask == nil` 的退化情況不洩漏原值：`NewSensitive(raw, nil)` 與「未綁定 mask 的 `UnmarshalJSON`」皆使 masked 為空字串 `""`，安全輸出路徑（String/JSON/Text/Log）不回傳 raw。 | 單元測試：`NewSensitive("secret", nil).String()==""`；對 `var s Sensitive[string]` 直接 `UnmarshalJSON([]byte(`"secret"`))` 後 `s.String()==""` 且 `s.Reveal()=="secret"`。 |
| AC-14 | 五個安全輸出 interface 由 value type 滿足，並有編譯期斷言（`var _ fmt.Stringer = Sensitive[string]{}` 等）。 | 編譯通過即驗證；程式碼審查確認斷言存在。 |
| AC-15 | 提供可執行 example（至少 `ExampleNewPhone` 及 JSON 序列化 example），由 `go test` 驗證輸出。 | `cd v3 && go test ./...` example 全數 PASS。 |
| AC-16 | 全部測試含 race detector 通過，新增程式碼有合理覆蓋率（核心型別 + 九建構子 + 五 interface + Equal + UnmarshalJSON 皆有測試）。 | `cd v3 && go test -race -cover ./...` 全綠。 |
| AC-17 | exported 識別字皆有繁中 GoDoc（第一句以名稱開頭）；通過 gofmt / go vet。 | `gofmt -l v3/` 無輸出；`cd v3 && go vet ./...` 無錯；程式碼審查確認 GoDoc。 |
