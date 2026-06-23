# Acceptance Criteria

| # | Criterion | Verification Method |
|---|---|---|
| AC-1 | `Struct()` 對外行為與現況完全一致：所有既有 `TestMaskerMarshaler_Struct`、`TestMaskerMarshaler_MapStructs`、`TestStruct_*` 測試在不修改斷言的情況下全數通過。 | `go test -run 'Struct\|MapStruct' ./...` 全綠；diff 確認既有測試斷言未被更動。 |
| AC-2 | 新增 type metadata cache：`typeCache`（`sync.Map`，key=`reflect.Type`，value=`*typeInfo`），`buildTypeInfo` 對單一 type 的欄位只迭代解析一次。 | 程式碼審查確認結構存在；AC-3 的測試佐證只建立一次。 |
| AC-3 | 同一 type 的 metadata 只解析一次：清空 cache 後對同型別多次呼叫 `Struct()`，metadata 僅建立 1 次（透過同 package 可觀測的內部計數或等效機制驗證）。 | `package masker` 內測試：`resetTypeCache()` 後連續呼叫 `Struct()` N 次，斷言解析次數 == 1。 |
| AC-4 | Cache 為 concurrency-safe：多 goroutine 並行對「相同」與「不同」type 呼叫 `Struct()`，無 data race 且輸出正確。 | `go test -race`：啟 ≥50 goroutine 並行呼叫 `Struct()`，`-race` 無告警，輸出與單執行緒一致。 |
| AC-5 | 新增 `func (m *MaskerMarshaler) Format(s interface{}) string`，簽章為 `*MaskerMarshaler` method。 | 編譯通過；`go doc` 顯示該方法含繁中 GoDoc。 |
| AC-6 | `Format()` 行為約束：nil input 回 `""` 且不 panic；任何 input 皆不 panic（含 nil pointer、空 struct）。 | 表格測試涵蓋 nil、`(*T)(nil)`、空 struct、含 nil pointer 欄位；以 `recover` 斷言無 panic。 |
| AC-7 | `Format()` 遮罩 parity：對同一 input，`Format()` 輸出中每個遮罩後的 string 欄位值，與 `Struct()` 結果對應欄位 byte-identical。 | 測試以已知 fixture，逐欄比對 `Format()` 子字串與 `Struct()` 結果欄位值相等。 |
| AC-8 | `Format()` 輸出為確定性：不含記憶體位址，多次呼叫結果相同；ptr-to-struct 遞迴 inline 展開（遮罩後），nil pointer 顯示 `<nil>`。 | 測試對含 ptr 欄位的 fixture 連呼叫 2 次斷言相等，且 `strings.Contains(out, "0x") == false`，巢狀內容已展開。 |
| AC-9 | `Format()` 配置數低於 `Struct()`（不建新 struct）：`testing.AllocsPerRun` 量得 `Format()` 的 allocs/op < `Struct()` 的 allocs/op。 | `testing.AllocsPerRun(100, ...)` 比較兩者，斷言 `formatAllocs < structAllocs`。 |
| AC-10 | 同 type 第二次以上呼叫 `Struct()` CPU 開銷降低 ≥ 50%（cache warm vs cold）。 | `go test -bench` 取 `BenchmarkStructCold` 與 `BenchmarkStruct`（warm）的 ns/op，warm ≤ cold 的 50%；數據貼入報告。 |
| AC-11 | Benchmark 齊備：`BenchmarkStruct`（warm）、`BenchmarkStructCold`、`BenchmarkFormat` 三者存在並可執行。 | `go test -bench='BenchmarkStruct\|BenchmarkStructCold\|BenchmarkFormat' -benchmem -run='^$' ./...` 正常輸出三項數據。 |
| AC-12 | 程式碼品質：`go build ./...`、`go vet ./...` 通過；新增 exported 識別字（`Format`）含繁中 GoDoc，第一句以識別字名稱開頭。 | `go build ./...` 與 `go vet ./...` 無錯誤；審查 GoDoc。 |
