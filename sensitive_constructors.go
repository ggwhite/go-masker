package masker

// 本檔提供九個內建的 Sensitive[string] 建構子，各自綁定對應的 v3 package-level 遮罩函式，
// 全部委派 NewSensitive，確保 masked 輸出與既有 v3 行為逐字一致。

// NewPhone 建立綁定手機遮罩（Mobile）的 Sensitive[string]。
// 注意與 NewTel（市話）不同，勿混用。
// Example:
//
//	s := masker.NewPhone("0987654321")
//	fmt.Println(s)        // 0987***321
//	_ = s.Reveal()        // 0987654321
func NewPhone(v string, opts ...SensitiveOption) Sensitive[string] {
	return NewSensitive(v, Mobile, opts...)
}

// NewEmail 建立綁定 email 遮罩（Email）的 Sensitive[string]。
// Example:
//
//	s := masker.NewEmail("ggw.chang@gmail.com")
//	fmt.Println(s) // ggw****@gmail.com
func NewEmail(v string, opts ...SensitiveOption) Sensitive[string] {
	return NewSensitive(v, Email, opts...)
}

// NewPassword 建立綁定密碼遮罩（Password）的 Sensitive[string]。
// Example:
//
//	s := masker.NewPassword("password")
//	fmt.Println(s) // **************
func NewPassword(v string, opts ...SensitiveOption) Sensitive[string] {
	return NewSensitive(v, Password, opts...)
}

// NewID 建立綁定證號遮罩（ID）的 Sensitive[string]。
// Example:
//
//	s := masker.NewID("A123456789")
//	fmt.Println(s) // A12345****
func NewID(v string, opts ...SensitiveOption) Sensitive[string] {
	return NewSensitive(v, ID, opts...)
}

// NewCredit 建立綁定信用卡號遮罩（Credit）的 Sensitive[string]。
// Example:
//
//	s := masker.NewCredit("4111111111111111")
//	fmt.Println(s) // 411111******1111
func NewCredit(v string, opts ...SensitiveOption) Sensitive[string] {
	return NewSensitive(v, Credit, opts...)
}

// NewName 建立綁定姓名遮罩（Name）的 Sensitive[string]。
// Example:
//
//	s := masker.NewName("John Doe")
//	fmt.Println(s) // J**n D**e
func NewName(v string, opts ...SensitiveOption) Sensitive[string] {
	return NewSensitive(v, Name, opts...)
}

// NewAddress 建立綁定地址遮罩（Address）的 Sensitive[string]。
// Example:
//
//	s := masker.NewAddress("台北市內湖區內湖路一段737巷1號1樓")
//	fmt.Println(s) // 台北市內湖區******
func NewAddress(v string, opts ...SensitiveOption) Sensitive[string] {
	return NewSensitive(v, Address, opts...)
}

// NewTel 建立綁定市話遮罩（Tel）的 Sensitive[string]。
// 注意與 NewPhone（手機）不同，勿混用。
// Example:
//
//	s := masker.NewTel("0227993078")
//	fmt.Println(s) // (02)2799-****
func NewTel(v string, opts ...SensitiveOption) Sensitive[string] {
	return NewSensitive(v, Tel, opts...)
}

// NewURL 建立綁定 URL 密碼段遮罩（URL）的 Sensitive[string]。
// Example:
//
//	s := masker.NewURL("http://john:password@localhost:3000")
//	fmt.Println(s) // http://john:xxxxx@localhost:3000
func NewURL(v string, opts ...SensitiveOption) Sensitive[string] {
	return NewSensitive(v, URL, opts...)
}
