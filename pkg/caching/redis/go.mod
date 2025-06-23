module github.com/tuan-dd/go-pkg/caching/redis

go 1.24.2

require (
	github.com/redis/go-redis/v9 v9.7.3
	github.com/tuan-dd/go-pkg/common v0.0.0-20250611030557-3fb6985f944f
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/uuid v1.6.0 // indirect
)

replace github.com/tuan-dd/go-pkg/common => ../../common
