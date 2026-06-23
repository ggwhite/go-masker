# Acceptance Criteria

| # | Criterion | Verification Method |
|---|---|---|
| AC-1 | `v3/go.mod` 存在，module path 為 `github.com/ggwhite/go-masker/v3`，且 `go 1.21`（含）以上 | `grep "module github.com/ggwhite/go-masker/v3" v3/go.mod` 且 `grep -E "^go 1\.(2[1-9]\|[3-9][0-9])" v3/go.mod` |
| AC-2 | v3 core package 零外部依賴：`v3/go.mod` 無任何 `require` 第三方套件 | 檢視 `v3/go.mod`，require 區塊為空或不存在 |
| AC-3 | 新 `Masker` interface 簽名為 `Mask(value string) string`，不含 mask char 參數 | `grep -A2 "type Masker interface" v3/masker.go` 顯示 `Mask(value string) string` |
| AC-4 | 提供 `WithMaskChar(c rune) Option` 與 `NewMaskerMarshaler(opts ...Option)`，且設定後遮罩字元生效；`WithMaskChar` 不污染 `DefaultMaskerMarshaler` | 單元測試：`NewMaskerMarshaler(WithMaskChar('#'))` 後 `Mobile` 遮罩輸出使用 `#`（如 `0987###321`）；未設定時為 `*`；同一程序內 default 的 `masker.Mobile(...)` 仍輸出 `*` 不受影響 |
| AC-4b | `WithMaskChar` 例外：`URLMasker` 不受影響（固定 `xxxxx`），`NoneMasker` 不適用 | 單元測試：`NewMaskerMarshaler(WithMaskChar('#'))` 下 `URL("http://u:p@host")` 仍為 `http://u:xxxxx@host`（密碼段固定 `xxxxx`，非 `#`）；`None(v)` 原樣回傳 |
| AC-5 | 12 種 masker 全部實作新 interface（mobile/email/password/name/addr/id/credit/tel/url/none/abuse/all），無殘留舊 `Marshal(s, i string) string` 方法 | `go build ./...` 通過；`grep -rn "Marshal(s" v3/*.go` 無 masker 型別命中 |
| AC-6 | 遮罩輸出與 v2 逐字一致（預設 mask char `*`）：`Mobile("0987654321")="0987***321"`、`Email("ggw@gmail.com")="ggw****@gmail.com"`、`Password("secret")="**************"`、`Credit("4111111111111111")="411111******1111"`、`Tel("0227993078")="(02)2799-****"`、`ID("A123456789")="A12345****"` 等 | 單元測試斷言各 masker 既有測資輸出不變（沿用 v2 test cases） |
| AC-7 | generic 動態 tag `first-N` / `last-N` 保留且行為不變（含 N 非負驗證、N 超長截斷） | 單元測試覆蓋 `first-3`、`last-4`、`first-0`、N 超過字串長度、非數字 N 回 error |
| AC-8 | `MaskerType` 常數採精簡命名 `TypeMobile`/`TypeEmail`/…，舊 `MaskerTypeXxx` 名稱不存在於 v3 | `grep -rnE "MaskerType(None\|Password\|Name\|Address\|Email\|Mobile\|Tel\|ID\|Credit\|URL\|Abuse\|Struct\|All\|MapStruct)" v3/*.go` 無命中；`grep "TypeMobile" v3/masker.go` 命中 |
| AC-9 | package-level 便利函式存在且回傳 `string`（不回 error）：`Mobile/Email/Password/Name/Address/ID/Credit/Tel/URL/Abuse/None/All` | 單元測試直接呼叫 `masker.Mobile("0987654321")` 等並斷言輸出 |
| AC-9b | `Abuse(v)` 在 default 空字典下回傳原值；載入詞典後可正確遮罩命中詞 | 單元測試：`masker.Abuse("hello")` 在無詞典時回 `"hello"`（原值）；另透過 `AbuseWordLoader` 載入詞典的 marshaler，斷言命中詞被遮罩 |
| AC-10 | 動態 `Mask(t MaskerType, value string) string` 存在，找不到 type 時行為一致且有測試明確定義（回原值或空字串） | 單元測試：`Mask(TypeMobile, v)` 正常；`Mask("nonexistent", v)` 行為符合文件定義 |
| AC-11 | `MaskerMarshaler.Marshal(t, value) (string, error)` 保留；新增 `MustMarshal(t, value) string`，未知 type 時 panic | 單元測試：`Marshal("bad", v)` 回 error；`MustMarshal("bad", v)` panic（`recover` 驗證）；合法 type 兩者輸出相同 |
| AC-12 | `Struct()` 在 v3 仍可用，對 tagged struct 遮罩結果與 v2 一致 | 單元測試：用含 `mask:"mobile"`/`mask:"email"` 等 tag 的 struct 跑 `Struct()`，斷言欄位被正確遮罩 |
| AC-13 | 全部測試通過且 race detector 乾淨 | `cd v3 && go test -race ./...` 全綠 |
| AC-14 | 所有 exported 識別字具備繁體中文 GoDoc，第一句以名稱開頭 | `go vet ./...` 通過；人工/Reviewer 抽查 exported 函式註解 |
| AC-15 | CI workflow 已涵蓋 v3：`v3/` 會被 build + test | 檢視 `.github/workflows/*.yml` 含針對 `v3/` 的 build/test 步驟 |
