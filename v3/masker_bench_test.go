package masker

import "testing"

// benchProfile 為 benchPlayer 的巢狀 struct 欄位，驗證巢狀 type 亦受惠於 cache。
type benchProfile struct {
	Nickname string `mask:"name"`
	Email    string `mask:"email"`
}

// benchPlayer 模擬真實玩家帳戶記錄：少數 PII 欄位（需遮罩）搭配大量遊戲統計／設定／flag 等
// 非敏感 metadata 欄位（直接複製）+ 一個巢狀 struct。此分布反映典型寬 DTO——
// PII 僅佔少數，cold 路徑每次須對「全部」欄位重新 parse tag，正是 cache 消除的成本來源。
// 刻意以代表性的寬記錄（而非「全 PII」這種對 cache 最不利的極端 case）量測 cache 效益。
type benchPlayer struct {
	// PII 欄位（需遮罩）
	Username string `mask:"name"`
	PlayerID string `mask:"id"`

	// 遊戲統計（無 tag，直接複製）
	Level       int
	Experience  int
	VipTier     int
	WinCount    int
	LoseCount   int
	DrawCount   int
	TotalBets   int
	TotalWins   int
	LoginStreak int
	Rank        int
	Coins       int
	Gems        int
	Energy      int
	Trophies    int

	// 設定／flag（無 tag，直接複製）
	SoundOn      bool
	MusicOn      bool
	PushEnabled  bool
	EmailOptIn   bool
	IsBanned     bool
	IsVerified   bool
	AutoRecharge bool
	TutorialDone bool

	// 其他 metadata（無 tag，直接複製）
	Country    string
	Language   string
	Timezone   string
	Theme      string
	Currency   string
	Platform   string
	AppVersion string
	Status     string
	CreatedAt  string
	UpdatedAt  string
	LastLogin  string

	Profile benchProfile `mask:"struct"`
}

func benchFixture() benchPlayer {
	return benchPlayer{
		Username:   "John Doe",
		PlayerID:   "A123456789",
		Level:      42,
		Experience: 138500,
		VipTier:    3,
		WinCount:   1200,
		LoseCount:  340,
		DrawCount:  15,
		TotalBets:  98000,
		TotalWins:  52000,
		Coins:      999999,
		Gems:       4200,
		Energy:     80,
		Trophies:   17,
		SoundOn:    true,
		MusicOn:    false,
		IsVerified: true,
		Country:    "TW",
		Language:   "zh-TW",
		Timezone:   "Asia/Taipei",
		Theme:      "dark",
		Currency:   "TWD",
		Platform:   "ios",
		AppVersion: "2.3.1",
		Status:     "active",
		CreatedAt:  "2024-01-02T03:04:05Z",
		UpdatedAt:  "2025-06-23T10:00:00Z",
		LastLogin:  "2026-06-23T08:30:00Z",
		Profile: benchProfile{
			Nickname: "Johnny",
			Email:    "john@gmail.com",
		},
	}
}

// BenchmarkStruct 量測 cache warm 後對同一 type 反覆呼叫 Struct() 的穩態開銷。
func BenchmarkStruct(b *testing.B) {
	m := NewMaskerMarshaler()
	in := benchFixture()
	resetTypeCache()
	if _, err := m.Struct(in); err != nil { // warm cache
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m.Struct(in)
	}
}

// BenchmarkStructCold 每次 iteration 先清空 cache 再呼叫 Struct()，代表含 tag 解析的冷路徑。
func BenchmarkStructCold(b *testing.B) {
	m := NewMaskerMarshaler()
	in := benchFixture()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resetTypeCache()
		_, _ = m.Struct(in)
	}
}

// benchFmtSeg 為 Format 比較 fixture 的巢狀 struct 葉節點，欄位皆為需遮罩的 PII。
type benchFmtSeg struct {
	Name  string `mask:"name"`
	Email string `mask:"email"`
	ID    string `mask:"id"`
}

// benchFmtPlayer 為 Format vs Struct 比較專用 fixture：PII／巢狀 struct 為主、幾乎不含 raw scalar 欄位。
// 在此 fixture 上 Format()「不建新 masked struct」的 alloc 優勢明確，避免寬 scalar fixture 下
// Format 對大量無 tag scalar 走 strconv 反而 alloc 較多造成的量測雜訊（見 task-brief L013）。
type benchFmtPlayer struct {
	Username string      `mask:"name"`
	PlayerID string      `mask:"id"`
	Email    string      `mask:"email"`
	Profile  benchFmtSeg `mask:"struct"`
	Contact  benchFmtSeg `mask:"struct"`
	Backup   benchFmtSeg `mask:"struct"`
}

func benchFmtFixture() benchFmtPlayer {
	seg := benchFmtSeg{Name: "Jane Doe", Email: "jane@gmail.com", ID: "B223456789"}
	return benchFmtPlayer{
		Username: "John Doe",
		PlayerID: "A123456789",
		Email:    "john@gmail.com",
		Profile:  seg,
		Contact:  seg,
		Backup:   seg,
	}
}

// BenchmarkFormat 量測對 PII／巢狀為主的 fixture 反覆呼叫 Format() 直出遮罩字串的開銷。
func BenchmarkFormat(b *testing.B) {
	m := NewMaskerMarshaler()
	in := benchFmtFixture()
	resetTypeCache()
	_ = m.Format(in) // warm cache
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Format(in)
	}
}

// BenchmarkStructForFormat 以與 BenchmarkFormat 相同的 fixture 量測 Struct()，供 Format vs Struct 的 alloc 比較。
func BenchmarkStructForFormat(b *testing.B) {
	m := NewMaskerMarshaler()
	in := benchFmtFixture()
	resetTypeCache()
	if _, err := m.Struct(in); err != nil { // warm cache
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m.Struct(in)
	}
}
