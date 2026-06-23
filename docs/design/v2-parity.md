# v2 功能完整清單（v3 必須全部承繼）

v3 開發過程中的 parity checklist。任何 feature 完成時都應對照此表確認沒有遺漏。

## Masker Types（14 種）

| Tag | 來源 | 說明 |
|-----|------|------|
| `none` | v1 | 不遮罩 |
| `password` | v1 | 固定 14 個 mask char |
| `name` | v1 | 遮罩每個字的中間字元 |
| `addr` | v1 | 遮罩末 6 字元 |
| `email` | v1 | 保留前 3 字元 + domain |
| `mobile` | v1 | 遮罩第 4–6 位 |
| `tel` | v1 | 格式化 `(XX)XXXX-****` |
| `id` | v1 | 遮罩第 7–10 位 |
| `credit` | v1 | 遮罩第 7–12 位 |
| `url` | PR #17 | 遮罩 URL password（`url.Redacted()`） |
| `abuse` | PR #32 | trie-based 敏感詞遮罩 |
| `all` | PR #36 | 全字元遮罩 |
| `first-N` | PR #36 | 動態 — 遮罩前 N 字元 |
| `last-N` | PR #36 | 動態 — 遮罩末 N 字元 |

## Structural Tags（2 種）

| Tag | 來源 | 說明 |
|-----|------|------|
| `struct` | v1 | 遞迴遮罩巢狀 struct（value 和 pointer） |
| `mapstruct` | PR #38 | 遞迴遮罩 map values（map/struct/ptr/slice/ptr-to-slice 組合） |

## Struct() 處理能力

- string fields
- nested struct（value type 和 pointer）
- `[]string` slice
- `[]struct`、`[]*struct`、`[]interface{}` slice（搭配 `struct` tag）
- map fields（搭配 `mapstruct` tag）：
  - `map[K]Struct`、`map[K]*Struct`
  - `map[K][]Struct`、`map[K]*[]Struct`
  - 巢狀 map（`map[K]map[K2]...Struct`）
- `interface{}` fields（搭配 `struct` tag）
- nil 值保留（nil pointer / nil slice / nil map）
- 非 struct 輸入回傳 error 不 panic（PR #35）
- unexported fields 跳過不處理

## API 表面

| API | 說明 |
|-----|------|
| `NewMaskerMarshaler()` | 建立含全部預設 masker 的實例 |
| `DefaultMaskerMarshaler` | 全域預設實例 |
| `Marshal(MaskerType, string)` | 單值遮罩 |
| `Struct(interface{})` | struct 遮罩 |
| `Register(MaskerType, Masker)` | 註冊/覆蓋 masker |
| `Unregister(MaskerType)` | 移除 masker |
| `Get(MaskerType)` | 取得 masker |
| `List()` | 列出所有 masker type |
| `SetMasker(string)` | 自訂 mask 字元 |

## Abuse 子系統

- `AbuseMasker` — trie 實作，支援 `NewAbuseMasker()` 和 `NewAbuseMaskerWithWords()`
- `AbuseTrie` — 底層 trie 結構
- `AbuseWordLoader` — `LoadFromFile(path)` 和 `LoadFromReader(io.Reader)`

## 非功能性

- Concurrency-safe（所有 `MaskerMarshaler` 公開方法用 `sync.RWMutex`）
- 零外部依賴
- Go 1.17+（v3 升至 1.21+）
