# go-masker v3 設計文件

## 定位

v3 從「遮罩工具函式庫」升級為「Go 敏感資料防護層」。

核心轉變：把遮罩從「開發者要記得做」變成「要洩漏反而要刻意寫 `.Reveal()`」。

## 為什麼需要 v3

v1 有 package-level 便利函式（`masker.Mobile(phone)`），用起來直覺。
v2 為了可擴充性砍了這些，只留 `MaskerMarshaler.Marshal(MaskerType, string)`——變得更通用但太囉嗦，導致使用者寧願自己寫 helper。

v2 的 `Masker` interface 把 mask char 綁進方法簽名（`Marshal(s, i string) string`），無法在不破壞 API 的前提下修正。

v3 的目標：**v1 的簡潔 + v2 的擴充性 + 型別層安全**。

## 設計原則

1. **零門檻上手** — `masker.Mobile(phone)` 一行搞定
2. **預設安全** — `Sensitive[T]` 讓洩漏變成顯式選擇
3. **零核心依賴** — core package 不依賴 zap/logrus 等外部套件
4. **效能友善** — reflect cache、避免不必要的 alloc
5. **向後可遷移** — v2 使用者有明確的 migration path

## Module 與 Go 版本

```
github.com/ggwhite/go-masker/v3
```

最低 Go 版本：**1.21**
- generics（Sensitive[T]）需要 1.18+
- `log/slog` 需要 1.21+
- 1.21 發布已兩年，合理的最低版本

## Package 結構

```
go-masker/v3/                  ← core module（零外部依賴）
├── masker.go                  ← Masker interface, MaskerMarshaler, Register/Get
├── convenience.go             ← Package-level 短函式：Mobile(), Email(), Password()...
├── sensitive.go               ← Sensitive[T] 泛型安全型別
├── struct.go                  ← Struct() + reflect type cache
├── format.go                  ← Format() 直出字串（不 alloc 新 struct）
├── generic.go                 ← AllMasker, first-N, last-N
├── <type>.go                  ← 各類型 masker（mobile, email, password...）
├── go.mod
│
go-masker/v3/zapfield/         ← 獨立 module（依賴 zap）
├── field.go                   ← zap Field helpers + Sensitive[T] adapter
├── interceptor.go             ← zap Core wrapper（P2）
├── go.mod                     ← depends on go-masker/v3 + go.uber.org/zap
```

`slog` 整合直接放 core（`log/slog` 是 stdlib），不需要獨立 module。

## API 設計

### 1. Masker Interface（新）

```go
// v2（要改掉）
type Masker interface {
    Marshal(maskChar, value string) string
}

// v3
type Masker interface {
    Mask(value string) string
}
```

mask char 從 interface 移到配置層：

```go
type Option func(*config)

func WithMaskChar(c rune) Option { ... }

func NewMaskerMarshaler(opts ...Option) *MaskerMarshaler { ... }
```

### 2. Package-level 便利函式

```go
// 一行搞定，零 reflection，零 error
masker.Mobile("0987654321")     // "0987***321"
masker.Email("ggw@gmail.com")   // "ggw****@gmail.com"
masker.Password("secret")       // "**************"
masker.Name("王大明")            // "王*明"
masker.Address("台北市...")      // "台北市內湖區******"
masker.ID("A123456789")         // "A12345****"
masker.Credit("4111...")        // "411111******1111"
masker.Tel("0227993078")        // "(02)2799-****"
masker.URL("http://u:p@host")   // "http://u:xxxxx@host"

// 動態 type（runtime 決定 masker 類型時用）
masker.Mask("mobile", "0987654321")  // 用 string 不用常數也行
```

底層直接呼叫各 masker 的 `Mask()`，不經 map lookup，不回 error。

### 3. Sensitive[T] — 型別層防洩漏

```go
type Sensitive[T any] struct {
    raw    T
    masked string
}

// 建構子：建立時就算好遮罩值，之後取值零成本
func NewPhone(v string) Sensitive[string]
func NewEmail(v string) Sensitive[string]
func NewPassword(v string) Sensitive[string]
func NewID(v string) Sensitive[string]
func NewCredit(v string) Sensitive[string]
// ...每個 masker type 一個建構子

// 通用建構子：用 MaskerType 或自訂 Masker
func NewSensitive[T any](raw T, maskFn func(T) string) Sensitive[T]

// 安全輸出（所有「無意識」的印出路徑都走遮罩值）
func (s Sensitive[T]) String() string           // fmt.Stringer
func (s Sensitive[T]) GoString() string         // fmt.GoStringer (%#v)
func (s Sensitive[T]) MarshalJSON() ([]byte, error)  // encoding/json
func (s Sensitive[T]) LogValue() slog.Value     // log/slog
func (s Sensitive[T]) MarshalText() ([]byte, error)  // encoding.TextMarshaler

// 顯式取原值
func (s Sensitive[T]) Reveal() T

// 比較（不暴露原值）
func (s Sensitive[T]) Equal(other Sensitive[T]) bool
```

