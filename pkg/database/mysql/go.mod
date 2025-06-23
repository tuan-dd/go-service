module github.com/tuan-dd/go-pkg/database/mysql

go 1.24.2

require (
	github.com/go-sql-driver/mysql v1.9.2
	github.com/tuan-dd/go-pkg/common v0.0.0-20250611030557-3fb6985f944f
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
)

replace github.com/tuan-dd/go-pkg/common => ../../common
