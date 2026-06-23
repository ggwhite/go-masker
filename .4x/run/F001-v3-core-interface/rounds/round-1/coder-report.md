# Coder Report — Round 1

## What Was Done

建立 go-masker v3 子 module，把 `Masker` interface 從 `Marshal(maskChar, value string) string` 改為 `Mask(value string) string`，遮罩字元改由 `MaskerMarshaler` 上的 functional option `WithMaskChar` 配置。全部 12 種 masker 遷移至新 interface，遮罩輸出與 v2 逐字一致（只改 API 形狀，遮罩演算法不動）。

具體完成項目：

- **module-setup**：新增 `v3/go.mod`（module path `github.com/ggwhite/go-masker/v3`、`go 1.21`、零外部依賴）。每個 masker 型別維持獨立檔案。
- **new-masker-interface**：`v3/masker.go` 定義新 `Masker` interface、`type Option func(*MaskerMarshaler)`、`WithMaskChar(c rune) Option`、`NewMaskerMarshaler(opts ...Option)`。各 masker 於建構時注入內部 `mask` 欄位；每個 marshaler 建立各自獨立的一組 masker 實例，`WithMaskChar` 不污染 `DefaultMaskerMarshaler`。
- **WithMaskChar 例外**：`URLMasker`（密碼段固定 `xxxxx`，走 stdlib `url.Redacted()`）與 `NoneMasker`（原樣回傳）不持有 mask char、不受影響。
- **migrate-maskers**：mobile/email/password/name/addr/id/credit/tel/url/none/abuse/all 全部改實作 `Mask(value string) string`，無殘留舊 `Marshal(s, i string)`。`MaskerType` 常數精簡為 `TypeXxx`（移除 `Masker` prefix）。
- **convenience-functions**：`v3/convenience.go` 提供 `Mobile/Email/Password/Name/Address/ID/Credit/Tel/URL/Abuse/None/All` 與動態 `Mask(t, value)`（未知型別回原值）。底層走 `DefaultMaskerMarshaler`。`Abuse(v)` 空字典回原值。
- **marshal-cleanup**：保留 `Marshal(t, value) (string, error)`，新增 `MustMarshal(t, value) string`（未知型別 panic）。`Struct()` 維持原樣可用（內部改呼叫 `masker.Mask(value)`），未做 reflect cache 重構（屬獨立 feature）。
- **CI**：`.github/workflows/go.yml` 新增 `build-v3` job，於 `v3/` 目錄以 go 1.21 跑 build + test。

## Files Changed

- `v3/go.mod` — 新 module 定義，go 1.21，零外部依賴
- `v3/masker.go` — 新 Masker interface、Option/WithMaskChar、MaskerMarshaler、Marshal/MustMarshal、Struct、精簡 MaskerType 常數、NewMaskerMarshaler、DefaultMaskerMarshaler、helper（strLoop/overlay）
- `v3/mobile.go`、`v3/email.go`、`v3/password.go`、`v3/name.go`、`v3/address.go`、`v3/id.go`、`v3/credit.go`、`v3/telephone.go` — 改 `Mask(value)`，注入 `mask` 欄位
- `v3/url.go` — `Mask(value)`，無 mask char（固定 xxxxx）
- `v3/none.go` — `Mask(value)`，原樣回傳
- `v3/generic.go` — `AllMasker.Mask` + `parseGenericMask`（first-N/last-N）
- `v3/abuse.go` — `AbuseMasker` 持有 mask char，`Mask(text)`；trie 結構不變
- `v3/abuse_loader.go` — loader 不變（無 mask char 依賴），註解改繁中
- `v3/convenience.go` — package-level 便利函式 + 動態 Mask
- `v3/*_test.go`（mobile/email/password/name/address/id/credit/telephone/url/generic/masker/convenience/abuse/abuse_loader）— 沿用 v2 測資 + v3 新增（WithMaskChar 隔離、URL/None 例外、MustMarshal panic、Struct、Abuse 空字典/載入詞典）
- `.github/workflows/go.yml` — 新增 build-v3 job（go 1.21，working-directory: v3）

## Verification

- `cd v3 && go build ./...` — 通過（無輸出）
- `cd v3 && go vet ./...` — 通過（No issues found）
- `cd v3 && go test -race -cover ./...` — 通過，91 個測試全綠，race detector 乾淨，coverage 63.9%
- `gofmt -l v3/` — 乾淨（無未格式化檔案）
- AC-3 `grep -A2 "type Masker interface" v3/masker.go` — 顯示 `Mask(value string) string` ✅
- AC-5 `grep -rn "Marshal(s" v3/*.go` — 無命中（無殘留舊簽名）✅
- AC-8 `grep -rnE "MaskerType(None|Password|...)" v3/*.go` — 無命中（舊常數已移除）✅

> 備註：coverage 63.9% 未達更高，主因為 `Struct()` 的 slice/map/ptr reflect 分支屬本 feature out-of-scope（reflect cache 為獨立 feature），未針對這些容器路徑補測；核心 interface、各 masker、convenience、Option 路徑皆已覆蓋。
