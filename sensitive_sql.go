package masker

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
)

// 編譯期斷言：Sensitive[T] 必須能直接被 database/sql 當作參數寫入與欄位掃描，
// 讓 GORM 等 ORM 直接宣告敏感欄位即可自動 scan/value，不需手動轉換。
var (
	_ driver.Valuer = Sensitive[string]{}
	_ sql.Scanner   = (*Sensitive[string])(nil)
)

// Value 實作 database/sql/driver.Valuer，回傳原始值供 ORM 寫入 DB（存明文，不是遮罩值）。
// 透過 driver.DefaultParameterConverter 把 raw 轉成合法的 driver.Value（如 int→int64），
// 若 T 為 driver 無法處理的型別則回傳 error。
func (s Sensitive[T]) Value() (driver.Value, error) {
	return driver.DefaultParameterConverter.ConvertValue(s.raw)
}

// Scan 實作 database/sql.Scanner，把 DB 欄位值寫進原值並用既有 mask 重算遮罩值。
// src 為 nil（DB NULL）時把原值設為 zero value，不 panic。
// 與 UnmarshalJSON 不同，Scan 容忍未綁定 mask 的 receiver（ORM 常以 zero value struct 接收查詢結果）：
// 未綁定時仍安全儲存原值，masked 維持空字串而非洩漏原值；綁定後（建構子預填）則重算遮罩值。
func (s *Sensitive[T]) Scan(src any) error {
	if src == nil {
		var zero T
		s.raw = zero
		s.masked = s.maskRaw()
		return nil
	}
	if err := convertAssign(&s.raw, src); err != nil {
		return err
	}
	s.masked = s.maskRaw()
	return nil
}

// maskRaw 用綁定的 mask 算遮罩值；mask 為 nil 時回傳空字串（安全預設，絕不回傳原值）。
func (s *Sensitive[T]) maskRaw() string {
	if s.mask == nil {
		return ""
	}
	return s.mask(s.raw)
}

// convertAssign 把 DB 驅動回傳的 src（int64/float64/bool/[]byte/string/time.Time 等）轉進 *dst。
// 設計仿 database/sql 的 convertAssign：先試直接賦值，再依目的型別做字串／數值轉換，
// 涵蓋 T 為 string、整數、浮點、bool、[]byte 等常見 DB 欄位型別。
func convertAssign[T any](dst *T, src any) error {
	dv := reflect.ValueOf(dst).Elem()
	dt := dv.Type()
	sv := reflect.ValueOf(src)

	// 同型別或可直接賦值（time.Time、bool、string、相同 named type 等）。
	if sv.Type().AssignableTo(dt) {
		dv.Set(sv)
		return nil
	}

	switch dt.Kind() {
	case reflect.String:
		if str, ok := asString(src); ok {
			dv.SetString(str)
			return nil
		}
	case reflect.Slice:
		if dt.Elem().Kind() == reflect.Uint8 {
			if b, ok := asBytes(src); ok {
				cp := make([]byte, len(b))
				copy(cp, b)
				dv.SetBytes(cp)
				return nil
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if n, ok := asInt64(src); ok {
			if dv.OverflowInt(n) {
				return fmt.Errorf("masker: scan value %v overflows Sensitive[%s]", src, dt)
			}
			dv.SetInt(n)
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if n, ok := asUint64(src); ok {
			if dv.OverflowUint(n) {
				return fmt.Errorf("masker: scan value %v overflows Sensitive[%s]", src, dt)
			}
			dv.SetUint(n)
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if f, ok := asFloat64(src); ok {
			if dv.OverflowFloat(f) {
				return fmt.Errorf("masker: scan value %v overflows Sensitive[%s]", src, dt)
			}
			dv.SetFloat(f)
			return nil
		}
	case reflect.Bool:
		if b, ok := asBool(src); ok {
			dv.SetBool(b)
			return nil
		}
	}
	return fmt.Errorf("masker: cannot scan %T into Sensitive[%s]", src, dt)
}

// asString 把 DB 來源轉成字串，支援 string 與 []byte（多數驅動回傳 []byte）。
func asString(src any) (string, bool) {
	switch v := src.(type) {
	case string:
		return v, true
	case []byte:
		return string(v), true
	}
	return "", false
}

// asBytes 取得 DB 來源的位元組形式，供 T 為 []byte 的欄位使用。
func asBytes(src any) ([]byte, bool) {
	switch v := src.(type) {
	case []byte:
		return v, true
	case string:
		return []byte(v), true
	}
	return nil, false
}

// asInt64 把 DB 來源轉成 int64，涵蓋驅動可能回傳的整數、bool 與字串數字。
func asInt64(src any) (int64, bool) {
	switch v := src.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case int32:
		return int64(v), true
	case int16:
		return int64(v), true
	case int8:
		return int64(v), true
	case bool:
		if v {
			return 1, true
		}
		return 0, true
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 64)
		return n, err == nil
	case string:
		n, err := strconv.ParseInt(v, 10, 64)
		return n, err == nil
	}
	return 0, false
}

// asUint64 把 DB 來源轉成 uint64，負數來源視為無法轉換。
func asUint64(src any) (uint64, bool) {
	switch v := src.(type) {
	case int64:
		if v < 0 {
			return 0, false
		}
		return uint64(v), true
	case int:
		if v < 0 {
			return 0, false
		}
		return uint64(v), true
	case bool:
		if v {
			return 1, true
		}
		return 0, true
	case []byte:
		n, err := strconv.ParseUint(string(v), 10, 64)
		return n, err == nil
	case string:
		n, err := strconv.ParseUint(v, 10, 64)
		return n, err == nil
	}
	return 0, false
}

// asFloat64 把 DB 來源轉成 float64，涵蓋浮點、整數與字串數字（如 MySQL DECIMAL 回傳 []byte）。
func asFloat64(src any) (float64, bool) {
	switch v := src.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int64:
		return float64(v), true
	case int:
		return float64(v), true
	case []byte:
		f, err := strconv.ParseFloat(string(v), 64)
		return f, err == nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		return f, err == nil
	}
	return 0, false
}

// asBool 把 DB 來源轉成 bool，支援 bool、0/1 整數與字串（透過 strconv.ParseBool）。
func asBool(src any) (bool, bool) {
	switch v := src.(type) {
	case bool:
		return v, true
	case int64:
		return v != 0, true
	case int:
		return v != 0, true
	case []byte:
		b, err := strconv.ParseBool(string(v))
		return b, err == nil
	case string:
		b, err := strconv.ParseBool(v)
		return b, err == nil
	}
	return false, false
}
