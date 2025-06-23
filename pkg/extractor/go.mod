module github.com/tuan-dd/go-pkg/extractor

go 1.24.2

require (
	github.com/tuan-dd/go-pkg/common v0.0.0-20250611030557-3fb6985f944f
	google.golang.org/grpc v1.72.1
)

require golang.org/x/sys v0.32.0 // indirect

replace github.com/tuan-dd/go-pkg/common => ../common
