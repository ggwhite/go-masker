package masker

// 本檔提供 package-level 便利函式，底層皆走 DefaultMaskerMarshaler（遮罩字元 '*'），
// 直接回傳 string 不回 error，方便最常見的單值遮罩場景。

// Mobile 以預設設定遮罩手機號碼。
// Example:
//
//	masker.Mobile("0987654321") // returns "0987***321"
func Mobile(value string) string {
	return DefaultMaskerMarshaler.MustMarshal(TypeMobile, value)
}

// Email 以預設設定遮罩 email。
// Example:
//
//	masker.Email("ggw.chang@gmail.com") // returns "ggw****@gmail.com"
func Email(value string) string {
	return DefaultMaskerMarshaler.MustMarshal(TypeEmail, value)
}

// Password 以預設設定遮罩密碼，固定回傳 14 個遮罩字元。
// Example:
//
//	masker.Password("password") // returns "**************"
func Password(value string) string {
	return DefaultMaskerMarshaler.MustMarshal(TypePassword, value)
}

// Name 以預設設定遮罩姓名。
// Example:
//
//	masker.Name("John Doe") // returns "J**n D**e"
func Name(value string) string {
	return DefaultMaskerMarshaler.MustMarshal(TypeName, value)
}

// Address 以預設設定遮罩地址。
// Example:
//
//	masker.Address("台北市內湖區內湖路一段737巷1號1樓") // returns "台北市內湖區******"
func Address(value string) string {
	return DefaultMaskerMarshaler.MustMarshal(TypeAddress, value)
}

// ID 以預設設定遮罩身分證、護照等證號。
// Example:
//
//	masker.ID("A123456789") // returns "A12345****"
func ID(value string) string {
	return DefaultMaskerMarshaler.MustMarshal(TypeID, value)
}

// Credit 以預設設定遮罩信用卡號。
// Example:
//
//	masker.Credit("4111111111111111") // returns "411111******1111"
func Credit(value string) string {
	return DefaultMaskerMarshaler.MustMarshal(TypeCredit, value)
}

// Tel 以預設設定遮罩市話。
// Example:
//
//	masker.Tel("0227993078") // returns "(02)2799-****"
func Tel(value string) string {
	return DefaultMaskerMarshaler.MustMarshal(TypeTel, value)
}

// URL 以預設設定遮罩 URL 的密碼段。
// Example:
//
//	masker.URL("http://john:password@localhost:3000") // returns "http://john:xxxxx@localhost:3000"
func URL(value string) string {
	return DefaultMaskerMarshaler.MustMarshal(TypeURL, value)
}

// Abuse 以預設設定遮罩敏感詞。
// DefaultMaskerMarshaler 的 AbuseMasker 預設為空字典，故未載入詞典時原樣回傳。
// Example:
//
//	masker.Abuse("hello") // returns "hello"（空字典）
func Abuse(value string) string {
	return DefaultMaskerMarshaler.MustMarshal(TypeAbuse, value)
}

// None 原樣回傳輸入值，不做任何遮罩。
// Example:
//
//	masker.None("secret") // returns "secret"
func None(value string) string {
	return DefaultMaskerMarshaler.MustMarshal(TypeNone, value)
}

// All 以預設設定把每個字元都換成遮罩字元。
// Example:
//
//	masker.All("secret") // returns "******"
func All(value string) string {
	return DefaultMaskerMarshaler.MustMarshal(TypeAll, value)
}

// Mask 以動態型別遮罩值，底層走 DefaultMaskerMarshaler。
// 找不到型別時回傳原值（與 MustMarshal 的 panic 行為不同，提供寬鬆的單值入口）。
// 同時支援動態 tag first-N / last-N。
// Example:
//
//	masker.Mask(masker.TypeMobile, "0987654321") // returns "0987***321"
//	masker.Mask("nonexistent", "value")          // returns "value"
func Mask(t MaskerType, value string) string {
	result, err := DefaultMaskerMarshaler.Marshal(t, value)
	if err != nil {
		return value
	}
	return result
}
