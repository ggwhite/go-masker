package masker

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// fieldMask 描述單一 struct 欄位在遮罩時所需的解析後 metadata。
// 由 buildTypeInfo 對該 type 的欄位只迭代一次後產出並快取，
// 之後 Struct() 與 Format() 直接讀此結構，不再於每次呼叫重新 parse tag 與分類 kind。
//
// 未匯出欄位（exported == false）的處理規則為「略過」：Struct() 不複製、保持 zero value，
// 與既有行為一致；無 mask tag 的欄位（hasTag == false）則「直接複製原值」。
type fieldMask struct {
	index    int          // 欄位在 struct 中的宣告順序 index
	exported bool         // 欄位是否為 exported；未匯出欄位 Struct() 略過保持 zero value
	hasTag   bool         // 欄位是否帶有非空的 mask tag
	tag      MaskerType   // 解析後的 mask tag 值；hasTag == false 時無意義
	kind     reflect.Kind // 欄位型別的 kind
	elemKind reflect.Kind // 針對 Slice／Ptr 欄位記錄 element／pointee 的 kind，其餘為 reflect.Invalid
}

// typeInfo 持有單一 struct type 的全部欄位 metadata，依宣告順序排列。
type typeInfo struct {
	fields []fieldMask
}

// typeCache 快取已解析的 struct type metadata。
// key 為 reflect.Type、value 為 *typeInfo，並以 LoadOrStore 處理 cache miss 競態，concurrency-safe。
var typeCache sync.Map

// buildTypeInfoCount 記錄 buildTypeInfo 實際執行次數，供同 package 測試驗證「同 type 只解析一次」。
// 以 atomic 操作避免與並行 -race 測試衝突。
var buildTypeInfoCount int64

// cachedTypeInfo 回傳 t 的欄位 metadata；cache miss 時呼叫 buildTypeInfo 解析一次並寫回，cache hit 直接回傳。
func cachedTypeInfo(t reflect.Type) *typeInfo {
	if v, ok := typeCache.Load(t); ok {
		return v.(*typeInfo)
	}
	info := buildTypeInfo(t)
	actual, _ := typeCache.LoadOrStore(t, info)
	return actual.(*typeInfo)
}

// buildTypeInfo 對 t 的欄位只迭代一次，解析 mask tag、分類 kind、標記匯出狀態，產出 typeInfo。
// 空字串 tag（mask:""）視為無 tag，與既有 Struct() 對 len(tag)==0 的處理一致。
func buildTypeInfo(t reflect.Type) *typeInfo {
	atomic.AddInt64(&buildTypeInfoCount, 1)
	n := t.NumField()
	fields := make([]fieldMask, 0, n)
	for i := 0; i < n; i++ {
		sf := t.Field(i)
		fm := fieldMask{
			index:    i,
			exported: sf.IsExported(),
			kind:     sf.Type.Kind(),
			elemKind: reflect.Invalid,
		}
		if mtag := sf.Tag.Get(tagName); len(mtag) > 0 {
			fm.hasTag = true
			fm.tag = MaskerType(mtag)
		}
		switch fm.kind {
		case reflect.Slice, reflect.Ptr:
			fm.elemKind = sf.Type.Elem().Kind()
		}
		fields = append(fields, fm)
	}
	return &typeInfo{fields: fields}
}

// resetTypeCache 清空 type metadata cache 與解析計數，僅供同 package 測試／benchmark 使用，不對外匯出。
func resetTypeCache() {
	typeCache = sync.Map{}
	atomic.StoreInt64(&buildTypeInfoCount, 0)
}
