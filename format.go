package masker

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Format 回傳 s 遮罩後的確定性字串表示，過程中不配置新的 masked struct，
// 適合「只需遮罩後字串、不需遮罩後物件」的記 log／debug 情境。
//
// 與 Struct() 共用同一條遮罩路徑：對每個 string 欄位的遮罩結果與 Struct() byte-identical，
// 並對 slice／map／巢狀 struct／interface 欄位套用相同遮罩語意，不直印任何敏感原值。
//
// 輸出採 Go %v struct 風格 {v0 v1 ... vn}，欄位依宣告順序以單一空白分隔；
// ptr-to-struct 以 &{...} 遞迴 inline 展開（遮罩後、不含記憶體位址），nil pointer 顯示 <nil>。
// 行為約束：nil input 回空字串、任何 input 皆不 panic、輸出為確定性（多次呼叫結果相同）。
//
// Example:
//
//	type Foo struct {
//		Name string `mask:"name"`
//		Foo  *Foo   `mask:"struct"`
//	}
//	m := masker.NewMaskerMarshaler()
//	log.Println(m.Format(&Foo{Name: "John Doe", Foo: &Foo{Name: "Jane Doe"}}))
//	// &{J**n D**e &{J**e D**e <nil>}}
func (m *MaskerMarshaler) Format(s interface{}) string {
	if s == nil {
		return ""
	}
	var b strings.Builder
	// 預先配置一塊緩衝，以單次分配取代 Builder 漸進成長造成的多次 realloc。
	b.Grow(256)
	m.writeRoot(&b, reflect.ValueOf(s))
	return b.String()
}

// writeRoot 將頂層值寫入 b：struct 套遮罩展開、ptr-to-struct 以 &{...} 展開，其餘以確定性方式輸出。
func (m *MaskerMarshaler) writeRoot(b *strings.Builder, v reflect.Value) {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			b.WriteString("<nil>")
			return
		}
		b.WriteByte('&')
		m.writeStructMasked(b, v.Elem())
	case reflect.Struct:
		m.writeStructMasked(b, v)
	default:
		writeVerbatim(b, v)
	}
}

// writeStructMasked 依 cache 的 typeInfo 將遮罩後的 struct 內容以 {f0 f1 ...} 寫入 b。
// 未匯出欄位輸出其 zero value（與 Struct() 保持 zero 的行為一致），無 tag 欄位以原值確定性輸出。
// 非 struct 值會退回 writeVerbatim，避免對 mask:"struct" 標錯型別的欄位 panic。
func (m *MaskerMarshaler) writeStructMasked(b *strings.Builder, sv reflect.Value) {
	if sv.Kind() != reflect.Struct {
		writeVerbatim(b, sv)
		return
	}
	info := cachedTypeInfo(sv.Type())
	b.WriteByte('{')
	for idx, f := range info.fields {
		if idx > 0 {
			b.WriteByte(' ')
		}
		if !f.exported {
			writeVerbatim(b, reflect.Zero(sv.Field(f.index).Type()))
			continue
		}
		fv := sv.Field(f.index)
		if !f.hasTag {
			writeVerbatim(b, fv)
			continue
		}
		// string 是最常見的遮罩欄位，內聯處理避免 reflect.Value 跨方法逃逸至 heap。
		if f.kind == reflect.String {
			if v, err := m.Marshal(f.tag, fv.String()); err == nil {
				b.WriteString(v)
			} else {
				b.WriteString(fv.String())
			}
			continue
		}
		m.writeMaskedField(b, fv, f)
	}
	b.WriteByte('}')
}