使用場景：

```go
type User struct {
    Name    string
    Phone   masker.Sensitive[string]  `json:"phone"`
    Email   masker.Sensitive[string]  `json:"email"`
}

u := User{
    Name:  "王大明",
    Phone: masker.NewPhone("0987654321"),
    Email: masker.NewEmail("ggw@gmail.com"),
}

fmt.Println(u.Phone)           // 0987***321
json.Marshal(u)                // {"phone":"0987***321","email":"ggw****@gmail.com"}
slog.Info("user", "data", u)   // phone=0987***321 email=ggw****@gmail.com

// 要用原值：顯式呼叫
sms.Send(u.Phone.Reveal())
```

### 4. Struct() + Reflect Cache

```go
// 快取 type metadata，只在首次遇到新 type 時 reflect
type typeInfo struct {
    fields []fieldMask  // 哪些 field 要 mask、用什麼 masker
}

var typeCache sync.Map  // key: reflect.Type → value: *typeInfo
```

`Struct()` 簽名不變，但內部用 cache 避免重複解析 tag。

新增 `Format()`：直接輸出遮罩後的字串表示，不 alloc 新 struct：

```go
func (m *MaskerMarshaler) Format(s any) string
```

### 5. zapfield Sub-module

```go
import "github.com/ggwhite/go-masker/v3/zapfield"

// Field helpers — 遮罩後包成 zap.Field
zapfield.Phone("phone", "0987654321")         // zap.String("phone", "0987***321")
zapfield.Email("email", "ggw@gmail.com")      // zap.String("email", "ggw****@gmail.com")

// Sensitive[T] adapter — 直接吃 Sensitive 型別
zapfield.Sensitive("phone", user.Phone)       // zap.String("phone", "0987***321")

// P2: Core wrapper — field name keyword 攔截
core := zapfield.WrapCore(originalCore, zapfield.InterceptRules{
    Keywords: []string{"phone", "password", "token"},
})
```

## MaskerType 常數（精簡命名）

```go
// v2
masker.MaskerTypeMobile    // 太長

// v3
masker.TypeMobile          // 精簡但保留 Type prefix 避免衝突
masker.TypeEmail
masker.TypePassword
// ...
```

## Migration Path（v2 → v3）

| v2 | v3 |
|---|---|
| `import "github.com/ggwhite/go-masker/v2"` | `import "github.com/ggwhite/go-masker/v3"` |
| `masker.MaskerTypeMobile` | `masker.TypeMobile` |
| `masker.DefaultMaskerMarshaler.Marshal(masker.MaskerTypeMobile, v)` | `masker.Mobile(v)` |
| `m.Maskers[t].Marshal(maskChar, v)` | `m.Maskers[t].Mask(v)` |
| 手動遮罩再 log | `masker.NewPhone(v)` 建立後自動安全 |

### Breaking Changes 清單

1. Module path: `/v2` → `/v3`
2. `Masker` interface: `Marshal(s, i string) string` → `Mask(value string) string`
3. `MaskerType` 常數重新命名（`MaskerTypeMobile` → `TypeMobile`）
4. 最低 Go 版本: 1.17 → 1.21
5. `MaskerMarshaler.Marshal()` 簽名可能調整（移除 error return 或改為 `MustMarshal`）

## 優先級

| 優先級 | Feature | 說明 |
|--------|---------|------|
| P0 | v3 core + 新 interface + 短函式 | 基礎，所有功能依賴此 |
| P0 | Sensitive[T] | 差異化功能，v3 存在的理由 |
| P1 | Struct() reflect cache | 效能改善 |
| P1 | zapfield sub-module | Logger 整合 |
| P2 | zap Core interceptor | 進階攔截 |
| P2 | Format() 直出字串 | 優化 |

## 未來可能（不在 v3 scope）

- protobuf message masking（issue #21 的延伸）
- `database/sql.Scanner` / `driver.Valuer` 整合
- OpenTelemetry attribute masking
