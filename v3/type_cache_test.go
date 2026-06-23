package masker

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
)

type cacheFoo struct {
	Name  string `mask:"name"`
	Email string `mask:"email"`
	ID    string `mask:"id"`
}

type cacheBar struct {
	Account string `mask:"name"`
	Secret  string `mask:"password"`
}

// TestBuildTypeInfoOnce 驗證 AC-3：清空 cache 後對同一 type 多次呼叫 Struct()，metadata 僅建立一次。
func TestBuildTypeInfoOnce(t *testing.T) {
	m := NewMaskerMarshaler()
	resetTypeCache()

	in := cacheFoo{Name: "John Doe", Email: "john@gmail.com", ID: "A123456789"}
	const n = 20
	for i := 0; i < n; i++ {
		if _, err := m.Struct(in); err != nil {
			t.Fatalf("Struct() unexpected error: %v", err)
		}
	}

	if got := atomic.LoadInt64(&buildTypeInfoCount); got != 1 {
		t.Errorf("buildTypeInfoCount = %d, want 1 (metadata should be parsed once)", got)
	}
}

// TestCachedTypeInfoReuse 驗證 AC-2：cachedTypeInfo 對同一 type 回傳同一個 *typeInfo 實例。
func TestCachedTypeInfoReuse(t *testing.T) {
	resetTypeCache()
	t1 := cachedTypeInfo(reflect.TypeOf(cacheFoo{}))
	t2 := cachedTypeInfo(reflect.TypeOf(cacheFoo{}))
	if t1 != t2 {
		t.Errorf("cachedTypeInfo returned different instances for same type")
	}
	if len(t1.fields) != 3 {
		t.Errorf("typeInfo.fields len = %d, want 3", len(t1.fields))
	}
}

// TestStructCacheConcurrent 驗證 AC-4：多 goroutine 並行對相同與不同 type 呼叫 Struct()，輸出正確且無 data race。
func TestStructCacheConcurrent(t *testing.T) {
	m := NewMaskerMarshaler()
	resetTypeCache()

	foo := cacheFoo{Name: "John Doe", Email: "john@gmail.com", ID: "A123456789"}
	bar := cacheBar{Account: "Jane Doe", Secret: "supersecret"}

	wantFoo, err := m.Struct(foo)
	if err != nil {
		t.Fatalf("baseline Struct(foo) error: %v", err)
	}
	wantBar, err := m.Struct(bar)
	if err != nil {
		t.Fatalf("baseline Struct(bar) error: %v", err)
	}

	const goroutines = 64
	var wg sync.WaitGroup
	errs := make(chan error, goroutines)
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				if g%2 == 0 {
					got, err := m.Struct(foo)
					if err != nil {
						errs <- err
						return
					}
					if !reflect.DeepEqual(got, wantFoo) {
						errs <- fmt.Errorf("foo mismatch: got %v want %v", got, wantFoo)
						return
					}
				} else {
					got, err := m.Struct(bar)
					if err != nil {
						errs <- err
						return
					}
					if !reflect.DeepEqual(got, wantBar) {
						errs <- fmt.Errorf("bar mismatch: got %v want %v", got, wantBar)
						return
					}
				}
			}
		}(g)
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatalf("concurrent Struct() error: %v", e)
		}
	}
}