// writeMaskedField 對單一帶 tag 欄位套用與 Struct() 相同的遮罩語意，並以確定性方式寫入 b。
func (m *MaskerMarshaler) writeMaskedField(b *strings.Builder, fv reflect.Value, f fieldMask) {
	switch f.kind {
	default:
		writeVerbatim(b, fv)
	case reflect.String:
		v, err := m.Marshal(f.tag, fv.String())
		if err != nil {
			// 與 Struct() 遇 masker 不存在會回 error 對應，Format 無法回 error，
			// 退回原值以保確定性且不 panic。
			b.WriteString(fv.String())
			return
		}
		b.WriteString(v)
	case reflect.Struct:
		if f.tag == MaskerTypeStruct {
			m.writeStructMasked(b, fv)
		} else {
			writeVerbatim(b, fv)
		}
	case reflect.Ptr:
		if fv.IsNil() {
			b.WriteString("<nil>")
			return
		}
		b.WriteByte('&')
		if f.tag == MaskerTypeStruct {
			m.writeStructMasked(b, fv.Elem())
		} else {
			writeVerbatim(b, fv.Elem())
		}
	case reflect.Map:
		if fv.IsNil() {
			b.WriteString("map[]")
			return
		}
		if f.tag != MaskerTypeMapStruct {
			writeVerbatim(b, fv)
			return
		}
		m.writeMaskedMap(b, fv)
	case reflect.Slice:
		if fv.IsNil() {
			b.WriteString("[]")
			return
		}
		m.writeMaskedSlice(b, fv, f)
	case reflect.Interface:
		if fv.IsNil() {
			b.WriteString("<nil>")
			return
		}
		if f.tag != MaskerTypeStruct {
			// Struct() 對非 struct tag 的 interface 欄位不複製（保持 zero），對應 <nil>。
			b.WriteString("<nil>")
			return
		}
		m.writeStructPtrOrValue(b, fv.Elem())
	}
}

// writeMaskedMap 以 mapstruct 語意遮罩每個 value，並以 key 排序後的確定性 map[k:v ...] 寫入 b。
func (m *MaskerMarshaler) writeMaskedMap(b *strings.Builder, fv reflect.Value) {
	type entry struct{ k, v string }
	entries := make([]entry, 0, fv.Len())
	iter := fv.MapRange()
	for iter.Next() {
		maskedValue, err := m.maskMapStructValue(iter.Value())
		if err != nil {
			maskedValue = iter.Value()
		}
		var kb, vb strings.Builder
		writeVerbatim(&kb, iter.Key())
		writeVerbatim(&vb, maskedValue)
		entries = append(entries, entry{kb.String(), vb.String()})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].k < entries[j].k })
	b.WriteString("map[")
	for i, e := range entries {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(e.k)
		b.WriteByte(':')
		b.WriteString(e.v)
	}
	b.WriteByte(']')
}

// writeMaskedSlice 依 element kind 套用與 Struct() 相同的 slice 遮罩語意，輸出 [v0 v1 ...]。
func (m *MaskerMarshaler) writeMaskedSlice(b *strings.Builder, fv reflect.Value, f fieldMask) {
	b.WriteByte('[')
	switch {
	case f.elemKind == reflect.String:
		for i := 0; i < fv.Len(); i++ {
			if i > 0 {
				b.WriteByte(' ')
			}
			v, err := m.Marshal(f.tag, fv.Index(i).String())
			if err != nil {
				v = fv.Index(i).String()
			}
			b.WriteString(v)
		}
	case f.elemKind == reflect.Struct && f.tag == MaskerTypeStruct:
		for i := 0; i < fv.Len(); i++ {
			if i > 0 {
				b.WriteByte(' ')
			}
			m.writeStructMasked(b, fv.Index(i))
		}
	case f.elemKind == reflect.Ptr && f.tag == MaskerTypeStruct:
		for i := 0; i < fv.Len(); i++ {
			if i > 0 {
				b.WriteByte(' ')
			}
			el := fv.Index(i)
			if el.IsNil() {
				b.WriteString("<nil>")
				continue
			}
			b.WriteByte('&')
			m.writeStructMasked(b, el.Elem())
		}
	case f.elemKind == reflect.Interface && f.tag == MaskerTypeStruct:
		for i := 0; i < fv.Len(); i++ {
			if i > 0 {
				b.WriteByte(' ')
			}
			el := fv.Index(i)
			if el.IsNil() {
				b.WriteString("<nil>")
				continue
			}
			m.writeStructPtrOrValue(b, el.Elem())
		}
	default:
		for i := 0; i < fv.Len(); i++ {
			if i > 0 {
				b.WriteByte(' ')
			}
			writeVerbatim(b, fv.Index(i))
		}
	}
	b.WriteByte(']')
}

