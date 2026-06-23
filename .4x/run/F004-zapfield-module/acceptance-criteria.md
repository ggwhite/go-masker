# Acceptance Criteria

| # | Criterion | Verification Method |
|---|---|---|
| AC-1 | `v3/zapfield/go.mod` 存在，module path 為 `github.com/ggwhite/go-masker/v3/zapfield`，`go 1.21`，且 require `github.com/ggwhite/go-masker/v3` 與 `go.uber.org/zap`，並有 `replace github.com/ggwhite/go-masker/v3 => ../`。 | `cat v3/zapfield/go.mod`；`grep` 檢查 module path、require、replace。 |
| AC-2 | zapfield 依賴隔離成立：`v3/go.mod`（core）**不含**任何 `go.uber.org/zap` 依賴。 | `! grep -q zap v3/go.mod`（無 match 才算過）。 |
| AC-3 | zapfield package 可編譯通過。 | `cd v3/zapfield && go build ./...` exit 0。 |
| AC-4 | 為下列每個 masker type 各提供一個 `func(key, value string) zap.Field` helper：Phone、Email、Password、Name、Address、ID、Credit、Tel、URL、Abuse、None、All。 | `go doc` 或 `grep -c "func .*(key, value string) zap.Field" v3/zapfield/field.go` ≥ 12；單元測試逐一呼叫。 |
| AC-5 | 每個 Field helper 回傳的 `zap.Field` 之 `Key` 等於傳入 key，且其字串值等於對應 core 函式輸出（逐字一致）。例：`zapfield.Phone("phone","0987654321")` 的值 == `masker.Mobile("0987654321")`（`"0987***321"`）。 | 單元測試：建構 field，斷言 `f.Key == key` 且 `f.String == masker.Xxx(value)`。 |
| AC-6 | 所有 Field helper 以 `zap.String` 建構（`f.Type == zapcore.StringType`），不得使用 `zap.Any`／reflection。 | 單元測試斷言 `f.Type == zapcore.StringType`；`! grep -E "zap.Any|zap.Reflect" v3/zapfield/*.go`。 |
| AC-7 | `Sensitive[T any](key string, s masker.Sensitive[T]) zap.Field` 存在，回傳 `zap.String(key, s.String())`，輸出值等於 `s.String()`（即 Sensitive 快取的 masked 值），且型別為 `zapcore.StringType`。 | 單元測試：`s := masker.NewPhone("0987654321")`；`f := zapfield.Sensitive("phone", s)`；斷言 `f.String == s.String()` 且 `f.Key == "phone"`。 |
| AC-8 | `Sensitive` adapter 不洩漏原值：對非 string T（如 `Sensitive[int]`）與 zero-value `Sensitive[T]`（mask 為 nil）皆只輸出 masked，不出現 raw。zero-value 情形輸出空字串。 | 單元測試：`var z masker.Sensitive[string]`；`f := zapfield.Sensitive("k", z)`；斷言 `f.String == ""`。另測一個自訂 `masker.NewSensitive[int]` 案例輸出 == `s.String()`。 |
| AC-9 | 全部 zapfield 測試通過（含 race detector）。 | `cd v3/zapfield && go test -race ./...` exit 0。 |
| AC-10 | 所有 exported 識別字（helper、adapter、package）皆有繁體中文 GoDoc，第一句以識別字名稱開頭。 | `go vet ./...` exit 0；人工/`go doc` 抽查。 |
| AC-11 | v3 core 未被修改、core 測試仍全綠（確認本 feature 未動 core）。 | `cd v3 && go test ./...` exit 0；`git status` 確認 core 檔無改動。 |
