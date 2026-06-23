# Task Brief — F001 v3 core: 新 Masker interface + package-level 短函式

## Goal

建立 go-masker v3 的核心基礎建設：在新的 `v3/` module 內，把 `Masker` interface 從 `Marshal(maskChar, value string) string` 改為 `Mask(value string) string`，mask char 改由 `MaskerMarshaler` 上的 functional option（`WithMaskChar`）配置；精簡 `MaskerType` 常數命名（`MaskerTypeMobile` → `TypeMobile`）；恢復 v1 風格的 package-level 便利函式（`masker.Mobile(v)` 等）；並把現有 12 種 masker（含 generic first-N/last-N）全部遷移到新 interface。所有遮罩「輸出結果」必須與 v2 完全一致——只改 API 形狀，不改遮罩邏輯。

此為 v3 第一個 feature，後續 feature（Sensitive[T]、Struct() reflect cache、zapfield/slog）都依賴此基礎。

## Tasks

### 1. module-setup — 建立 v3 module 結構
- 在 repo 內新增 `v3/` 子目錄，放置獨立 `go.mod`，module path：`github.com/ggwhite/go-masker/v3`，`go 1.21`。
- v3 core package 名稱維持 `masker`，**零外部依賴**（go.mod 不得出現 require 第三方套件）。
- 每個 masker 型別維持獨立檔案：`v3/masker.go`、`v3/mobile.go`、`v3/email.go`、`v3/password.go`、`v3/name.go`、`v3/address.go`、`v3/id.go`、`v3/credit.go`、`v3/telephone.go`、`v3/url.go`、`v3/none.go`、`v3/abuse.go`、`v3/abuse_loader.go`、`v3/generic.go`、`v3/convenience.go`。
- CI：更新既有 GitHub Actions workflow（或新增 job），讓 `v3/` 也跑 `go build` + `go test`；最低 Go 版本對齊 1.21。

### 2. new-masker-interface — 新 Masker interface + functional options
- 在 `v3/masker.go` 定義新 interface：
  ```go
  type Masker interface {
      Mask(value string) string
  }
  ```
- mask char 從方法簽名移除，改為配置層：
  - `type Option func(*MaskerMarshaler)`（或等效配置物件）
  - `func WithMaskChar(c rune) Option`
  - `func NewMaskerMarshaler(opts ...Option) *MaskerMarshaler`
- **契約**：未指定 option 時，預設 mask char 為 `'*'`，輸出與 v2 一致。`WithMaskChar` 設定後，該 marshaler 下**會套用 mask char 的** masker 改用新字元。
- **`WithMaskChar` 例外（必須遵守，否則違反「輸出與 v2 逐字一致」紅線）**：
  - `URLMasker` 使用 stdlib `url.Redacted()`，其密碼段固定輸出 `xxxxx`（寫死於標準函式庫，非本套件 mask char），**不受 `WithMaskChar` 影響**。即使 `WithMaskChar('#')`，`URL("http://u:p@host")` 仍為 `http://u:xxxxx@host`。coder 不得為了統一 mask char 而改寫 URL 遮罩邏輯。
  - `NoneMasker` 不做任何遮罩、原樣回傳，`WithMaskChar` 對它無意義、亦不適用。
  - `AbuseMasker` 以詞典命中後遮罩，若未載入詞典（default 空 trie）回原值，遮罩字元行為見下方 Task 4 的 Abuse 說明。
- 各 masker 實作如何取得 mask char 由 coder 決定（例如建構時注入內部欄位），但 interface 對外只暴露 `Mask(value string) string`。
- **實例化策略（避免污染 default）**：`DefaultMaskerMarshaler`（mask char `'*'`）與 `NewMaskerMarshaler(WithMaskChar('#'))` 必須各自持有獨立的一組 masker 實例，不可共用全域單例——否則在某個 marshaler 設定 `WithMaskChar` 會連帶改到 default 的輸出。每個 marshaler 在建構時自行建立其 masker 集合。