// writeStructPtrOrValue 將一個具體值（可能是 struct 或 ptr-to-struct）以遮罩後確定性形式寫入 b。
// 對應 Struct() 處理 interface／slice-of-interface 時依底層是否為 pointer 而展開的兩條分支。
func (m *MaskerMarshaler) writeStructPtrOrValue(b *strings.Builder, ev reflect.Value) {
	if ev.Kind() == reflect.Ptr {
		if ev.IsNil() {
			b.WriteString("<nil>")
			return
		}
		b.WriteByte('&')
		m.writeStructMasked(b, ev.Elem())
		return
	}
	m.writeStructMasked(b, ev)
}

// writeVerbatim 將 v 以確定性方式（不含記憶體位址）寫入 b，不套用任何遮罩。
// 用於無 tag 欄位、未匯出欄位 zero value 等「Struct() 直接複製」的情境：
// pointer 一律以 &<展開> 或 <nil> 呈現、struct 遞迴展開、scalar 走 fmt %v。
func writeVerbatim(b *strings.Builder, v reflect.Value) {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			b.WriteString("<nil>")
			return
		}
		b.WriteByte('&')
		writeVerbatim(b, v.Elem())
	case reflect.Interface:
		if v.IsNil() {
			b.WriteString("<nil>")
			return
		}
		writeVerbatim(b, v.Elem())
	case reflect.Struct:
		b.WriteByte('{')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				b.WriteByte(' ')
			}
			writeVerbatim(b, v.Field(i))
		}
		b.WriteByte('}')
	case reflect.Slice, reflect.Array:
		if v.Kind() == reflect.Slice && v.IsNil() {
			b.WriteString("[]")
			return
		}
		b.WriteByte('[')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				b.WriteByte(' ')
			}
			writeVerbatim(b, v.Index(i))
		}
		b.WriteByte(']')
	case reflect.Map:
		if v.IsNil() {
			b.WriteString("map[]")
			return
		}
		writeVerbatimMap(b, v)
	case reflect.String:
		b.WriteString(v.String())
	case reflect.Bool:
		b.WriteString(strconv.FormatBool(v.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b.WriteString(strconv.FormatInt(v.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		b.WriteString(strconv.FormatUint(v.Uint(), 10))
	case reflect.Float32:
		b.WriteString(strconv.FormatFloat(v.Float(), 'g', -1, 32))
	case reflect.Float64:
		b.WriteString(strconv.FormatFloat(v.Float(), 'g', -1, 64))
	default:
		// 其餘罕見 kind（complex、chan、func 等）：傳入 reflect.Value，
		// fmt 會輸出其持有值，未匯出欄位亦不會 panic。
		fmt.Fprintf(b, "%v", v)
	}
}

// writeVerbatimMap 以 key 排序輸出 map[k:v ...]，確保確定性。
func writeVerbatimMap(b *strings.Builder, v reflect.Value) {
	type entry struct{ k, v string }
	entries := make([]entry, 0, v.Len())
	iter := v.MapRange()
	for iter.Next() {
		var kb, vb strings.Builder
		writeVerbatim(&kb, iter.Key())
		writeVerbatim(&vb, iter.Value())
		entries = append(entries, entry{kb.String(), vb.String()})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].k < entries[j].k })
	b.WriteString("map[")
	for i, e := range entries {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(e.k)
		b.WriteByte(':')
		b.WriteString(e.v)
	}
	b.WriteByte(']')
}
