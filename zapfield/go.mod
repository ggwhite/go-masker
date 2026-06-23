module github.com/ggwhite/go-masker/v3/zapfield

go 1.22

require (
	github.com/ggwhite/go-masker/v3 v3.0.0
	go.uber.org/zap v1.28.0
)

require go.uber.org/multierr v1.10.0 // indirect

replace github.com/ggwhite/go-masker/v3 => ../
