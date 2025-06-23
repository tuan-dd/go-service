module github.com/tuan-dd/go-service/url

go 1.24.4

require (
	entgo.io/ent v0.14.4
	github.com/go-sql-driver/mysql v1.9.2
	github.com/gofiber/fiber/v3 v3.0.0-beta.4
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/stretchr/testify v1.10.0
	github.com/tuan-dd/go-pkg/appLogger v0.0.0-20250527015517-e519fcfb8809
	github.com/tuan-dd/go-pkg/caching/memory v0.0.0-20250527015517-e519fcfb8809
	github.com/tuan-dd/go-pkg/caching/redis v0.0.0-20250527015517-e519fcfb8809
	github.com/tuan-dd/go-pkg/common v0.0.0-20250611030557-3fb6985f944f
	github.com/tuan-dd/go-pkg/database/mysql v0.0.0-20250527015517-e519fcfb8809
	github.com/tuan-dd/go-pkg/extractor v0.0.0-20250527015517-e519fcfb8809
	github.com/tuan-dd/go-pkg/orm/ent v0.0.0-20250527015517-e519fcfb8809
	github.com/tuan-dd/go-pkg/settings v0.0.0-20250611030121-71f669231a62
	go.uber.org/zap v1.27.0
	golang.org/x/sys v0.32.0
)

require (
	ariga.io/atlas v0.31.1-0.20250212144724-069be8033e83 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/casbin/ent-adapter v1.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/go-jose/go-jose/v4 v4.1.0 // indirect
	github.com/go-openapi/inflect v0.21.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/gofiber/schema v1.2.0 // indirect
	github.com/gofiber/utils/v2 v2.0.0-beta.7 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/hcl/v2 v2.23.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/philhofer/fwd v1.1.3-0.20240916144458-20a13a1f6b7c // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/redis/go-redis/v9 v9.7.3 // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/samber/lo v1.49.1 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/sqids/sqids-go v0.4.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/tinylib/msgp v1.2.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.58.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/zclconf/go-cty v1.16.2 // indirect
	github.com/zclconf/go-cty-yaml v1.1.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/mod v0.23.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	golang.org/x/tools v0.30.0 // indirect
	google.golang.org/grpc v1.72.1 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/tuan-dd/go-pkg/appLogger => ../../pkg/appLogger
	github.com/tuan-dd/go-pkg/caching/memory => ../../pkg/caching/memory
	github.com/tuan-dd/go-pkg/caching/redis => ../../pkg/caching/redis
	github.com/tuan-dd/go-pkg/common => ../../pkg/common
	github.com/tuan-dd/go-pkg/database/mysql => ../../pkg/database/mysql
	github.com/tuan-dd/go-pkg/extractor => ../../pkg/extractor
	github.com/tuan-dd/go-pkg/fiberzap-logger => ../../pkg/fiberzap-logger
	github.com/tuan-dd/go-pkg/orm/ent => ../../pkg/orm/ent
	github.com/tuan-dd/go-pkg/settings => ../../pkg/settings
)
