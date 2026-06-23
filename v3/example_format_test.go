package masker_test

import (
	"fmt"

	masker "github.com/ggwhite/go-masker/v3"
)

type formatExampleFoo struct {
	Name  string            `mask:"name"`
	Email string            `mask:"email"`
	Self  *formatExampleFoo `mask:"struct"`
}

// ExampleMaskerMarshaler_Format 示範 Format() 直接輸出遮罩後的確定性字串：
// ptr-to-struct 以 &{...} 遞迴展開（不含記憶體位址），nil pointer 顯示 <nil>。
func ExampleMaskerMarshaler_Format() {
	m := masker.NewMaskerMarshaler()
	foo := &formatExampleFoo{
		Name:  "John Doe",
		Email: "john@gmail.com",
		Self: &formatExampleFoo{
			Name:  "Jane Doe",
			Email: "jane@gmail.com",
		},
	}
	fmt.Println(m.Format(foo))
	// Output: &{J**n D**e joh****@gmail.com &{J**e D**e jan****@gmail.com <nil>}}
}
