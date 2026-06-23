# Task Brief — F002 Sensitive[T] 泛型安全型別

## Goal

在 v3 module（`github.com/ggwhite/go-masker/v3`，Go 1.21）新增 `Sensitive[T]` 泛型安全型別，把遮罩從「開發者記得做」變成「洩漏原值要刻意呼叫 `.Reveal()`」。

核心精神：任何「自動」輸出路徑（`fmt` 列印、`%#v`、`json.Marshal`、`encoding.TextMarshaler`、`slog` 結構化日誌）一律輸出**遮罩值**；原始值只有透過唯一出口 `Reveal()` 才能取得。建構時就算好 masked string 並快取，之後取值零額外成本。

底層遮罩運算必須沿用 F001 已完成的 v3 package-level 函式（`Mobile/Email/Password/Name/Address/ID/Credit/Tel/URL`，皆走 `DefaultMaskerMarshaler`），**不得自行重寫遮罩邏輯**，確保 masked 輸出與既有 v3 行為逐字一致。

## Tasks (numbered, specific)

1. **`Sensitive[T]` 核心 struct（subtask: sensitive-struct）** — 在 `v3/sensitive.go` 定義泛型 struct：
   ```go
   type Sensitive[T any] struct {
       raw    T                 // unexported，唯一原值來源
       masked string            // 建構時算好並快取的遮罩字串
       mask   func(T) string    // 綁定的遮罩函式，供 UnmarshalJSON 重算 masked 用
   }
   ```
   - `raw`、`masked`、`mask` 三個 field 全部 unexported（規則要求）。
   - 提供 `func (s Sensitive[T]) Reveal() T`：回傳 `s.raw`，是取得原值的**唯一**途徑（value receiver，回傳 raw 本身）。

2. **通用建構子 `NewSensitive`（subtask: constructors）** — 在 `v3/sensitive.go`：
   ```go
   func NewSensitive[T any](raw T, maskFn func(T) string) Sensitive[T]
   ```
   - 建構時立即執行 `maskFn(raw)` 算出 masked 並快取，同時保存 `maskFn` 於 `mask` field。
   - 所有內建建構子都委派給它。`maskFn` 為 nil 時的行為見 Task 7 的 edge case 規範。

3. **安全輸出 interface 實作（subtask: safe-output）** — 在 `v3/sensitive.go`，皆為 value receiver、皆只回傳 `s.masked`（絕不碰 `s.raw`）：
   - `func (s Sensitive[T]) String() string` — 實作 `fmt.Stringer`，回傳 masked。
   - `func (s Sensitive[T]) GoString() string` — 實作 `fmt.GoStringer`，回傳 masked（防 `%#v` 洩漏 struct 內部）。
   - `func (s Sensitive[T]) MarshalJSON() ([]byte, error)` — 實作 `encoding/json.Marshaler`，把 masked 字串做 JSON 編碼後輸出（須正確處理跳脫，建議 `json.Marshal(s.masked)`）。
   - `func (s Sensitive[T]) MarshalText() ([]byte, error)` — 實作 `encoding.TextMarshaler`，回傳 `[]byte(s.masked)`。
   - `func (s Sensitive[T]) LogValue() slog.Value` — 實作 `log/slog.LogValuer`，回傳 `slog.StringValue(s.masked)`。
   - 須在檔案以 `var _ fmt.Stringer = Sensitive[string]{}`（及其餘四個 interface）做編譯期斷言，確保 value type 滿足介面。

4. **內建 masker type 建構子（subtask: constructors）** — 在 `v3/sensitive_constructors.go`，每個皆 `func(v string) Sensitive[string]`，委派 `NewSensitive` 並綁定對應的 v3 package-level 函式：
   | 建構子 | 綁定的遮罩函式（v3） | 對應 MaskerType |
   |--------|---------------------|-----------------|
   | `NewPhone(v string)`    | `Mobile`  | `TypeMobile` |
   | `NewEmail(v string)`    | `Email`   | `TypeEmail` |
   | `NewPassword(v string)` | `Password`| `TypePassword` |
   | `NewID(v string)`       | `ID`      | `TypeID` |
   | `NewCredit(v string)`   | `Credit`  | `TypeCredit` |
   | `NewName(v string)`     | `Name`    | `TypeName` |
   | `NewAddress(v string)`  | `Address` | `TypeAddress` |
   | `NewTel(v string)`      | `Tel`     | `TypeTel` |
   | `NewURL(v string)`      | `URL`     | `TypeURL` |
   - 注意命名對應：`NewPhone` → `Mobile`（手機），`NewTel` → `Tel`（市話），兩者不同，不可混用。
   - 每個建構子須有繁中 GoDoc，第一句以識別字名稱開頭，附使用範例。

