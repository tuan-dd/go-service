module github.com/tuan-dd/go-pkg/appLogger

go 1.24.2

require (
	github.com/tuan-dd/go-pkg/common v0.0.0-20250611030557-3fb6985f944f
	github.com/tuan-dd/go-pkg/settings v0.0.0-20250611030121-71f669231a62
	go.uber.org/zap v1.27.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
)

replace (
	github.com/tuan-dd/go-pkg/common => ../common
	github.com/tuan-dd/go-pkg/settings => ../settings
)
