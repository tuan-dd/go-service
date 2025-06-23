package global

import (
	"github.com/tuan-dd/go-pkg/appLogger"
	"github.com/tuan-dd/go-pkg/caching/memory"
	redisClient "github.com/tuan-dd/go-pkg/caching/redis"
	"github.com/tuan-dd/go-pkg/database/mysql"
	"github.com/tuan-dd/go-pkg/extractor"
	"github.com/tuan-dd/go-pkg/settings"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/core"
)

type (
	Config struct {
		AllowOrigins string                   `mapstructure:"ALLOW_ORIGINS"`
		Server       *settings.ServerSetting  `mapstructure:"SERVER"`
		Logger       *appLogger.LoggerConfig  `mapstructure:"LOGGER"`
		Cache        *redisClient.CacheConfig `mapstructure:"CACHE"`
		SQL          *mysql.SQLConfig         `mapstructure:"SQL"`
	}
	JwtSigning struct {
		PrivateKey  string `mapstructure:"PRIVATE_KEY"`
		Issuer      string `mapstructure:"ISSUER"`
		ExpiresTime int    `mapstructure:"EXPIRES_TIME"`
	}
)

var (
	Log         *appLogger.Logger
	AppConfig   *Config
	App         *core.HttpServer
	SQLDB       *mysql.Connection
	Orm         *core.EntClient
	Cache       *redisClient.CacheClient
	MemoryCache *memory.MemoryCache
	Ext         extractor.Extractor
)