5. **`Equal` 安全比較（subtask: equality）** — 在 `v3/sensitive.go`：
   ```go
   func (s Sensitive[T]) Equal(other Sensitive[T]) bool
   ```
   - 比較兩個 `Sensitive[T]` 的**原始值**是否相等，過程不暴露原值（method 對同型別可直接存取 `other.raw`）。
   - `T` 為任意型別，使用 `reflect.DeepEqual(s.raw, other.raw)` 比較（避免限制 `T comparable`，也支援 slice/struct 等）。
   - 只比 `raw`，不比 `masked` / `mask`（不同 raw 可能遮罩成相同字串，比 masked 會誤判）。

6. **`UnmarshalJSON` 支援（subtask: unmarshal）** — 在 `v3/sensitive.go`：
   ```go
   func (s *Sensitive[T]) UnmarshalJSON(data []byte) error
   ```
   - pointer receiver（需寫回 receiver）。
   - 把 `data` 解碼進 `s.raw`（`json.Unmarshal(data, &s.raw)`）。
   - 解碼後用 receiver 既有的 `s.mask` 重算 `s.masked`，使 round-trip 後 masked 與原值一致。
   - **設計前提**：UnmarshalJSON 依賴 receiver 已綁定 `mask`。常見用法是先用建構子預填 struct 欄位（綁定 maskFn），再 `json.Unmarshal` 進該 struct，欄位的既有 `mask` 會被保留並用於重算。AC/測試須涵蓋此 round-trip 流程。

7. **`mask` 為 nil 的 edge case** — 規範並測試：當 `Sensitive[T]` 為 zero value（未經建構子，`mask == nil`）時：
   - `NewSensitive` 收到 nil maskFn：masked 設為空字串 `""`（安全預設，絕不回傳原值）。
   - `UnmarshalJSON` 遇 `s.mask == nil`：解碼 raw 後，masked 設為 `""`（不得直接把 raw 當 masked 輸出而洩漏）。
   - 此規範確保「未刻意綁定遮罩」的退化情況下，安全輸出路徑仍不洩漏原值。

8. **測試** — 新增 `v3/sensitive_test.go` 與 `v3/sensitive_constructors_test.go`，含 `ExampleNewPhone` 等可執行 example。詳見 acceptance-criteria.md 與 test-strategy.yaml。

## Scope (which files/dirs to modify)

- 新增 `v3/sensitive.go` — `Sensitive[T]` struct、`Reveal`、`NewSensitive`、五個安全輸出 interface、`Equal`、`UnmarshalJSON`、編譯期介面斷言。
- 新增 `v3/sensitive_constructors.go` — 九個內建建構子。
- 新增 `v3/sensitive_test.go`、`v3/sensitive_constructors_test.go` — 測試與 examples。
- 僅在 `v3/` 內作業。

## Out of Scope

- 不修改 v2（repo root）任何檔案。
- 不修改 F001 既有 v3 檔案（`masker.go`、`convenience.go`、各 masker 檔）——只新增檔案、複用既有 package-level 函式。
- 不改既有遮罩演算法 / 輸出格式（masked 結果須與既有 v3 函式逐字一致）。
- 不實作 struct tag 整合（`Sensitive[T]` 作為 struct 欄位時走自身 Marshaler interface，不經 `mask:""` tag）。
- 不支援 `Abuse` / `None` / `All` / `Mask` 等動態型別的專屬建構子（僅九個內建建構子 + 通用 `NewSensitive`）。
- 不處理並行可變狀態（`Sensitive[T]` 為 immutable value type，建構後不改 raw）。