### 3. migrate-maskers — 遷移所有既有 masker
- 將下列 masker 從 `Marshal(maskChar, value string) string` 改為實作新 `Mask(value string) string`，**遮罩演算法與輸出逐字不變**：
  `MobileMasker`、`EmailMasker`、`PasswordMasker`、`NameMasker`、`AddressMasker`、`IDMasker`、`CreditMasker`、`TelephoneMasker`、`URLMasker`、`NoneMasker`、`AbuseMasker`、`AllMasker`。
- generic 動態 tag `first-N` / `last-N` 保留，行為不變（包含 N 為非負整數驗證、N 超過長度時截斷）。
- `MaskerType` 常數重新命名為精簡版（保留 `Type` prefix 避免衝突）：
  `TypeNone`、`TypePassword`、`TypeName`、`TypeAddress`、`TypeEmail`、`TypeMobile`、`TypeTel`、`TypeID`、`TypeCredit`、`TypeURL`、`TypeAbuse`、`TypeStruct`、`TypeAll`、`TypeMapStruct`。
- `NewMaskerMarshaler()` 與 `DefaultMaskerMarshaler` 用新常數註冊全部 masker。

### 4. convenience-functions — Package-level 便利函式
- 在 `v3/convenience.go` 提供頂層函式，底層走 `DefaultMaskerMarshaler`（mask char `'*'`），直接回傳 `string`（不回 error）：
  `Mobile(v) `、`Email(v)`、`Password(v)`、`Name(v)`、`Address(v)`、`ID(v)`、`Credit(v)`、`Tel(v)`、`URL(v)`、`Abuse(v)`、`None(v)`、`All(v)`。
- 提供動態版 `func Mask(t MaskerType, value string) string`（找不到 type 時的行為需在 AC 明確：回傳原值或空字串，二擇一並一致）。
- 命名對齊 design 範例：信用卡用 `Credit`、市話用 `Tel`（非 v1 的 `CreditCard`/`Telephone`）。
- **`Abuse(v)` 的預設行為**：`DefaultMaskerMarshaler` 註冊的 `AbuseMasker` 預設為空 trie（無詞典），故 `masker.Abuse(v)` 在未載入詞典時**回傳原值**。本 feature 不負責提供 default 詞典；`Abuse(v)` 的存在僅為 API 對稱性。tester 對 `Abuse(v)` 的斷言以「空字典回原值」為準，或先透過 `AbuseWordLoader` 載入詞典後再斷言遮罩結果（見 AC）。

### 5. marshal-cleanup — Marshal/MustMarshal 精簡
- `MaskerMarshaler.Marshal(t MaskerType, value string) (string, error)` 保留（給需要錯誤處理的呼叫者）。
- 新增 `MustMarshal(t MaskerType, value string) string`：找不到 masker type 時 panic，方便不想處理 error 的場景。
- `Struct()` 在本 feature **維持原樣可用**（其內部對 `Marshal` 的呼叫需配合新 interface 調整），但不做 reflect cache 重構（屬獨立 feature）。

## Scope

僅在新建的 `v3/` 子目錄內工作：
- `v3/go.mod`、`v3/go.sum`（若有）
- `v3/masker.go`、`v3/convenience.go`、`v3/generic.go`
- `v3/mobile.go`、`v3/email.go`、`v3/password.go`、`v3/name.go`、`v3/address.go`、`v3/id.go`、`v3/credit.go`、`v3/telephone.go`、`v3/url.go`、`v3/none.go`、`v3/abuse.go`、`v3/abuse_loader.go`
- 對應 `*_test.go`
- CI workflow 檔（`.github/workflows/*.yml`）新增 v3 job

## Out of Scope

- `Sensitive[T]` 泛型安全型別（獨立 feature）
- `Struct()` 的 reflect type cache 重構與 `Format()`（獨立 feature）— 本 feature 只需 `Struct()` 仍能運作
- zapfield sub-module、`log/slog` 整合（獨立 feature）
- 修改 root（v2）程式碼——v2 維持原狀，不動
- 改變任何 masker 的遮罩輸出結果
