package masker

import (
	"encoding"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
)

// 編譯期斷言：value type 必須滿足所有安全輸出介面，避免日後改動破壞自動遮罩路徑。
var (
	_ fmt.Stringer           = Sensitive[string]{}
	_ fmt.GoStringer         = Sensitive[string]{}
	_ json.Marshaler         = Sensitive[string]{}
	_ encoding.TextMarshaler = Sensitive[string]{}
	_ slog.LogValuer         = Sensitive[string]{}
)

// Sensitive 是泛型安全型別，把遮罩從「開發者記得做」變成「洩漏原值要刻意呼叫 Reveal」。
// 任何自動輸出路徑（fmt 列印、%#v、json.Marshal、encoding.TextMarshaler、slog）一律輸出遮罩值，
// 原始值只有透過唯一出口 Reveal 才能取得。建構時就算好 masked 並快取，之後取值零額外成本。
type Sensitive[T any] struct {
	raw    T              // 唯一原值來源
	masked string         // 建構時算好並快取的遮罩字串
	mask   func(T) string // 綁定的遮罩函式，供 UnmarshalJSON 重算 masked 用
}

// NewSensitive 建立綁定指定遮罩函式的 Sensitive[T]，是所有內建建構子的共同底層。
// 建構時立即執行 maskFn(raw) 算出 masked 並快取，同時保存 maskFn 供 UnmarshalJSON 重算用。
// maskFn 為 nil 時 masked 設為空字串（安全預設，絕不回傳原值）。
// Example:
//
//	s := masker.NewSensitive("0987654321", masker.Mobile)
//	fmt.Println(s) // 0987***321
func NewSensitive[T any](raw T, maskFn func(T) string) Sensitive[T] {
	masked := ""
	if maskFn != nil {
		masked = maskFn(raw)
	}
	return Sensitive[T]{raw: raw, masked: masked, mask: maskFn}
}

// Reveal 回傳原始值，是取得原值的唯一途徑。呼叫它即代表刻意洩漏原值。
func (s Sensitive[T]) Reveal() T {
	return s.raw
}

// String 實作 fmt.Stringer，回傳遮罩值（絕不碰原值）。
func (s Sensitive[T]) String() string {
	return s.masked
}

// GoString 實作 fmt.GoStringer，回傳遮罩值，防止 %#v 洩漏 struct 內部。
func (s Sensitive[T]) GoString() string {
	return s.masked
}

// MarshalJSON 實作 encoding/json.Marshaler，把遮罩值做 JSON 編碼後輸出。
func (s Sensitive[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.masked)
}

// MarshalText 實作 encoding.TextMarshaler，回傳遮罩值的 byte 切片。
func (s Sensitive[T]) MarshalText() ([]byte, error) {
	return []byte(s.masked), nil
}

// LogValue 實作 log/slog.LogValuer，讓結構化日誌輸出遮罩值。
func (s Sensitive[T]) LogValue() slog.Value {
	return slog.StringValue(s.masked)
}

// Equal 比較兩個 Sensitive[T] 的原始值是否相等，過程不暴露原值。
// 只比 raw（不比 masked／mask），因不同 raw 可能遮罩成相同字串。
// 使用 reflect.DeepEqual 以支援 slice／struct 等非 comparable 的 T。
func (s Sensitive[T]) Equal(other Sensitive[T]) bool {
	return reflect.DeepEqual(s.raw, other.raw)
}

// UnmarshalJSON 實作 encoding/json.Unmarshaler，把 data 解碼進原值並用既有 mask 重算遮罩值。
// 依賴 receiver 已綁定 mask：常見用法是先用建構子預填欄位（綁定 maskFn），再 json.Unmarshal 進該 struct，
// 欄位既有的 mask 會被保留並用於重算。mask 為 nil 時回傳 error，避免靜默產生空遮罩值。
func (s *Sensitive[T]) UnmarshalJSON(data []byte) error {
	if s.mask == nil {
		return fmt.Errorf("masker: cannot unmarshal into Sensitive[T] without a bound mask function; use a constructor (e.g. NewPhone) first")
	}
	if err := json.Unmarshal(data, &s.raw); err != nil {
		return err
	}
	s.masked = s.mask(s.raw)
	return nil
}
